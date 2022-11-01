package webserver

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

func Run(port int, tp trace.TracerProvider, db *sql.DB) {
	addr := fmt.Sprintf(":%d", port)
	model := NewModel(context.TODO(), tp, db)
	serveMux := SetupHTTPServer(model, tp)

	logrus.Infof("listening on port %s", addr)
	utils.DoOrDie(http.ListenAndServe(addr, serveMux))
}
