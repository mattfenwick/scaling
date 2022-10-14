package webserver

import (
	"fmt"
	"github.com/mattfenwick/scaling/pkg/utils"
	"net/http"
)

func Run(port int) {
	addr := fmt.Sprintf(":%d", port)
	model := &Model{}
	serveMux := SetupHTTPServer(model)

	utils.DoOrDie(http.ListenAndServe(addr, serveMux))
}
