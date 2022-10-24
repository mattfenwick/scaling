package webserver

import (
	"context"
	"github.com/google/uuid"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/gunparse/pkg/example"
	"github.com/mattfenwick/scaling/pkg/parse"
	"github.com/mattfenwick/scaling/pkg/telemetry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"sync"
	"time"
)

type Document struct {
	Id        string
	Raw       string
	ParseTree *example.JsonValue
	Parsed    any
	IsValid   bool
}

type Action struct {
	Name string
	F    func() error
}

type Model struct {
	Documents map[string]*Document
	Live      bool
	Ready     bool
	tp        trace.TracerProvider
	tracer    trace.Tracer
	actions   chan *Action
}

func NewModel(tp trace.TracerProvider, ctx context.Context) *Model {
	actions := make(chan *Action, 1)
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
			case a := <-actions:
				start := time.Now()
				err := a.F()
				telemetry.RecordEventLoopDuration(a.Name, err, start)
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
	action := func() error {
		result, err = m.unsafeDocumentUpload(ctx, request)
		wg.Done()
		return err
	}

	_, span := m.tracer.Start(ctx, "run action")
	defer span.End()

	select {
	case m.actions <- &Action{F: action, Name: "upload document"}:
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

func (m *Model) unsafeDocumentUpload(ctx context.Context, request *UploadDocumentRequest) (*UploadDocumentResponse, error) {
	id := uuid.New().String()
	if _, ok := m.Documents[id]; ok {
		return nil, errors.Errorf("cannot create doc with uuid %s: id already found", id)
	}

	logrus.Debugf("attemping to parse object: %s", request.Document)

	parsed, parseErr := json.ParseString[any](request.Document)

	//var parseResult pkg.Result[example.ParseError, *pkg.Pair[int, int], rune, *example.JsonValue]
	parseResult := parse.JsonValue(request.Document)

	var jsonValue *example.JsonValue
	if parseResult.Success != nil {
		jsonValue = parseResult.Success.Value.Result
		if parseErr != nil {
			logrus.Errorf("raw: %s\ntree: %s\nerror: %s\n\n", request.Document, json.MustMarshalToString(parseResult.Success.Value.Result), parseErr.Error())
			return nil, errors.Errorf("document inconsistency: successfully parsed tree, but unable to parse into JsonValue")
		}
	} else {
		if parseErr == nil {
			logrus.Errorf("raw: %s\nerror: %s\nJsonValue: %s\n\n", request.Document, json.MustMarshalToString(parseResult.Error.Value), json.MustMarshalToString(parsed))
			return nil, errors.Errorf("document inconsistency: unable to parse tree, but successfully parsed into JsonValue")
		}
	}

	m.Documents[id] = &Document{
		Id:        id,
		Raw:       request.Document,
		ParseTree: jsonValue,
		Parsed:    parsed,
		IsValid:   parseResult.Success != nil,
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
	action := func() error {
		result, err = m.unsafeDocumentFetch(ctx, request)
		wg.Done()
		return err
	}

	_, span := m.tracer.Start(ctx, "run action")
	defer span.End()

	select {
	case m.actions <- &Action{F: action, Name: "fetch document"}:
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

func (m *Model) unsafeDocumentFetch(ctx context.Context, request *GetDocumentRequest) (*GetDocumentResponse, error) {
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

func (m *Model) DocumentsFetchAll(ctx context.Context) (*GetAllDocumentsResponse, error) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	docs := map[string]*Document{}

	action := func() error {
		for id, doc := range m.Documents {
			docs[id] = &Document{
				Id:        doc.Id,
				Raw:       doc.Raw,
				ParseTree: doc.ParseTree,
				Parsed:    doc.Parsed,
				IsValid:   doc.IsValid,
			}
		}
		wg.Done()
		return nil
	}

	select {
	case m.actions <- &Action{F: action, Name: "fetch all documents"}:
		return &GetAllDocumentsResponse{
			Documents: docs,
		}, nil
	default:
		return nil, errors.Errorf("service unavailable")
	}
}

func (m *Model) DocumentsFind(ctx context.Context, request *FindDocumentsRequest) (*FindDocumentsResponse, error) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var items []*FindDocumentsResponseItem

	action := func() error {
		for id, doc := range m.Documents {
			paths := FindKeyInJson(doc.Parsed, []any{}, request.Key)
			if len(paths) > 0 {
				items = append(items, &FindDocumentsResponseItem{
					DocumentId: id,
					Paths:      paths,
				})
			}
		}
		wg.Done()
		return nil
	}

	select {
	case m.actions <- &Action{F: action, Name: "find documents"}:
		return &FindDocumentsResponse{
			Matches: items,
		}, nil
	default:
		return nil, errors.Errorf("service unavailable")
	}
}

func (m *Model) IsLive(ctx context.Context) bool {
	return m.Live
}

func (m *Model) IsReady(ctx context.Context) bool {
	return m.Ready
}
