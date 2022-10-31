package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type User struct {
	UserId    uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
}

func NewUser(name string, email string) *User {
	return &User{UserId: uuid.New(), Name: name, Email: email, CreatedAt: time.Now()}
}

func InsertUser(ctx context.Context, db *sql.DB, user *User) error {
	_, err := db.ExecContext(ctx,
		"INSERT INTO users (user_id, name, email, created_at) VALUES ($1, $2, $3, $4)",
		user.UserId,
		user.Name,
		user.Email,
		user.CreatedAt,
	)
	return errors.Wrapf(err, "unable to insert user")
}

func ReadAllUsers(ctx context.Context, db *sql.DB) ([]*User, error) {
	rows, err := db.QueryContext(ctx, "select * from users")
	if err != nil {
		return nil, errors.Wrapf(err, "unable to issue query")
	}

	var records []*User

	for rows.Next() {
		var record User
		err = rows.Scan(&record.UserId, &record.Name, &record.Email, &record.CreatedAt)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load row")
		}
		records = append(records, &record)
	}

	if closeErr := rows.Close(); closeErr != nil {
		return nil, errors.Wrapf(err, "unable to close")
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "row iteration problem")
	}

	return records, nil
}

type Follower struct {
	FolloweeUserId uuid.UUID
	FollowerUserId uuid.UUID
	CreatedAt      time.Time
}

func NewFollower(followeeId uuid.UUID, followerId uuid.UUID) *Follower {
	return &Follower{FolloweeUserId: followeeId, FollowerUserId: followerId, CreatedAt: time.Now()}
}

func InsertFollower(ctx context.Context, db *sql.DB, follower *Follower) error {
	_, err := db.ExecContext(ctx,
		"INSERT INTO followers (followee_user_id, follower_user_id, created_at) VALUES ($1, $2, $3)",
		follower.FolloweeUserId,
		follower.FollowerUserId,
		follower.CreatedAt,
	)
	return errors.Wrapf(err, "unable to insert follower")
}

func ReadAllFollowers(ctx context.Context, db *sql.DB) ([]*Follower, error) {
	rows, err := db.QueryContext(ctx, "select * from followers")
	if err != nil {
		return nil, errors.Wrapf(err, "unable to issue query")
	}

	var records []*Follower

	for rows.Next() {
		var record Follower
		err = rows.Scan(&record.FolloweeUserId, &record.FollowerUserId, &record.CreatedAt)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load row")
		}
		records = append(records, &record)
	}

	if closeErr := rows.Close(); closeErr != nil {
		return nil, errors.Wrapf(err, "unable to close")
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "row iteration problem")
	}

	return records, nil
}

type Message struct {
	MessageId    uuid.UUID
	SenderUserId uuid.UUID
	Content      string
	CreatedAt    time.Time
}

func NewMessage(senderUserId uuid.UUID, content string) *Message {
	return &Message{MessageId: uuid.New(), SenderUserId: senderUserId, Content: content, CreatedAt: time.Now()}
}

func InsertMessage(ctx context.Context, db *sql.DB, message *Message) error {
	_, err := db.ExecContext(ctx,
		"INSERT INTO messages (message_id, sender_user_id, content, created_at) VALUES ($1, $2, $3, $4)",
		message.MessageId,
		message.SenderUserId,
		message.Content,
		message.CreatedAt,
	)
	return errors.Wrapf(err, "unable to insert follower")
}

func ReadAllMessages(ctx context.Context, db *sql.DB) ([]*Message, error) {
	rows, err := db.QueryContext(ctx, "select * from messages")
	if err != nil {
		return nil, errors.Wrapf(err, "unable to issue query")
	}

	var records []*Message

	for rows.Next() {
		var record Message
		err = rows.Scan(&record.MessageId, &record.SenderUserId, &record.Content, &record.CreatedAt)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load row")
		}
		records = append(records, &record)
	}

	if closeErr := rows.Close(); closeErr != nil {
		return nil, errors.Wrapf(err, "unable to close")
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "row iteration problem")
	}

	return records, nil
}

type Upvote struct {
	UpvoteId  uuid.UUID
	UserId    uuid.UUID
	MessageId uuid.UUID
	CreatedAt time.Time
}

func NewUpvote(userId uuid.UUID, messageId uuid.UUID) *Upvote {
	return &Upvote{UpvoteId: uuid.New(), UserId: userId, MessageId: messageId, CreatedAt: time.Now()}
}

func InsertUpvote(ctx context.Context, db *sql.DB, upvote *Upvote) error {
	_, err := db.ExecContext(ctx,
		"INSERT INTO upvotes (upvote_id, user_id, message_id, created_at) VALUES ($1, $2, $3, $4)",
		upvote.UpvoteId,
		upvote.UserId,
		upvote.MessageId,
		upvote.CreatedAt,
	)
	return errors.Wrapf(err, "unable to insert follower")
}

func ReadAllUpvotes(ctx context.Context, db *sql.DB) ([]*Upvote, error) {
	rows, err := db.QueryContext(ctx, "select * from upvotes")
	if err != nil {
		return nil, errors.Wrapf(err, "unable to issue query")
	}

	var records []*Upvote

	for rows.Next() {
		var record Upvote
		err = rows.Scan(&record.UpvoteId, &record.UserId, &record.MessageId, &record.CreatedAt)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load row")
		}
		records = append(records, &record)
	}

	if closeErr := rows.Close(); closeErr != nil {
		return nil, errors.Wrapf(err, "unable to close")
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "row iteration problem")
	}

	return records, nil
}
