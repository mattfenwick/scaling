package webserver

import (
	"context"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/sirupsen/logrus"
	"net/url"
)

type Model struct {
}

func (m *Model) Respond(ctx context.Context, path string, method string, body []byte, values url.Values) (string, int, error) {
	logrus.Infof("responder: handling request: %s to %s, body %+v", method, path, body)
	return json.MustMarshalToString(map[string]string{"status": "TODO"}), 500, nil
}
