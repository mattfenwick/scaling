package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	readFollowersOfQueryTemplate = `
	select 
		users.*
	from followers
	inner join 
		users
	on 
		followers.follower_user_id = users.user_id 
	where 
		followers.followee_user_id = $1`

	readMessagesFromUserAndFollowersTemplate = `
	with userids as (
		select 
			$1 as user_id
		union
		select 
			follower_user_id
		from followers
		where followee_user_id = $1
	  ),
	  upvote_counts as (
		select message_id, count(*) as upvotes
		from messages
		group by message_id
	  )
	select
		messages.message_id,
		messages.sender_user_id,
		messages.content,
		upvote_counts.upvotes,
		messages.created_at
	from messages
	inner join 
		userids
	on 
		messages.sender_user_id = userids.user_id
	left join
	    upvote_counts
	on
	  	messages.message_id = upvote_counts.message_id`
)

var (
	loadUser = func(rows *sql.Rows, record *User) error {
		return rows.Scan(&record.UserId, &record.Name, &record.Email, &record.CreatedAt)
	}
	loadSingleUser = func(rows *sql.Row, record *User) error {
		return rows.Scan(&record.UserId, &record.Name, &record.Email, &record.CreatedAt)
	}

	loadMessage = func(rows *sql.Rows, record *Message) error {
		return rows.Scan(&record.MessageId, &record.SenderUserId, &record.Content, &record.CreatedAt)
	}

	loadMessageUpvoteCount = func(rows *sql.Rows, record *MessageUpvoteCount) error {
		return rows.Scan(&record.MessageId, &record.SenderUserId, &record.Content, &record.UpvoteCount, &record.CreatedAt)
	}
)

// Users

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

func GetUser(ctx context.Context, db *sql.DB, userId uuid.UUID) (*User, error) {
	// TODO consider using a prepared statement
	//   https://go.dev/doc/database/prepared-statements
	return ReadSingle(ctx, db, loadSingleUser, `SELECT * FROM users WHERE user_id = $1`, userId.String())
}

func ReadAllUsers(ctx context.Context, db *sql.DB) ([]*User, error) {
	return ReadMany(ctx, db, loadUser, "select * from users")
}

func ReadUserById(ctx context.Context, db *sql.DB, userId uuid.UUID) (*User, error) {
	rows := db.QueryRowContext(ctx, "select * from users where user_id = $1", userId)

	var record User
	err := rows.Scan(&record.UserId, &record.Name, &record.Email, &record.CreatedAt)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load row")
	}
	return &record, nil
}

// Followers

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
	process := func(rows *sql.Rows, record *Follower) error {
		return rows.Scan(&record.FolloweeUserId, &record.FollowerUserId, &record.CreatedAt)
	}
	return ReadMany(ctx, db, process, "select * from followers")
}

func ReadFollowersOf(ctx context.Context, db *sql.DB, userId uuid.UUID) ([]*User, error) {
	return ReadMany(ctx, db, loadUser, readFollowersOfQueryTemplate, userId)
}

// Messages

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
	return ReadMany(ctx, db, loadMessage, "select * from messages")
}

type MessageUpvoteCount struct {
	MessageId    uuid.UUID
	SenderUserId uuid.UUID
	Content      string
	UpvoteCount  int
	CreatedAt    time.Time
}

func ReadMessagesForUser(ctx context.Context, db *sql.DB, userId uuid.UUID) ([]*MessageUpvoteCount, error) {
	return ReadMany(ctx, db, loadMessageUpvoteCount, readMessagesFromUserAndFollowersTemplate, userId)
}

// Upvotes

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
	process := func(rows *sql.Rows, record *Upvote) error {
		return rows.Scan(&record.UpvoteId, &record.UserId, &record.MessageId, &record.CreatedAt)
	}
	return ReadMany(ctx, db, process, "select * from upvotes")
}
