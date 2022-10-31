package database

// import (
// 	"context"
// 	"database/sql"

// 	"github.com/google/uuid"
// 	"github.com/pkg/errors"
// )

// type UserFollowersView struct {
// 	User      *User
// 	Followers []*Follower
// }

// func ReadUserFollowersView(ctx context.Context, db *sql.DB, userId uuid.UUID) (*UserFollowersView, error) {
// 	user, err := ReadUserById(ctx, db, userId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	followers, err := ReadFollowersOf(ctx, db, userId)
// 	return &UserFollowersView{User: user, Followers: followers}, nil
// }

// func ReadUserFollowersOf(ctx context.Context, db *sql.DB, userId uuid.UUID) ([]*User, error) {
// 	user, err := ReadUserById(ctx, db, userId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	followers, err := ReadFollowersOf(ctx, db, userId)
// 	return &UserFollowersView{User: user, Followers: followers}, nil
// }

// type UserMessageView struct {
// 	UserId uuid.UUID
// }

// func ReadUserMessageView(ctx context.Context, db *sql.DB, userId uuid.UUID) (*UserMessageView, error) {
// 	// start with user id
// 	// find followers
// 	// find messages from all followers
// 	// also find messages from user
// 	// end up with sorted:
// 	// - messages, by date
// 	// - each message includes:
// 	//   author
// 	//   content
// 	//   upvotes (count, and list of people who upvoted)
// 	return nil, errors.Errorf("TODO")
// }
