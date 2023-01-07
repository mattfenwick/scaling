package webserver

import "github.com/google/uuid"

// users

type CreateUserRequest struct {
	Name  string
	Email string
}

type CreateUserResponse struct {
	UserId uuid.UUID
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
	Users []GetUserResponse
}

type SearchUsersRequest struct {
	NamePattern  string
	EmailPattern string
	// TODO other stuff?
}

type SearchUsersResponse struct {
	Users []GetUserResponse
}

type GetUserMessagesRequest struct {
	UserId uuid.UUID
	// TODO paginate
}

type GetUserMessagesResponse struct {
	UserId   uuid.UUID
	Messages []GetMessageResponse
}

type GetUserTimelineRequest struct {
	UserId uuid.UUID
	// TODO paginate
}

type GetUserTimelineResponse struct {
	UserId   uuid.UUID
	Messages []GetMessageResponse
}

// messages

type CreateMessageRequest struct {
	SenderUserId uuid.UUID
	Content      string
}

type CreateMessageResponse struct {
	MessageId uuid.UUID
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
}

type SearchMessagesRequest struct {
}

type SearchMessagesResponse struct {
	Messages []GetMessageResponse
}

// follow/upvote

type FollowRequest struct {
	FolloweeUserId uuid.UUID
	FollowerUserId uuid.UUID
}

type FollowResponse struct {
}

type CreateUpvoteRequest struct {
	UserId    uuid.UUID
	MessageId uuid.UUID
}

type CreateUpvoteResponse struct {
	UpvoteId uuid.UUID
}

type GetFollowersOfUserRequest struct {
	UserId string
	// TODO paginate
}

type GetFollowersOfUserResponse struct {
	Followers []struct {
		UserId uuid.UUID
		Name   string
		Email  string
	}
}
