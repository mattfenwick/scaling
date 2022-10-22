package webserver

import (
	"context"
	"github.com/google/uuid"
	"github.com/mattfenwick/gunparse/pkg"
	"github.com/mattfenwick/gunparse/pkg/example"
	"github.com/mattfenwick/scaling/pkg/parse"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"sync"
)

type Document struct {
	Id      string
	Raw     string
	Parsed  *example.Object
	IsValid bool
}

type Model struct {
	Documents map[string]*Document
	Live      bool
	Ready     bool
	tp        trace.TracerProvider
	tracer    trace.Tracer
	actions   chan func()
}

func NewModel(tp trace.TracerProvider, ctx context.Context) *Model {
	actions := make(chan func(), 1)
	m := &Model{
		Documents: map[string]*Document{},
		Live:      true,
		Ready:     true,
		tp:        tp,
		tracer:    tp.Tracer("model"),
		actions:   actions,
	}
	go func() {
		for {
			select {
			case f := <-actions:
				f()
			case <-ctx.Done():
				return
			}
		}
	}()
	return m
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

func (m *Model) DocumentUpload(ctx context.Context, request *UploadDocumentRequest) (*UploadDocumentResponse, error) {
	wg := sync.WaitGroup{}
	var result *UploadDocumentResponse
	var err error
	wg.Add(1)
	action := func() {
		result, err = m.noncurrentDocumentUpload(ctx, request)
		wg.Done()
	}

	_, span := m.tracer.Start(ctx, "run action")
	defer span.End()

	select {
	case m.actions <- action:
		wg.Wait()
		span.AddEvent("finished action run")
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		return result, err
	default:
		logrus.Warnf("service unavailable")
		span.SetStatus(codes.Error, "service unavailable")
		return nil, errors.Errorf("service unavailable")
	}
}

func (m *Model) noncurrentDocumentUpload(ctx context.Context, request *UploadDocumentRequest) (*UploadDocumentResponse, error) {
	id := uuid.New().String()
	if _, ok := m.Documents[id]; ok {
		return nil, errors.Errorf("cannot create doc with uuid %s: id already found", id)
	}

	logrus.Debugf("attemping to parse object: %s", request.Document)
	var parseResult pkg.Result[example.ParseError, *pkg.Pair[int, int], rune, *example.Object]
	parseResult = parse.JsonObject(request.Document)
	var obj *example.Object
	if parseResult.Success != nil {
		obj = parseResult.Success.Value.Result
	}
	m.Documents[id] = &Document{
		Id:      id,
		Raw:     request.Document,
		Parsed:  obj,
		IsValid: parseResult.Success != nil,
	}

	return &UploadDocumentResponse{
		DocumentId: id,
	}, nil
}

func (m *Model) DocumentFetch(ctx context.Context, request *GetDocumentRequest) (*GetDocumentResponse, error) {
	wg := sync.WaitGroup{}
	var result *GetDocumentResponse
	var err error
	wg.Add(1)
	action := func() {
		result, err = m.noncurrentDocumentFetch(ctx, request)
		wg.Done()
	}

	_, span := m.tracer.Start(ctx, "run action")
	defer span.End()

	select {
	case m.actions <- action:
		wg.Wait()
		span.AddEvent("finished action run")
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		return result, err
	default:
		logrus.Warnf("service unavailable")
		span.SetStatus(codes.Error, "service unavailable")
		return nil, errors.Errorf("service unavailable")
	}
}

func (m *Model) noncurrentDocumentFetch(ctx context.Context, request *GetDocumentRequest) (*GetDocumentResponse, error) {
	_, childSpan := m.tracer.Start(ctx, "document fetch")
	defer childSpan.End()

	id := request.DocumentId
	if id == "" {
		return nil, errors.Errorf("invalid id: empty")
	}
	doc, ok := m.Documents[id]
	if !ok {
		return nil, errors.Errorf("document %s not found", id)
	}
	return &GetDocumentResponse{
		Document: doc,
	}, nil
}

func (m *Model) DocumentUnsafeFetchAll(ctx context.Context) (*UnsafeGetDocumentsResponse, error) {
	return &UnsafeGetDocumentsResponse{
		Documents: m.Documents,
	}, nil
}

func (m *Model) IsLive(ctx context.Context) bool {
	return m.Live
}

func (m *Model) IsReady(ctx context.Context) bool {
	return m.Ready
}
