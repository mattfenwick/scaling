package webserver

import "github.com/google/uuid"

type CreateUserRequest struct {
	Name  string
	Email string
}

type CreateUserResponse struct {
	UserId uuid.UUID
}

type FollowRequest struct {
	FolloweeUserId uuid.UUID
	FollowerUserId uuid.UUID
}

type FollowResponse struct {
}

type CreateMessageRequest struct {
	SenderUserId uuid.UUID
	Content      string
}

type CreateMessageResponse struct {
	MessageId uuid.UUID
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

type GetMessagesForUserRequest struct {
	UserId string
	// TODO paginate
}

type GetMessagesForUserResponse struct {
	Messages []struct {
		MessageId    uuid.UUID
		SenderUserId uuid.UUID
		Content      string
		UpvoteCount  int
	}
}

// create document

type UploadDocumentRequest struct {
	Document string
}

type UploadDocumentResponse struct {
	DocumentId string
}

// fetch single document

type GetDocumentRequest struct {
	DocumentId string
}

type GetDocumentResponse struct {
	Document *Document
}

// fetch all documents

type GetAllDocumentsResponse struct {
	Documents map[string]*Document
}

// find documents

type FindDocumentsRequest struct {
	Key string
}

type FindDocumentsResponseItem struct {
	DocumentId string
	Paths      [][]any
}

type FindDocumentsResponse struct {
	Matches []*FindDocumentsResponseItem
}
