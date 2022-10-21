package webserver

import (
	"github.com/go-resty/resty/v2"
	"github.com/mattfenwick/scaling/pkg/utils"
)

type Client struct {
	URL   string
	Resty *resty.Client
}

func NewClient(url string) *Client {
	return &Client{
		URL:   url,
		Resty: resty.New().SetBaseURL(url),
	}
}

func (c *Client) UploadDocument(request *UploadDocumentRequest) (*UploadDocumentResponse, error) {
	out, _, err := utils.RestyIssueRequest[UploadDocumentResponse](c.Resty, "POST", DocumentsPath, request, nil)
	return out, err
}

func (c *Client) GetDocument(request *GetDocumentRequest) (*GetDocumentResponse, error) {
	out, _, err := utils.RestyIssueRequest[GetDocumentResponse](c.Resty, "GET", DocumentsPath, nil, map[string]string{"documentId": request.DocumentId})
	return out, err
}

func (c *Client) UnsafeGetAllDocuments() (*UnsafeGetDocumentsResponse, error) {
	out, _, err := utils.RestyIssueRequest[UnsafeGetDocumentsResponse](c.Resty, "GET", UnsafeDocumentsPath, nil, nil)
	return out, err
}
