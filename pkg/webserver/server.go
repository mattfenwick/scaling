package webserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/telemetry"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/sirupsen/logrus"
)

type Responder interface {
	Sleep(ctx context.Context, seconds string) error

	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	CreateMessage(context.Context, *CreateMessageRequest) (*CreateMessageResponse, error)
	Follow(context.Context, *FollowRequest) (*FollowResponse, error)
	CreateUpvote(context.Context, *CreateUpvoteRequest) (*CreateUpvoteResponse, error)

	IsLive(context.Context) bool
	IsReady(context.Context) bool

	DocumentsFetchAll(context.Context) (*GetAllDocumentsResponse, error)
	DocumentsFind(context.Context, *FindDocumentsRequest) (*FindDocumentsResponse, error)
	DocumentFetch(context.Context, *GetDocumentRequest) (*GetDocumentResponse, error)
	DocumentUpload(context.Context, *UploadDocumentRequest) (*UploadDocumentResponse, error)
	Dump(ctx context.Context) (string, error)
}

func RequestHandler(r *http.Request, process func(ctx context.Context, body string, urlParams url.Values) (any, error)) (int, any, error) {
	logrus.Infof("handling request: %s to %s", r.Method, r.URL.Path)

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
			http.Error(w, "method not allowed", 405)
			return
		}

		code, response, err = RequestHandler(r, handler)

		logrus.Debugf("response code: %d; err? %t", code, err != nil)
		if err != nil {
			logrus.Errorf("http error: %s to %s, code %d, error %+v", r.Method, r.URL.Path, code, err)
			http.Error(w, err.Error(), code)
			return
		} else if response == nil {
			code = 404
			logrus.Errorf("http not found: %s to %s, code %d, error %+v", r.Method, r.URL.Path, code, err)
			http.NotFound(w, r)
			return
		}

		header := w.Header()
		w.WriteHeader(code)
		header.Set(http.CanonicalHeaderKey("content-type"), "application/json")
		body := json.MustMarshalToString(response)
		n, err := fmt.Fprint(w, body)
		if err != nil {
			logrus.Errorf("unable to print response body: %s", err.Error())
		} else if n < len(body) {
			logrus.Errorf("failed to print full body: %d / %d", n, len(body))
		} else {
			logrus.Infof("wrote %d / %d bytes to response successfully", n, len(body))
		}
	}
}

const (
	// kubernetes
	LivenessPath  = "/liveness"
	ReadinessPath = "/readiness"

	// core model
	UsersPath     = "/users"
	MessagesPath  = "/messages"
	FollowersPath = "/followers"
	UpvotesPath   = "/upvotes"

	// hacks
	DumpPath  = "/dump"
	SleepPath = "/sleep"

	// documents
	DocumentsPath     = "/documents"
	AllDocumentsPath  = "/documents/all"
	FindDocumentsPath = "/documents/find"
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
	serveMux.Handle(UsersPath, otelhttp.NewHandler(http.HandlerFunc(Handler(1000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"POST": func(ctx context.Context, body string, values url.Values) (any, error) {
				user, err := json.ParseString[CreateUserRequest](body)
				if err != nil {
					return nil, err
				}
				return responder.CreateUser(ctx, user)
			},
		})), "handle create user"))

	serveMux.Handle(MessagesPath, otelhttp.NewHandler(http.HandlerFunc(Handler(1000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"POST": func(ctx context.Context, body string, values url.Values) (any, error) {
				message, err := json.ParseString[CreateMessageRequest](body)
				if err != nil {
					return nil, err
				}
				return responder.CreateMessage(ctx, message)
			},
		})), "handle create message"))

	serveMux.Handle(FollowersPath, otelhttp.NewHandler(http.HandlerFunc(Handler(1000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"POST": func(ctx context.Context, body string, values url.Values) (any, error) {
				follow, err := json.ParseString[FollowRequest](body)
				if err != nil {
					return nil, err
				}
				return responder.Follow(ctx, follow)
			},
		})), "handle follow"))

	serveMux.Handle(UpvotesPath, otelhttp.NewHandler(http.HandlerFunc(Handler(1000,
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

	// documents
	serveMux.Handle(DocumentsPath, otelhttp.NewHandler(http.HandlerFunc(Handler(10000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"GET": func(ctx context.Context, body string, values url.Values) (any, error) {
				return responder.DocumentFetch(ctx, &GetDocumentRequest{
					DocumentId: values.Get("id"),
				})
			},
			"POST": func(ctx context.Context, body string, values url.Values) (any, error) {
				udr, err := json.ParseString[UploadDocumentRequest](body)
				if err != nil {
					return nil, err
				}
				return responder.DocumentUpload(ctx, udr)
			},
		})), "handle document"))

	serveMux.Handle(AllDocumentsPath, otelhttp.NewHandler(http.HandlerFunc(Handler(0,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"GET": func(ctx context.Context, body string, values url.Values) (any, error) {
				return responder.DocumentsFetchAll(ctx)
			},
		})), "handle fetch all documents"))

	serveMux.Handle(FindDocumentsPath, otelhttp.NewHandler(http.HandlerFunc(Handler(10000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"POST": func(ctx context.Context, body string, values url.Values) (any, error) {
				fdr, err := json.ParseString[FindDocumentsRequest](body)
				if err != nil {
					return nil, err
				}
				return responder.DocumentsFind(ctx, fdr)
			},
		})), "handle find documents"))

	return serveMux
}
