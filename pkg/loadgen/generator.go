package loadgen

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/mattfenwick/scaling/pkg/telemetry"
	"github.com/mattfenwick/scaling/pkg/webserver"
	"github.com/sirupsen/logrus"
)

type Generator struct {
	Client  *webserver.Client
	Actions chan func()
	UserIds []uuid.UUID
}

func NewGenerator(ctx context.Context, client *webserver.Client) *Generator {
	g := &Generator{
		Client:  client,
		Actions: make(chan func()),
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case action := <-g.Actions:
				action()
			}
		}
	}()
	return g
}

func (g *Generator) CreateUsers(ctx context.Context, limit int) {
	childCtx, childCancel := context.WithCancel(ctx)
	defer childCancel()

	stamp := int(time.Now().Unix())
	users := GenerateUsers(childCtx, stamp)
	messages := GenerateMessages(childCtx, stamp)

	for i := 0; i < limit; i++ {
		nextUser := <-users
		start := time.Now()
		resp, err := g.Client.CreateUser(childCtx, &webserver.CreateUserRequest{Name: nextUser[0], Email: nextUser[1]})
		telemetry.RecordClientApiRequestDuration("create user", err, start)
		if err != nil {
			logrus.Errorf("unable to create user: %+v", err)
		} else {
			logrus.Infof("created user of id: %s", resp.UserId)
			g.Actions <- func() {
				g.UserIds = append(g.UserIds, resp.UserId)
			}
			g.CreateMessages(childCtx, resp.UserId, messages, rand.Intn(100))
		}
	}
}

func (g *Generator) CreateMessages(ctx context.Context, userId uuid.UUID, messages <-chan string, count int) {
	for i := 0; i < count; i++ {
		logrus.Infof("creating message %d of %d for user %s", i+1, count, userId.String())
		resp, err := g.Client.CreateMessage(ctx, &webserver.CreateMessageRequest{SenderUserId: userId, Content: <-messages})
		if err != nil {
			logrus.Errorf("unable to create message: %+v", err)
		} else {
			logrus.Infof("created message of id: %s", resp.MessageId)
		}
	}
}
