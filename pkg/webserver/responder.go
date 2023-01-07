package webserver

import "context"

/*

/user
 - POST => create user
 - GET => get user by uuid
/user/timeline
 - GET => by user's uuid
/user/messages
 - GET => by user's uuid

/users
 - POST => fancy search
 - GET => naive list of users

/message
 - POST => create
 - GET => by message's uuid
/messages
 - POST => fancy search (TODO)
 - GET => naive list of messages

/follow
 - POST
/followers
 - GET => by user's uuid (TODO)

/upvote
 - POST (TODO)

*/

type Responder interface {
	Sleep(ctx context.Context, seconds string) error

	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	GetUserTimeline(context.Context, *GetUserTimelineRequest) (*GetUserTimelineResponse, error) // TODO paginate
	GetUserMessages(context.Context, *GetUserMessagesRequest) (*GetUserMessagesResponse, error) // TODO paginate
	GetUsers(context.Context, *GetUsersRequest) (*GetUsersResponse, error)                      // TODO paginate
	SearchUsers(context.Context, *SearchUsersRequest) (*SearchUsersResponse, error)             // TODO paginate

	CreateMessage(context.Context, *CreateMessageRequest) (*CreateMessageResponse, error)
	GetMessage(context.Context, *GetMessageRequest) (*GetMessageResponse, error)
	GetMessages(context.Context, *GetMessagesRequest) (*GetMessagesResponse, error)          // TODO pagniate
	SearchMessages(context.Context, *SearchMessagesRequest) (*SearchMessagesResponse, error) // TODO paginate

	Follow(context.Context, *FollowRequest) (*FollowResponse, error)
	GetFollowers(context.Context, *GetFollowersOfUserRequest) (*GetFollowersOfUserResponse, error)
	CreateUpvote(context.Context, *CreateUpvoteRequest) (*CreateUpvoteResponse, error)

	IsLive(context.Context) bool
	IsReady(context.Context) bool

	Dump(ctx context.Context) (string, error)
}
