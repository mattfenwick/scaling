package client

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/webserver"
	"github.com/sirupsen/logrus"
)

func RunSmallBatchOfRequests(host string, port int) error {
	url := fmt.Sprintf("http://%s:%d", host, port)
	client := webserver.NewClient(url)

	docs := []string{
		"abcdef",
		"123456",
		"{}",
		"[]",
		"1",
		"null",
		"34",
		`"qrs"`,
		"true",
		"false",
	}
	for _, doc := range docs {
		resp, err := client.UploadDocument(&webserver.UploadDocumentRequest{
			Document: doc,
		})
		if err != nil {
			return err
		}

		id := resp.DocumentId
		logrus.Infof("resp: %s", json.MustMarshalToString(resp))
		logrus.Infof("id: %s", id)

		fetchedDoc, err := client.GetDocument(&webserver.GetDocumentRequest{
			DocumentId: id,
		})
		if err != nil {
			return err
		}

		logrus.Infof("fetched doc: %s", json.MustMarshalToString(fetchedDoc.Document))

		allDocs, err := client.UnsafeGetAllDocuments()
		if err != nil {
			return err
		}
		//logrus.Infof("all docs: %d\n", len(allDocs.Documents))
		logrus.Infof("all docs: %d\n[%s]\n", len(allDocs.Documents), json.MustMarshalToString(allDocs))
	}
	return nil
}
