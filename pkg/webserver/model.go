package webserver

import (
	"context"
	"github.com/google/uuid"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/pkg/errors"
)

type Model struct {
	Documents map[string]string
	IsLive    bool
	IsReady   bool
}

func NewModel() *Model {
	return &Model{
		IsLive:  true,
		IsReady: true,
	}
}

//func (m *Model) Respond(ctx context.Context, path string, method string, body []byte, values url.Values) (string, int, error) {
//	pathPieces := getPathPieces(path)
//	logrus.Infof("responder: handling request: %s to %s, body %+v, path pieces [%+v]", method, path, body, pathPieces)
//	if slice.EqualSlicePairwise[string]()(pathPieces, []string{"liveness"}) {
//		if m.IsLive {
//			return "liveness", 200, nil
//		} else {
//			return "not live", 500, nil
//		}
//	} else if slice.EqualSlicePairwise[string]()(pathPieces, []string{"readiness"}) {
//		if m.IsReady {
//			return "readiness", 200, nil
//		} else {
//			return "not ready", 500, nil
//		}
//	} else if slice.EqualSlicePairwise[string]()(pathPieces, []string{"hack", "wait"}) {
//		secondsString := values.Get("seconds")
//		seconds, err := strconv.Atoi(secondsString)
//		if err != nil {
//			return "invalid seconds", 400, err
//		}
//		if seconds < 0 || seconds > 10 {
//			return "seconds out of bounds", 400, nil
//		}
//		time.Sleep(time.Duration(seconds) * time.Second)
//		return fmt.Sprintf("waited %d seconds", seconds), 200, nil
//	} else if slice.EqualSlicePairwise[string]()(pathPieces, []string{"hack", "kill"}) {
//		if m.IsLive {
//			m.IsLive = false
//			return "killed", 200, nil
//		}
//		return "can't kill, already dead", 400, nil
//	}
//	return json.MustMarshalToString(map[string]string{"status": "TODO"}), 500, nil
//}

func (m *Model) DocumentUpload(ctx context.Context, doc string) (string, error) {
	id := uuid.New().String()
	if _, ok := m.Documents[id]; ok {
		return "", errors.Errorf("cannot create doc with uuid %s: id already found", id)
	}
	m.Documents[id] = doc
	return json.MustMarshalToString(map[string]string{"id": id}), nil
}

func (m *Model) DocumentFetch(ctx context.Context, id string) (string, error) {
	if id == "" {
		return "", errors.Errorf("invalid id: empty")
	}
	doc, ok := m.Documents[id]
	if !ok {
		return "", errors.Errorf("document %s not found", id)
	}
	return json.MustMarshalToString(map[string]string{"document": doc}), nil
}

func (m *Model) DocumentUnsafeFetchAll(ctx context.Context) (string, error) {
	return json.MustMarshalToString(map[string]interface{}{"documents": m.Documents}), nil
}

func (m *Model) LivenessCode(ctx context.Context) int {
	if m.IsLive {
		return 200
	} else {
		return 500
	}
}

func (m *Model) ReadinessCode(ctx context.Context) int {
	if m.IsReady {
		return 200
	} else {
		return 500
	}
}
