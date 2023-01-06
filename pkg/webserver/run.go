package webserver

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/mattfenwick/scaling/pkg/database"
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

func Run(port int, tp trace.TracerProvider, db *sql.DB) {
	addr := fmt.Sprintf(":%d", port)

	rootContext := context.Background()
	ctx, cancel := context.WithTimeout(rootContext, 10*time.Second)
	defer cancel()
	utils.DoOrDie(database.InitializeSchema(ctx, db))

	model := NewModel(rootContext, tp, db)
	serveMux := SetupHTTPServer(model, tp)

	logrus.Infof("listening on port %s", addr)
	utils.DoOrDie(http.ListenAndServe(addr, serveMux))
}
