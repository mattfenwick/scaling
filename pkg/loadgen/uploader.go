package loadgen

import (
	"context"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/telemetry"
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/mattfenwick/scaling/pkg/webserver"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
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
	docs := []string{
		"abcdef",
		`{"abc": [123, 456]}`,
		"{}",
		`{"": [true, false, null, 0, -1, "abc", [], {}]}`,
		"1",
		"null",
		"34",
		`"qrs"`,
		"true",
		"false",
	}
	for _, doc := range docs {
		resp, err := u.Client.UploadDocument(&webserver.UploadDocumentRequest{
			Document: doc,
		})
		if err != nil {
			return err
		}

		id := resp.DocumentId
		logrus.Infof("resp: %s", json.MustMarshalToString(resp))
		logrus.Infof("id: %s", id)

		fetchedDoc, err := u.Client.GetDocument(&webserver.GetDocumentRequest{
			DocumentId: id,
		})
		if err != nil {
			return err
		}

		logrus.Infof("fetched doc: %s", json.MustMarshalToString(fetchedDoc.Document))

		allDocs, err := u.Client.GetAllDocuments()
		if err != nil {
			return err
		}
		//logrus.Infof("all docs: %d\n", len(allDocs.Documents))
		logrus.Infof("all docs: %d\n[%s]\n", len(allDocs.Documents), json.MustMarshalToString(allDocs))
	}
	return nil
}

func (u *Uploader) RunRandomUploadsByKeyCount(keyCounts []int) {
	// TODO concurrency?
	for _, c := range keyCounts {
		resp, err := u.Client.UploadDocument(&webserver.UploadDocumentRequest{
			Document: json.MustMarshalToString(GenerateByNumberOfKeys(c)),
		})
		utils.DoOrDie(err)

		logrus.Debugf("resp: %s", json.MustMarshalToString(resp))
	}
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

				keys := keyCounts[i%len(keyCounts)]
				start := time.Now()
				resp, err := u.Client.UploadDocument(&webserver.UploadDocumentRequest{
					Document: json.MustMarshalToString(GenerateByNumberOfKeys(keys)),
				})
				telemetry.RecordClientApiRequestDuration("upload", err, start)
				if err == nil {
					logrus.Infof("worker %d generated document %d, key count %d: %s", workerId, i, keys, resp.DocumentId)
				} else {
					logrus.Errorf("worker %d unable to generate document %d, key count %d: %s", workerId, i, keys, err.Error())
				}

				time.Sleep(time.Duration(pauseMilliseconds) * time.Millisecond)
			}
		}(workerIdMutable)
	}
}
