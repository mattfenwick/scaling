package webserver

import (
	"context"
	"fmt"
	"github.com/mattfenwick/scaling/pkg/utils"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func Run(port int, tp trace.TracerProvider) {
	addr := fmt.Sprintf(":%d", port)
	model := NewModel(tp, context.TODO())
	serveMux := SetupHTTPServer(model, tp)

	utils.DoOrDie(http.ListenAndServe(addr, serveMux))
}
