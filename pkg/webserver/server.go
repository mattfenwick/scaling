package webserver

import (
	"context"
	"fmt"
	"github.com/mattfenwick/collections/pkg/json"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Function struct {
	Name string
	Args []int
}

type FunctionResult struct {
	Value int
}

type Responder interface {
	RunFunctionHttp(ctx context.Context, function *Function) (*FunctionResult, error)

	NotFound(w http.ResponseWriter, r *http.Request)
	Error(w http.ResponseWriter, r *http.Request, err error, statusCode int)
}

func SetupHTTPServer(responder Responder) {
	handleJob := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := trace.SpanFromContext(ctx)
		span.AddEvent("handling-function")

		logrus.Infof("handling function request")
		switch r.Method {
		case "POST":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				logrus.Errorf("unable to read body for RunFunction POST: %s", err.Error())
				responder.Error(w, r, err, 400)
				return
			}
			logrus.Debugf("request body: <%s>", body)

			f, err := json.Parse[Function](body)
			if err != nil {
				logrus.Errorf("unable to ummarshal JSON for RunFunction POST: %s", err.Error())
				responder.Error(w, r, err, 400)
				return
			}
			jobStatus, err := responder.RunFunctionHttp(ctx, f)
			if err != nil {
				http.Error(w, err.Error(), 400)
			} else {
				header := w.Header()
				header.Set(http.CanonicalHeaderKey("content-type"), "application/json")
				fmt.Fprint(w, json.MustMarshalToString(jobStatus))
			}
		default:
			responder.NotFound(w, r)
		}
	}
	http.Handle("/function", otelhttp.NewHandler(http.HandlerFunc(handleJob), "function"))
}
