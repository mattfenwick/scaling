package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	getFollowersOfQueryTemplate = `
	select 
		users.*
	from followers
	inner join 
		users
	on 
		followers.follower_user_id = users.user_id 
	where 
		followers.followee_user_id = $1`

	getUserMessagesTemplate = `
	with upvote_counts as (
		select message_id, count(*) as upvotes
		from upvotes
		group by message_id
	  )
	select
		messages.message_id,
		messages.sender_user_id,
		messages.content,
		coalesce(upvote_counts.upvotes, 0),
		messages.created_at
	from messages
	left join
	    upvote_counts
	on
	  	messages.message_id = upvote_counts.message_id
	where
		messages.sender_user_id = $1`

	getUserTimelineTemplate = `
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
		from upvotes
		group by message_id
	  )
	select
		messages.message_id,
		messages.sender_user_id,
		messages.content,
		coalesce(upvote_counts.upvotes, 0),
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
	loadSingleMessage = func(rows *sql.Row, record *Message) error {
		return rows.Scan(&record.MessageId, &record.SenderUserId, &record.Content, &record.CreatedAt)
	}

	loadTimelineMessage = func(rows *sql.Rows, record *TimelineMessage) error {
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

func GetUsers(ctx context.Context, db *sql.DB) ([]*User, error) {
	return ReadMany(ctx, db, loadUser, "select * from users")
}

func regexWrap(s string) string {
	return "%" + s + "%"
}

func SearchUsers(ctx context.Context, db *sql.DB, namePattern string, emailPattern string) ([]*User, error) {
	return ReadMany(ctx, db, loadUser,
		"select * from users where name ilike $1 and email ilike $2",
		regexWrap(namePattern),
		regexWrap(emailPattern))
}

type TimelineMessage struct {
	MessageId    uuid.UUID
	SenderUserId uuid.UUID
	Content      string
	UpvoteCount  int
	CreatedAt    time.Time
}

func GetUserTimeline(ctx context.Context, db *sql.DB, userId uuid.UUID) ([]*TimelineMessage, error) {
	return ReadMany(ctx, db, loadTimelineMessage, getUserTimelineTemplate, userId)
}

func GetUserMessages(ctx context.Context, db *sql.DB, userId uuid.UUID) ([]*TimelineMessage, error) {
	return ReadMany(ctx, db, loadTimelineMessage, getUserMessagesTemplate, userId)
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

func GetMessage(ctx context.Context, db *sql.DB, messageId uuid.UUID) (*Message, error) {
	return ReadSingle(ctx, db, loadSingleMessage, `SELECT * FROM messages WHERE message_id = $1`, messageId.String())
}

func GetMessages(ctx context.Context, db *sql.DB) ([]*Message, error) {
	return ReadMany(ctx, db, loadMessage, "select * from messages")
}

func SearchMessages(ctx context.Context, db *sql.DB, literalString string) ([]*Message, error) {
	return ReadMany(ctx, db, loadMessage,
		"select * from messages where position($1 in content) > 0",
		literalString)
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

func GetFollowers(ctx context.Context, db *sql.DB) ([]*Follower, error) {
	process := func(rows *sql.Rows, record *Follower) error {
		return rows.Scan(&record.FolloweeUserId, &record.FollowerUserId, &record.CreatedAt)
	}
	return ReadMany(ctx, db, process, "select * from followers")
}

func GetFollowersOfUser(ctx context.Context, db *sql.DB, userId uuid.UUID) ([]*User, error) {
	return ReadMany(ctx, db, loadUser, getFollowersOfQueryTemplate, userId)
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

// debug

func GetTableSizes(ctx context.Context, db *sql.DB) (map[string]int, error) {
	tableNames := []string{
		"users",
		"messages",
		"followers",
		"upvotes",
		"topics",
		"pings",
	}
	process := func(row *sql.Row, out *int) error {
		return errors.Wrapf(row.Scan(out), "unable to fetch row")
	}
	rowCounts := map[string]int{}
	for _, table := range tableNames {
		count, err := ReadSingle(ctx, db, process, fmt.Sprintf(`SELECT count(*) FROM %s`, table))
		if err != nil {
			return nil, err
		}
		if count == nil {
			return nil, errors.Errorf("unable to get row size for table %s", table)
		}
		rowCounts[table] = *count
	}
	return rowCounts, nil
}
