package webserver

import "github.com/google/uuid"

// users

type CreateUserRequest struct {
	Name  string
	Email string
}

type CreateUserResponse struct {
	UserId  uuid.UUID
	Request *CreateUserRequest
}

type GetUserRequest struct {
	UserId uuid.UUID
}

type GetUserResponse struct {
	UserId uuid.UUID
	Name   string
	Email  string
}

type GetUsersRequest struct {
	// TODO limit?  paginate?
}

type GetUsersResponse struct {
	Users   []GetUserResponse
	Request *GetUsersRequest
}

type SearchUsersRequest struct {
	NamePattern  string
	EmailPattern string
	// TODO other stuff?
}

type SearchUsersResponse struct {
	Users   []GetUserResponse
	Request *SearchUsersRequest
}

type GetUserMessagesRequest struct {
	UserId uuid.UUID
	// TODO paginate
}

type GetUserMessagesResponse struct {
	UserId   uuid.UUID
	Messages []GetMessageResponse
	Request  *GetUserMessagesRequest
}

type GetUserTimelineRequest struct {
	UserId uuid.UUID
	// TODO paginate
}

type GetUserTimelineResponse struct {
	UserId   uuid.UUID
	Messages []GetMessageResponse
	Request  *GetUserTimelineRequest
}

// messages

type CreateMessageRequest struct {
	SenderUserId uuid.UUID
	Content      string
}

type CreateMessageResponse struct {
	MessageId uuid.UUID
	Request   *CreateMessageRequest
}

type GetMessageRequest struct {
	MessageId uuid.UUID
}

type GetMessageResponse struct {
	MessageId    uuid.UUID
	SenderUserId uuid.UUID
	Content      string
	UpvoteCount  int
}

type GetMessagesRequest struct {
}

type GetMessagesResponse struct {
	Messages []GetMessageResponse
	Request  *GetMessagesRequest
}

type SearchMessagesRequest struct {
	LiteralString string
}

type SearchMessagesResponse struct {
	Messages []GetMessageResponse
	Request  *SearchMessagesRequest
}

// follow/upvote

type FollowRequest struct {
	FolloweeUserId uuid.UUID
	FollowerUserId uuid.UUID
}

type FollowResponse struct {
	Request *FollowRequest
}

type CreateUpvoteRequest struct {
	UserId    uuid.UUID
	MessageId uuid.UUID
}

type CreateUpvoteResponse struct {
	UpvoteId uuid.UUID
	Request  *CreateUpvoteRequest
}

type GetFollowersOfUserRequest struct {
	UserId uuid.UUID
	// TODO paginate
}

type GetFollowersOfUserResponse struct {
	Followers []GetUserResponse
	Request   *GetFollowersOfUserRequest
}
