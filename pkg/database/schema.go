package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	uuidOsspExtention = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`

	usersTable = `
CREATE TABLE IF NOT EXISTS users (
    user_id uuid NOT NULL, -- DEFAULT uuid_generate_v4() NOT NULL,
    name varchar(80) NOT NULL,
    email varchar(80) NOT NULL,
    created_at timestamp NOT NULL, -- DEFAULT NOW() NOT NULL,
    CONSTRAINT users_pk PRIMARY KEY (user_id)
);
`

	followersTable = `
CREATE TABLE IF NOT EXISTS followers (
    followee_user_id uuid NOT NULL references users(user_id),
    follower_user_id uuid NOT NULL references users(user_id),
    created_at timestamp NOT NULL, -- DEFAULT NOW() NOT NULL,
    CONSTRAINT followers_pk PRIMARY KEY (followee_user_id, follower_user_id)
);
`

	messagesTable = `
CREATE TABLE IF NOT EXISTS messages (
    message_id uuid NOT NULL, -- DEFAULT uuid_generate_v4() NOT NULL,
    sender_user_id uuid NOT NULL references users(user_id),
    content varchar(200) NOT NULL,
    created_at timestamp NOT NULL, -- DEFAULT NOW() NOT NULL,
    CONSTRAINT messages_pk PRIMARY KEY (message_id)
);
`

	upvotesTable = `
CREATE TABLE IF NOT EXISTS upvotes (
    upvote_id uuid NOT NULL, -- DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL references users(user_id),
    message_id uuid NOT NULL references messages(message_id),
    created_at timestamp NOT NULL, -- DEFAULT NOW() NOT NULL,
    CONSTRAINT upvotes_pk PRIMARY KEY (upvote_id)
);
`

	// ?? derived tables ??

	topicsTable = `
CREATE TABLE IF NOT EXISTS topics (
   topic_id uuid NOT NULL, -- DEFAULT uuid_generate_v4() NOT NULL,
   name varchar(80) NOT NULL,
   description varchar(80) NOT NULL,
--   created_at timestamp DEFAULT NOW() NOT NULL,
   CONSTRAINT topics_pk PRIMARY KEY (topic_id)
);
`

	pingsTable = `
CREATE TABLE IF NOT EXISTS pings (
   ping_id uuid NOT NULL, -- DEFAULT uuid_generate_v4() NOT NULL,
   user_id uuid NOT NULL references users(user_id),
   message_id uuid NOT NULL references messages(message_id),
--   created_at timestamp DEFAULT NOW() NOT NULL,
   CONSTRAINT pings_pk PRIMARY KEY (ping_id)
);
`
)

func DoesDatabaseExist(ctx context.Context, db *sql.DB, databaseName string) (bool, error) {
	process := func(row *sql.Row, out *int) error {
		return errors.Wrapf(row.Scan(out), "unable to fetch row")
	}
	count, err := ReadSingle(ctx, db, process, fmt.Sprintf(`SELECT count(*) FROM pg_database WHERE datname='%s'`, databaseName))
	if err != nil {
		return false, err
	}
	return *count == 1, nil
}

func CreateDatabase(ctx context.Context, db *sql.DB, databaseName string) error {
	// TODO why doesn't this work?
	// _, err := RunStatement(ctx, db, `create database "$1" encoding UTF8`, databaseName)
	_, err := RunStatement(ctx, db, fmt.Sprintf(`create database "%s" encoding UTF8`, databaseName))
	return err
}

func CreateDatabaseIfNotExists(ctx context.Context, db *sql.DB, databaseName string) error {
	logrus.Debugf("creating database '%s' if not exists", databaseName)
	exists, err := DoesDatabaseExist(ctx, db, databaseName)
	if err != nil {
		return err
	}
	if exists {
		logrus.Debugf("skipping creation of database %s, already exists", databaseName)
		return nil
	}
	logrus.Debugf("database '%s' does not exist: creating", databaseName)
	_, err = RunStatement(ctx, db, fmt.Sprintf(`create database "%s" encoding UTF8`, databaseName))
	return err
}

func InitializeSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, uuidOsspExtention)
	if err != nil {
		return errors.Wrapf(err, "unable to create extension")
	}
	for _, table := range []string{usersTable, followersTable, messagesTable, upvotesTable, topicsTable, pingsTable} {
		_, err = RunStatement(ctx, db, table)
		if err != nil {
			return err
		}
	}
	return nil
}
