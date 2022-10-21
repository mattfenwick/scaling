package webserver

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	out, _, err := IssueRequest[UploadDocumentResponse](c.Resty, "POST", "/documents", request, nil)
	return out, err
}

func (c *Client) GetDocument(request *GetDocumentRequest) (*GetDocumentResponse, error) {
	out, _, err := IssueRequest[GetDocumentResponse](c.Resty, "GET", "/documents", nil, map[string]string{"documentId": request.DocumentId})
	return out, err
}

func IssueRequest[A any](restyClient *resty.Client, verb string, path string, body interface{}, queryParams map[string]string) (*A, string, error) {
	var err error
	request := restyClient.R()
	if body != nil {
		reqBody := json.MustMarshalToString(body)
		log.Tracef("request body: %s", reqBody)
		request = request.SetBody(body)
	}

	request = request.SetQueryParams(queryParams)

	urlPath := fmt.Sprintf("%s/%s", restyClient.BaseURL, path)
	log.Debugf("issuing %s to %s", verb, urlPath)

	var resp *resty.Response
	switch verb {
	case "GET":
		resp, err = request.Get(path)
	case "POST":
		resp, err = request.Post(path)
	case "PUT":
		resp, err = request.Put(path)
	case "DELETE":
		resp, err = request.Delete(path)
	default:
		return nil, "", errors.Errorf("unrecognized http verb %s to %s", verb, path)
	}
	if err != nil {
		return nil, "", errors.Wrapf(err, "unable to issue %s to %s", verb, path)
	}

	respBody, statusCode := resp.String(), resp.StatusCode()
	log.Debugf("response code %d from %s to %s", statusCode, verb, urlPath)
	log.Tracef("response body: %s", respBody)

	if !resp.IsSuccess() {
		return nil, respBody, errors.Errorf("bad status code for %s to path %s: %d, response %s", verb, path, statusCode, respBody)
	}

	out, err := json.ParseString[A](respBody)
	return out, respBody, err
}
