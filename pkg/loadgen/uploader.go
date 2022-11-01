package loadgen

import (
	"context"
	"math/rand"
	"time"

	"github.com/mattfenwick/scaling/pkg/telemetry"
	"github.com/mattfenwick/scaling/pkg/webserver"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Uploader struct {
	Client *webserver.Client
}

func NewUploader(client *webserver.Client) *Uploader {
	return &Uploader{
		Client: client,
	}
}

func (u *Uploader) RunCannedUploads() error {
	ctx := context.Background()
	childCtx, childCancel := context.WithCancel(ctx)
	defer childCancel()
	users := GenerateUsers(childCtx, int(time.Now().Unix()))

	for i := 0; i < 10; i++ {

		nextUser := <-users
		resp, err := u.Client.CreateUser(&webserver.CreateUserRequest{Name: nextUser[0], Email: nextUser[1]})
		if err != nil {
			return err
		}

		logrus.Infof("created user of id: %s", resp.UserId)

		// TODO get user, add upvotes, messages, followers, etc.
		// u.Client.GetUser(resp.UserId)
	}
	return nil
}

func (u *Uploader) RunContinuous(ctx context.Context, keyCounts []int, workers int, pauseMilliseconds int) {
	for workerIdMutable := 0; workerIdMutable < workers; workerIdMutable++ {
		go func(workerId int) {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			for i := rand.Intn(len(keyCounts)); ; i++ {
				select {
				case <-ctx.Done():
					return
				default:
				}

				start := time.Now()
				// resp, err := u.Client.UploadDocument(&webserver.UploadDocumentRequest{
				// 	Document: json.MustMarshalToString(GenerateByNumberOfKeys(keys)),
				// })
				err := errors.Errorf("TODO")
				panic(err)

				telemetry.RecordClientApiRequestDuration("upload", err, start)
				if err == nil {

				} else {

				}

				time.Sleep(time.Duration(pauseMilliseconds) * time.Millisecond)
			}
		}(workerIdMutable)
	}
}
