package loadgen

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

var (
	MessageOpen = []string{
		"aloha",
		"bonjour",
		"cool news today",
		"good morning",
		"hey friends",
	}
	MessageBody = []string{
		"apples taste good!",
		"bananas are yellow",
		"cedar trees grow large",
		"darkness happens at night",
		"eureka, I wrote a readme",
	}
	MessageClose = []string{
		"au revoir",
		"bye",
		"check ya later",
		"farewell",
		"goodbye",
		"hasta la vista",
		"later, taters",
		"so long",
	}
)

type MessageState struct {
	Open      int
	Middle    int
	Close     int
	Iteration int
	Stamp     int
}

func (m *MessageState) GetMessage() string {
	first, middle, last := MessageOpen[m.Open], MessageBody[m.Middle], MessageClose[m.Close]
	return fmt.Sprintf("%s %s %s (%d, %d)", first, middle, last, m.Iteration, m.Stamp)
}

func (m *MessageState) Increment() {
	logrus.Infof("message state before: %d , %d , %d , %d ", m.Close, m.Middle, m.Open, m.Iteration)
	defer func() {
		logrus.Infof("message state  after: %d , %d , %d , %d ", m.Close, m.Middle, m.Open, m.Iteration)
	}()
	m.Close++
	if m.Close < len(MessageClose) {
		return
	}

	m.Close = 0
	m.Middle++
	if m.Middle < len(MessageBody) {
		return
	}

	m.Middle = 0
	m.Open++
	if m.Open < len(MessageOpen) {
		return
	}

	m.Open = 0
	m.Iteration++
}

func GenerateMessages(ctx context.Context, stamp int) <-chan string {
	out := make(chan string)
	go func() {
		state := &MessageState{Stamp: stamp}
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			out <- state.GetMessage()
			state.Increment()
		}
	}()
	return out
}
