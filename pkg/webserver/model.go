package webserver

import (
	"context"
	"database/sql"
	"strconv"
	"sync"
	"time"

	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/mattfenwick/scaling/pkg/database"
	"github.com/mattfenwick/scaling/pkg/telemetry"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

type Action struct {
	Name string
	F    func() error
}

type Model struct {
	Live    bool
	Ready   bool
	db      *sql.DB
	tp      trace.TracerProvider
	tracer  trace.Tracer
	actions chan *Action
}

func NewModel(ctx context.Context, tp trace.TracerProvider, db *sql.DB) *Model {
	actions := make(chan *Action, 1)
	m := &Model{
		Live:    true,
		Ready:   true,
		db:      db,
		tp:      tp,
		tracer:  tp.Tracer("model"),
		actions: actions,
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

func (m *Model) Dump(ctx context.Context) (string, error) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var out string

	action := func() error {
		out = json.MustMarshalToString([]string{"TODO"})
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

func (m *Model) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	user, err := database.GetUser(ctx, m.db, req.UserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	return &GetUserResponse{UserId: user.UserId, Name: user.Name, Email: user.Email}, nil
}

func (m *Model) GetUsers(ctx context.Context, req *GetUsersRequest) (*GetUsersResponse, error) {
	users, err := database.GetUsers(ctx, m.db)
	if err != nil {
		return nil, err
	}
	f := func(d *database.User) GetUserResponse {
		return GetUserResponse{
			UserId: d.UserId,
			Name:   d.Name,
			Email:  d.Email,
		}
	}
	return &GetUsersResponse{Users: slice.Map(f, users)}, nil
}

func (m *Model) CreateMessage(ctx context.Context, req *CreateMessageRequest) (*CreateMessageResponse, error) {
	newMessage := database.NewMessage(req.SenderUserId, req.Content)
	err := database.InsertMessage(ctx, m.db, newMessage)
	if err != nil {
		return nil, err
	}
	return &CreateMessageResponse{MessageId: newMessage.MessageId}, nil
}

func (m *Model) GetMessage(context.Context, *GetMessageRequest) (*GetMessageResponse, error) {
	return nil, errors.Errorf("unimplemented")
}

func (m *Model) GetMessages(context.Context, *GetMessagesRequest) (*GetMessagesResponse, error) {
	return nil, errors.Errorf("unimplemented")
}

func (m *Model) Follow(ctx context.Context, req *FollowRequest) (*FollowResponse, error) {
	newFollower := database.NewFollower(req.FolloweeUserId, req.FollowerUserId)
	err := database.InsertFollower(ctx, m.db, newFollower)
	if err != nil {
		return nil, err
	}
	return &FollowResponse{}, nil
}

func (m *Model) GetFollowers(context.Context, *GetFollowersOfUserRequest) (*GetFollowersOfUserResponse, error) {
	return nil, errors.Errorf("unimplemented")
}

func (m *Model) CreateUpvote(ctx context.Context, req *CreateUpvoteRequest) (*CreateUpvoteResponse, error) {
	newUpvote := database.NewUpvote(req.UserId, req.MessageId)
	err := database.InsertUpvote(ctx, m.db, newUpvote)
	if err != nil {
		return nil, err
	}
	return &CreateUpvoteResponse{UpvoteId: newUpvote.UpvoteId}, nil
}

func (m *Model) GetUserTimeline(context.Context, *GetUserTimelineRequest) (*GetUserTimelineResponse, error) {
	return nil, errors.Errorf("unimplemented")
}

func (m *Model) GetUserMessages(context.Context, *GetUserMessagesRequest) (*GetUserMessagesResponse, error) {
	return nil, errors.Errorf("unimplemented")
}

func (m *Model) SearchUsers(context.Context, *SearchUsersRequest) (*SearchUsersResponse, error) {
	return nil, errors.Errorf("unimplemented")
}

func (m *Model) SearchMessages(context.Context, *SearchMessagesRequest) (*SearchMessagesResponse, error) {
	return nil, errors.Errorf("unimplemented")
}
