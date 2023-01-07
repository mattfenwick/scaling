package webserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/telemetry"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/sirupsen/logrus"
)

func isNil(v any) bool {
	// source:
	//   https://gist.github.com/miguelmota/faca748b3c8598f2abf322b51b542d24
	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}

func RequestHandler(r *http.Request, process func(ctx context.Context, body string, urlParams url.Values) (any, error)) (int, any, error) {
	logrus.Debugf("handling request: %s to %s", r.Method, r.URL.Path)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return 400, nil, err
	}

	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	childCtx, childCancel := context.WithTimeout(ctx, 5*time.Second)
	defer childCancel()

	span.AddEvent("start process")
	response, err := process(childCtx, string(body), r.URL.Query())
	span.AddEvent("finish process")

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 500, nil, err
	}

	return 200, response, nil
}

func Handler(maxSize int64, methodHandlers map[string]func(ctx context.Context, body string, values url.Values) (any, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var code int
		var response any
		var err error

		start := time.Now()
		defer func() {
			telemetry.RecordAPIDuration(r.URL.Path, r.Method, code, start)
		}()

		if r.ContentLength > maxSize {
			logrus.Errorf("content length too large")
			http.Error(w, "content length too large", 400)
			return
		}

		handler, ok := methodHandlers[r.Method]
		if !ok {
			logrus.Errorf("method %s not allowed for %s", r.Method, r.URL.Path)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		code, response, err = RequestHandler(r, handler)
		logrus.Debugf("handled %s to %s: response %+v (is nil? %t) (provisional code %d), err %+v", r.Method, r.URL.Path, response, isNil(response), code, err)

		logrus.Debugf("response code: %d; err? %t", code, err != nil)
		if err != nil {
			logrus.Errorf("http error: %s to %s, code %d, error %+v", r.Method, r.URL.Path, code, err)
			http.Error(w, err.Error(), code)
			return
		} else if isNil(response) { // response == nil {
			code = 404
			logrus.Errorf("http not found: %s to %s, code %d, error %+v", r.Method, r.URL.Path, code, err)
			http.NotFound(w, r)
			return
		}

		header := w.Header()
		w.WriteHeader(code)
		// header.Set(http.CanonicalHeaderKey("content-type"), "application/json")
		header.Set("content-type", "application/json")
		body := json.MustMarshalToString(response)
		n, err := fmt.Fprint(w, body)
		if err != nil {
			logrus.Errorf("unable to print response body: %s", err.Error())
		} else if n < len(body) {
			logrus.Errorf("failed to print full body: %d / %d", n, len(body))
		} else {
			logrus.Debugf("wrote %d / %d bytes to response successfully", n, len(body))
		}
	}
}

const (
	// kubernetes
	LivenessPath  = "/liveness"
	ReadinessPath = "/readiness"

	// core model
	UserPath         = "/user"
	UserTimelinePath = "/user/timeline"
	UserMessagesPath = "/user/messages"
	UsersPath        = "/users"
	MessagePath      = "/message"
	MessagesPath     = "/messages"
	FollowPath       = "/follow"
	FollowersPath    = "/followers"
	UpvotePath       = "/upvote"

	// hacks
	DumpPath  = "/dump"
	SleepPath = "/sleep"
)

func SetupHTTPServer(responder Responder, tp trace.TracerProvider) *http.ServeMux {
	serveMux := http.NewServeMux()
	//serveMux.Handle("/", otelhttp.NewHandler(http.HandlerFunc(handler), "handle"))

	// kubernetes
	serveMux.Handle(LivenessPath, otelhttp.NewHandler(http.HandlerFunc(Handler(0,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"GET": func(ctx context.Context, body string, values url.Values) (any, error) {
				if responder.IsLive(ctx) {
					return "", nil
				} else {
					return nil, errors.Errorf("not live")
				}
			},
		})), "handle liveness"))

	serveMux.Handle(ReadinessPath, otelhttp.NewHandler(http.HandlerFunc(Handler(0,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"GET": func(ctx context.Context, body string, values url.Values) (any, error) {
				if responder.IsReady(ctx) {
					return "", nil
				} else {
					return nil, errors.Errorf("not ready")
				}
			},
		})), "handle readiness"))

	// core model
	serveMux.Handle(UserPath, otelhttp.NewHandler(http.HandlerFunc(Handler(1000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"POST": func(ctx context.Context, body string, values url.Values) (any, error) {
				user, err := json.ParseString[CreateUserRequest](body)
				if err != nil {
					return nil, err
				}
				return responder.CreateUser(ctx, user)
			},
			"GET": func(ctx context.Context, body string, values url.Values) (any, error) {
				userId, err := uuid.Parse(values.Get("userid"))
				if err != nil {
					return nil, errors.Wrapf(err, "unable to parse uuid from '%s'", values.Get("userid"))
				}
				request := &GetUserRequest{
					UserId: userId,
				}
				return responder.GetUser(ctx, request)
			},
		})), "handle user"))

	serveMux.Handle(UsersPath, otelhttp.NewHandler(http.HandlerFunc(Handler(1000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"GET": func(ctx context.Context, body string, values url.Values) (any, error) {
				req := &GetUsersRequest{}
				return responder.GetUsers(ctx, req)
			},
			"POST": func(ctx context.Context, body string, values url.Values) (any, error) {
				req, err := json.ParseString[SearchUsersRequest](body)
				if err != nil {
					return nil, err
				}
				return responder.SearchUsers(ctx, req)
			},
		})), "handle users"))

	serveMux.Handle(UserTimelinePath, otelhttp.NewHandler(http.HandlerFunc(Handler(1000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"POST": func(ctx context.Context, body string, values url.Values) (any, error) {
				req, err := json.ParseString[GetUserTimelineRequest](body)
				if err != nil {
					return nil, err
				}
				return responder.GetUserTimeline(ctx, req)
			},
		})), "handle user timeline"))

	serveMux.Handle(UserMessagesPath, otelhttp.NewHandler(http.HandlerFunc(Handler(1000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"POST": func(ctx context.Context, body string, values url.Values) (any, error) {
				req, err := json.ParseString[GetUserMessagesRequest](body)
				if err != nil {
					return nil, err
				}
				return responder.GetUserMessages(ctx, req)
			},
		})), "handle user messages"))

	serveMux.Handle(MessagePath, otelhttp.NewHandler(http.HandlerFunc(Handler(1000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"POST": func(ctx context.Context, body string, values url.Values) (any, error) {
				message, err := json.ParseString[CreateMessageRequest](body)
				if err != nil {
					return nil, err
				}
				return responder.CreateMessage(ctx, message)
			},
		})), "handle create message"))

	serveMux.Handle(MessagesPath, otelhttp.NewHandler(http.HandlerFunc(Handler(1000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"GET": func(ctx context.Context, body string, values url.Values) (any, error) {
				req, err := json.ParseString[GetMessagesRequest](body) // TODO wrong, use values
				if err != nil {
					return nil, err
				}
				return responder.GetMessages(ctx, req)
			},
			"POST": func(ctx context.Context, body string, values url.Values) (any, error) {
				req, err := json.ParseString[SearchMessagesRequest](body)
				if err != nil {
					return nil, err
				}
				return responder.SearchMessages(ctx, req)
			},
		})), "handle messages"))

	serveMux.Handle(FollowPath, otelhttp.NewHandler(http.HandlerFunc(Handler(1000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"POST": func(ctx context.Context, body string, values url.Values) (any, error) {
				follow, err := json.ParseString[FollowRequest](body)
				if err != nil {
					return nil, err
				}
				return responder.Follow(ctx, follow)
			},
		})), "handle follow"))

	serveMux.Handle(FollowersPath, otelhttp.NewHandler(http.HandlerFunc(Handler(1000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"GET": func(ctx context.Context, body string, values url.Values) (any, error) {
				userId, err := uuid.Parse(values.Get("userid"))
				if err != nil {
					return nil, errors.Wrapf(err, "unable to parse uuid from '%s'", values.Get("userid"))
				}
				return responder.GetFollowers(ctx, &GetFollowersOfUserRequest{UserId: userId})
			},
		})), "handle followers"))

	serveMux.Handle(UpvotePath, otelhttp.NewHandler(http.HandlerFunc(Handler(1000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"POST": func(ctx context.Context, body string, values url.Values) (any, error) {
				upvote, err := json.ParseString[CreateUpvoteRequest](body)
				if err != nil {
					return nil, err
				}
				return responder.CreateUpvote(ctx, upvote)
			},
		})), "handle create upvote"))

	// hacks
	serveMux.Handle(DumpPath, otelhttp.NewHandler(http.HandlerFunc(Handler(0,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"GET": func(ctx context.Context, body string, values url.Values) (any, error) {
				return responder.Dump(ctx)
			},
		})), "handle dump"))

	serveMux.Handle(SleepPath, otelhttp.NewHandler(http.HandlerFunc(Handler(0,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"GET": func(ctx context.Context, body string, values url.Values) (any, error) {
				return "", responder.Sleep(ctx, values.Get("seconds"))
			},
		})), "handle sleep"))

	return serveMux
}
