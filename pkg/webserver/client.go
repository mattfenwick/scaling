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

// users

func (c *Client) CreateUser(ctx context.Context, request *CreateUserRequest) (*CreateUserResponse, error) {
	out, _, err := utils.RestyIssueRequest[CreateUserResponse](ctx, c.Resty, "POST", UserPath, request, nil)
	return out, err
}

func (c *Client) GetUser(ctx context.Context, request *GetUserRequest) (*GetUserResponse, error) {
	params := map[string]string{"userid": request.UserId.String()}
	out, _, err := utils.RestyIssueRequest[GetUserResponse](ctx, c.Resty, "GET", UserPath, nil, params)
	return out, err
}

func (c *Client) GetUsers(ctx context.Context, request *GetUsersRequest) (*GetUsersResponse, error) {
	out, _, err := utils.RestyIssueRequest[GetUsersResponse](ctx, c.Resty, "GET", UsersPath, nil, nil)
	return out, err
}

func (c *Client) SearchUsers(ctx context.Context, request *SearchUsersRequest) (*SearchUsersResponse, error) {
	out, _, err := utils.RestyIssueRequest[SearchUsersResponse](ctx, c.Resty, "POST", UsersPath, request, nil)
	return out, err
}

func (c *Client) GetUserTimeline(ctx context.Context, request *GetUserTimelineRequest) (*GetUserTimelineResponse, error) {
	out, _, err := utils.RestyIssueRequest[GetUserTimelineResponse](ctx, c.Resty, "POST", UserTimelinePath, request, nil)
	return out, err
}

func (c *Client) GetUserMessages(ctx context.Context, request *GetUserMessagesRequest) (*GetUserMessagesResponse, error) {
	out, _, err := utils.RestyIssueRequest[GetUserMessagesResponse](ctx, c.Resty, "POST", UserMessagesPath, request, nil)
	return out, err
}

// messages

func (c *Client) CreateMessage(ctx context.Context, request *CreateMessageRequest) (*CreateMessageResponse, error) {
	out, _, err := utils.RestyIssueRequest[CreateMessageResponse](ctx, c.Resty, "POST", MessagePath, request, nil)
	return out, err
}

func (c *Client) GetMessage(ctx context.Context, request *GetMessageRequest) (*GetMessageResponse, error) {
	params := map[string]string{"messageid": request.MessageId.String()}
	out, _, err := utils.RestyIssueRequest[GetMessageResponse](ctx, c.Resty, "GET", MessagePath, nil, params)
	return out, err
}

func (c *Client) GetMessages(ctx context.Context, request *GetMessagesRequest) (*GetMessagesResponse, error) {
	out, _, err := utils.RestyIssueRequest[GetMessagesResponse](ctx, c.Resty, "GET", MessagesPath, request, nil)
	return out, err
}

func (c *Client) SearchMessages(ctx context.Context, request *SearchMessagesRequest) (*SearchMessagesResponse, error) {
	out, _, err := utils.RestyIssueRequest[SearchMessagesResponse](ctx, c.Resty, "POST", MessagesPath, request, nil)
	return out, err
}

// follow/upvote

func (c *Client) FollowUser(ctx context.Context, request *FollowRequest) (*FollowResponse, error) {
	out, _, err := utils.RestyIssueRequest[FollowResponse](ctx, c.Resty, "POST", FollowPath, request, nil)
	return out, err
}

func (c *Client) GetFollowers(ctx context.Context, request *GetFollowersOfUserRequest) (*GetFollowersOfUserResponse, error) {
	params := map[string]string{"userid": request.UserId.String()}
	out, _, err := utils.RestyIssueRequest[GetFollowersOfUserResponse](ctx, c.Resty, "GET", FollowersPath, nil, params)
	return out, err
}

func (c *Client) UpvoteMessage(ctx context.Context, request *CreateUpvoteRequest) (*CreateUpvoteResponse, error) {
	out, _, err := utils.RestyIssueRequest[CreateUpvoteResponse](ctx, c.Resty, "POST", UpvotePath, request, nil)
	return out, err
}
