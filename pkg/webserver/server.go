package webserver

import (
	"context"
	"fmt"
	"github.com/mattfenwick/scaling/pkg/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

type Responder interface {
	Respond(ctx context.Context, path string, method string, body []byte, values url.Values) (string, int, error)
}

func SetupHTTPServer(responder Responder) *http.ServeMux {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var code int
		var response string

		start := time.Now()
		defer func() {
			telemetry.RecordAPIDuration(r.URL.Path, r.Method, code, start)
		}()

		logrus.Infof("handling request: %s to %s", r.Method, r.URL.Path)
		ctx := r.Context()
		span := trace.SpanFromContext(ctx)
		span.AddEvent("handler")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logrus.Errorf("http server error %s with code %d (%s to %s, r.Method, r.URL.Path)", err.Error(), 400, r.Method, r.URL.Path)
			http.Error(w, err.Error(), 400)
			return
		}

		response, code, err = responder.Respond(ctx, r.URL.Path, r.Method, body, r.URL.Query())
		logrus.Debugf("response: %s; code: %d; err? %t", response, code, err != nil)
		if err != nil {
			logrus.Errorf("http error: %s to %s, code %d, error %+v", r.Method, r.URL.Path, code, err)
			http.Error(w, err.Error(), code)
		} else if code == 404 {
			logrus.Errorf("http not found: %s to %s, code %d, error %+v", r.Method, r.URL.Path, code, err)
			http.NotFound(w, r)
		}

		header := w.Header()
		w.WriteHeader(code)
		header.Set(http.CanonicalHeaderKey("content-type"), "application/json")
		_, _ = fmt.Fprint(w, response)
	}

	serveMux := http.NewServeMux()
	//serveMux.HandleFunc("/", handler)
	serveMux.Handle("/", otelhttp.NewHandler(http.HandlerFunc(handler), "handle"))
	return serveMux
}
