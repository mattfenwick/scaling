package webserver

import (
	"context"
	"database/sql"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/database"
	"github.com/mattfenwick/scaling/pkg/telemetry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Document struct {
	Id     string
	Raw    string
	Parsed any
	Error  string
}

type Action struct {
	Name string
	F    func() error
}

type Model struct {
	Documents map[string]*Document
	Live      bool
	Ready     bool
	db        *sql.DB
	tp        trace.TracerProvider
	tracer    trace.Tracer
	actions   chan *Action
}

func NewModel(ctx context.Context, tp trace.TracerProvider, db *sql.DB) *Model {
	actions := make(chan *Action, 1)
	m := &Model{
		Documents: map[string]*Document{},
		Live:      true,
		Ready:     true,
		db:        db,
		tp:        tp,
		tracer:    tp.Tracer("model"),
		actions:   actions,
	}
	go func() {
		for {
			//logrus.Debugf("state: %s", json.MustMarshalToString(m))
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

	logrus.Debugf("attempting to parse object of length %d", len(request.Document))
	logrus.Tracef("attempting to parse object: %s", request.Document)

	parsed, parseErr := json.ParseString[any](request.Document)

	errorString := ""
	var derefedParsed any
	if parseErr != nil {
		errorString = parseErr.Error()
	} else {
		derefedParsed = *parsed
	}

	m.Documents[id] = &Document{
		Id:     id,
		Raw:    request.Document,
		Parsed: derefedParsed,
		Error:  errorString,
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
				Id:     doc.Id,
				Raw:    doc.Raw,
				Parsed: doc.Parsed,
				Error:  doc.Error,
			}
		}
		wg.Done()
		return nil
	}

	select {
	case m.actions <- &Action{F: action, Name: "fetch all documents"}:
		wg.Wait()
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
			logrus.Debugf("looking for key '%s' in document %s", request.Key, id)
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
		wg.Wait()
		return &FindDocumentsResponse{
			Matches: items,
		}, nil
	default:
		return nil, errors.Errorf("service unavailable")
	}
}

func (m *Model) Dump(ctx context.Context) (string, error) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var out string

	action := func() error {
		out = json.MustMarshalToString(m.Documents)
		wg.Done()
		return nil
	}

	select {
	case m.actions <- &Action{F: action, Name: "dump"}:
		wg.Wait()
		return out, nil
	default:
		return "", errors.Errorf("service unavailable")
	}
}

func (m *Model) IsLive(ctx context.Context) bool {
	return m.Live
}

func (m *Model) IsReady(ctx context.Context) bool {
	return m.Ready
}

func (m *Model) Sleep(ctx context.Context, milliseconds string) error {
	ms, err := strconv.Atoi(milliseconds)
	if err != nil {
		return errors.Wrapf(err, "unable to parse milliseconds: '%s'", milliseconds)
	}
	if ms <= 0 || ms > 5000 {
		return errors.Errorf("milliseconds '%d' out of range", ms)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	action := func() error {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return nil
	}

	select {
	case m.actions <- &Action{F: action, Name: "sleep"}:
		wg.Wait()
		return nil
	default:
		return errors.Errorf("service unavailable")
	}
}

func (m *Model) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	newUser := database.NewUser(req.Name, req.Email)
	err := database.InsertUser(ctx, m.db, newUser)
	if err != nil {
		return nil, err
	}
	return &CreateUserResponse{UserId: newUser.UserId}, nil
}

func (m *Model) CreateMessage(ctx context.Context, req *CreateMessageRequest) (*CreateMessageResponse, error) {
	newMessage := database.NewMessage(req.SenderUserId, req.Content)
	err := database.InsertMessage(ctx, m.db, newMessage)
	if err != nil {
		return nil, err
	}
	return &CreateMessageResponse{MessageId: newMessage.MessageId}, nil
}

func (m *Model) Follow(ctx context.Context, req *FollowRequest) (*FollowResponse, error) {
	newFollower := database.NewFollower(req.FolloweeUserId, req.FollowerUserId)
	err := database.InsertFollower(ctx, m.db, newFollower)
	if err != nil {
		return nil, err
	}
	return &FollowResponse{}, nil
}

func (m *Model) CreateUpvote(ctx context.Context, req *CreateUpvoteRequest) (*CreateUpvoteResponse, error) {
	newUpvote := database.NewUpvote(req.UserId, req.MessageId)
	err := database.InsertUpvote(ctx, m.db, newUpvote)
	if err != nil {
		return nil, err
	}
	return &CreateUpvoteResponse{UpvoteId: newUpvote.UpvoteId}, nil
}
