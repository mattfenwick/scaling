package webserver

import (
	"context"
	"fmt"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/telemetry"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

type Responder interface {
	DocumentsFetchAll(context.Context) (*GetAllDocumentsResponse, error)
	DocumentsFind(context.Context, *FindDocumentsRequest) (*FindDocumentsResponse, error)
	DocumentFetch(context.Context, *GetDocumentRequest) (*GetDocumentResponse, error)
	DocumentUpload(context.Context, *UploadDocumentRequest) (*UploadDocumentResponse, error)

	IsLive(context.Context) bool
	IsReady(context.Context) bool
}

func RequestHandler(r *http.Request, process func(ctx context.Context, body string, urlParams url.Values) (any, error)) (int, any, error) {
	logrus.Infof("handling request: %s to %s", r.Method, r.URL.Path)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return 400, nil, err
	}

	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	span.AddEvent("start process")
	response, err := process(ctx, string(body), r.URL.Query())
	span.AddEvent("finish process")

	logrus.Debugf("response: %s; err? %t", response, err != nil)
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

		logrus.Debugf("response: %+v; code: %d; err? %t", response, code, err != nil)
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
		_, _ = fmt.Fprint(w, json.MustMarshalToString(response))
	}
}

const (
	LivenessPath      = "/liveness"
	ReadinessPath     = "/readiness"
	DocumentsPath     = "/documents"
	AllDocumentsPath  = "/documents/all"
	FindDocumentsPath = "/documents/find"
)

func SetupHTTPServer(responder Responder, tp trace.TracerProvider) *http.ServeMux {
	serveMux := http.NewServeMux()
	//serveMux.Handle("/", otelhttp.NewHandler(http.HandlerFunc(handler), "handle"))

	serveMux.Handle(LivenessPath, otelhttp.NewHandler(http.HandlerFunc(Handler(10000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"GET": func(ctx context.Context, body string, values url.Values) (any, error) {
				if responder.IsLive(ctx) {
					return "", nil
				} else {
					return nil, errors.Errorf("not live")
				}
			},
		})), "handle liveness"))

	serveMux.Handle(ReadinessPath, otelhttp.NewHandler(http.HandlerFunc(Handler(10000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"GET": func(ctx context.Context, body string, values url.Values) (any, error) {
				if responder.IsReady(ctx) {
					return "", nil
				} else {
					return nil, errors.Errorf("not ready")
				}
			},
		})), "handle readiness"))

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

	serveMux.Handle(FindDocumentsPath, otelhttp.NewHandler(http.HandlerFunc(Handler(0,
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
