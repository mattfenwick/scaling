package main

import (
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/mattfenwick/scaling/pkg/webserver"
	"github.com/sirupsen/logrus"
)

func main() {
	client := webserver.NewClient("http://localhost:8765")

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
		utils.DoOrDie(err)
		id := resp.DocumentId
		logrus.Infof("resp: %s", json.MustMarshalToString(resp))
		logrus.Infof("id: %s", id)

		allDocs, err := client.UnsafeGetAllDocuments()
		utils.DoOrDie(err)
		logrus.Infof("all docs: %d\n[%s]\n", len(allDocs.Documents), json.MustMarshalToString(allDocs))
	}

}
