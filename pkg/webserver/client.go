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
		Resty: resty.New().SetBaseURL(url).SetTransport(utils.OtelTransport()),
	}
}

func (c *Client) CreateUser(request *CreateUserRequest) (*CreateUserResponse, error) {
	out, _, err := utils.RestyIssueRequest[CreateUserResponse](c.Resty, "POST", UsersPath, request, nil)
	return out, err
}

func (c *Client) CreateMessage(request *CreateMessageRequest) (*CreateMessageResponse, error) {
	out, _, err := utils.RestyIssueRequest[CreateMessageResponse](c.Resty, "POST", MessagesPath, request, nil)
	return out, err
}

func (c *Client) CreateFollower(request *FollowRequest) (*FollowResponse, error) {
	out, _, err := utils.RestyIssueRequest[FollowResponse](c.Resty, "POST", FollowersPath, request, nil)
	return out, err
}

func (c *Client) CreateUpvote(request *CreateUpvoteRequest) (*CreateUpvoteResponse, error) {
	out, _, err := utils.RestyIssueRequest[CreateUpvoteResponse](c.Resty, "POST", UpvotesPath, request, nil)
	return out, err
}

func (c *Client) UploadDocument(request *UploadDocumentRequest) (*UploadDocumentResponse, error) {
	out, _, err := utils.RestyIssueRequest[UploadDocumentResponse](c.Resty, "POST", DocumentsPath, request, nil)
	return out, err
}

func (c *Client) GetDocument(request *GetDocumentRequest) (*GetDocumentResponse, error) {
	out, _, err := utils.RestyIssueRequest[GetDocumentResponse](c.Resty, "GET", DocumentsPath, nil, map[string]string{"id": request.DocumentId})
	return out, err
}

func (c *Client) GetAllDocuments() (*GetAllDocumentsResponse, error) {
	out, _, err := utils.RestyIssueRequest[GetAllDocumentsResponse](c.Resty, "GET", AllDocumentsPath, nil, nil)
	return out, err
}

func (c *Client) FindDocuments(request *FindDocumentsRequest) (*FindDocumentsResponse, error) {
	out, _, err := utils.RestyIssueRequest[FindDocumentsResponse](c.Resty, "POST", FindDocumentsPath, request, nil)
	return out, err
}
