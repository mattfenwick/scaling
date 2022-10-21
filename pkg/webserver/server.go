package webserver

import (
	"context"
	"fmt"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

type Responder interface {
	DocumentUnsafeFetchAll(context.Context) (*UnsafeGetDocumentsResponse, error)
	DocumentFetch(context.Context, *GetDocumentRequest) (*GetDocumentResponse, error)
	DocumentUpload(context.Context, *UploadDocumentRequest) (*UploadDocumentResponse, error)

	LivenessCode(context.Context) int
	ReadinessCode(context.Context) int
}

func RequestHandler(r *http.Request, process func(ctx context.Context, body string, urlParams url.Values) (any, error)) (int, any, error) {
	var code int
	var response any

	start := time.Now()
	defer func() {
		telemetry.RecordAPIDuration(r.URL.Path, r.Method, code, start)
	}()

	logrus.Infof("handling request: %s to %s", r.Method, r.URL.Path)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		code = 400
		logrus.Errorf("http server error %s with code %d (%s to %s, r.Method, r.URL.Path)", err.Error(), code, r.Method, r.URL.Path)
		return code, nil, err
	}

	ctx := r.Context()
	//span := trace.SpanFromContext(ctx)
	//span.AddEvent("handler")
	response, err = process(ctx, string(body), r.URL.Query())
	logrus.Debugf("response: %s; code: %d; err? %t", response, code, err != nil)
	if err != nil {
		code = 500
		logrus.Errorf("http error: %s to %s, code %d, error %+v", r.Method, r.URL.Path, code, err)
		return code, nil, err
	}

	code = 200
	return code, response, nil
}

func Handler(maxSize int64, methodHandlers map[string]func(ctx context.Context, body string, values url.Values) (any, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		code, response, err := RequestHandler(r, handler)

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
	LivenessPath        = "/liveness"
	ReadinessPath       = "/readiness"
	DocumentsPath       = "/documents"
	UnsafeDocumentsPath = "/unsafe/documents"
)

func SetupHTTPServer(responder Responder) *http.ServeMux {
	serveMux := http.NewServeMux()
	//serveMux.Handle("/", otelhttp.NewHandler(http.HandlerFunc(handler), "handle"))

	serveMux.Handle(LivenessPath, otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(responder.LivenessCode(r.Context()))
	}), "handle liveness"))

	serveMux.Handle(ReadinessPath, otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(responder.ReadinessCode(r.Context()))
	}), "handle readiness"))

	serveMux.Handle(DocumentsPath, otelhttp.NewHandler(http.HandlerFunc(Handler(10000,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"GET": func(ctx context.Context, body string, values url.Values) (any, error) {
				return responder.DocumentFetch(ctx, &GetDocumentRequest{
					DocumentId: values.Get("id"),
				})
			},
			"POST": func(ctx context.Context, body string, values url.Values) (any, error) {
				return responder.DocumentUpload(ctx, &UploadDocumentRequest{
					Document: body,
				})
			},
		})), "handle document"))

	serveMux.Handle(UnsafeDocumentsPath, otelhttp.NewHandler(http.HandlerFunc(Handler(0,
		map[string]func(ctx context.Context, body string, values url.Values) (any, error){
			"GET": func(ctx context.Context, body string, values url.Values) (any, error) {
				return responder.DocumentUnsafeFetchAll(ctx)
			},
		})), "handle unsafe document"))

	return serveMux
}
