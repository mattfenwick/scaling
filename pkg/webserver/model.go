package webserver

import (
	"context"
	"fmt"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/sirupsen/logrus"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Model struct {
	IsLive  bool
	IsReady bool
}

func NewModel() *Model {
	return &Model{
		IsLive:  true,
		IsReady: true,
	}
}

func (m *Model) Respond(ctx context.Context, path string, method string, body []byte, values url.Values) (string, int, error) {
	pathPieces := getPathPieces(path)
	logrus.Infof("responder: handling request: %s to %s, body %+v, path pieces [%+v]", method, path, body, pathPieces)
	if slice.EqualSlicePairwise[string]()(pathPieces, []string{"liveness"}) {
		if m.IsLive {
			return "liveness", 200, nil
		} else {
			return "not live", 500, nil
		}
	} else if slice.EqualSlicePairwise[string]()(pathPieces, []string{"readiness"}) {
		if m.IsReady {
			return "readiness", 200, nil
		} else {
			return "not ready", 500, nil
		}
	} else if slice.EqualSlicePairwise[string]()(pathPieces, []string{"hack", "wait"}) {
		secondsString := values.Get("seconds")
		seconds, err := strconv.Atoi(secondsString)
		if err != nil {
			return "invalid seconds", 400, err
		}
		if seconds < 0 || seconds > 10 {
			return "seconds out of bounds", 400, nil
		}
		time.Sleep(time.Duration(seconds) * time.Second)
		return fmt.Sprintf("waited %d seconds", seconds), 200, nil
	} else if slice.EqualSlicePairwise[string]()(pathPieces, []string{"hack", "kill"}) {
		if m.IsLive {
			m.IsLive = false
			return "killed", 200, nil
		}
		return "can't kill, already dead", 400, nil
	}
	return json.MustMarshalToString(map[string]string{"status": "TODO"}), 500, nil
}

func getPathPieces(path string) []string {
	return slice.Filter(func(p string) bool { return len(p) > 0 }, strings.Split(path, "/"))
}
