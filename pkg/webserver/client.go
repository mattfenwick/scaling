package webserver

import (
	"context"

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

func (c *Client) CreateUser(ctx context.Context, request *CreateUserRequest) (*CreateUserResponse, error) {
	out, _, err := utils.RestyIssueRequest[CreateUserResponse](ctx, c.Resty, "POST", UsersPath, request, nil)
	return out, err
}

func (c *Client) CreateMessage(ctx context.Context, request *CreateMessageRequest) (*CreateMessageResponse, error) {
	out, _, err := utils.RestyIssueRequest[CreateMessageResponse](ctx, c.Resty, "POST", MessagesPath, request, nil)
	return out, err
}

func (c *Client) CreateFollower(ctx context.Context, request *FollowRequest) (*FollowResponse, error) {
	out, _, err := utils.RestyIssueRequest[FollowResponse](ctx, c.Resty, "POST", FollowersPath, request, nil)
	return out, err
}

func (c *Client) CreateUpvote(ctx context.Context, request *CreateUpvoteRequest) (*CreateUpvoteResponse, error) {
	out, _, err := utils.RestyIssueRequest[CreateUpvoteResponse](ctx, c.Resty, "POST", UpvotesPath, request, nil)
	return out, err
}

func (c *Client) UploadDocument(ctx context.Context, request *UploadDocumentRequest) (*UploadDocumentResponse, error) {
	out, _, err := utils.RestyIssueRequest[UploadDocumentResponse](ctx, c.Resty, "POST", DocumentsPath, request, nil)
	return out, err
}

func (c *Client) GetDocument(ctx context.Context, request *GetDocumentRequest) (*GetDocumentResponse, error) {
	out, _, err := utils.RestyIssueRequest[GetDocumentResponse](ctx, c.Resty, "GET", DocumentsPath, nil, map[string]string{"id": request.DocumentId})
	return out, err
}

func (c *Client) GetAllDocuments(ctx context.Context) (*GetAllDocumentsResponse, error) {
	out, _, err := utils.RestyIssueRequest[GetAllDocumentsResponse](ctx, c.Resty, "GET", AllDocumentsPath, nil, nil)
	return out, err
}

func (c *Client) FindDocuments(ctx context.Context, request *FindDocumentsRequest) (*FindDocumentsResponse, error) {
	out, _, err := utils.RestyIssueRequest[FindDocumentsResponse](ctx, c.Resty, "POST", FindDocumentsPath, request, nil)
	return out, err
}
