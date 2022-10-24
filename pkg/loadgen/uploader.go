package loadgen

import (
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/mattfenwick/scaling/pkg/webserver"
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
