package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"net/url"

	"github.com/pkg/errors"
)

func Connect(user string, password string, host string, database string) (*sql.DB, error) {
	dbConnectionURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user,
		url.QueryEscape(password),
		host,
		5432,
		database)
	logrus.Infof("attempting to connect to postgres at '%s'", dbConnectionURL)
	db, err := sql.Open("postgres", dbConnectionURL)
	if err != nil {
		return db, errors.Wrap(err, "unable to open postgres connection")
	}
	if err := db.PingContext(context.TODO()); err != nil {
		return db, errors.Wrap(err, "unable to ping postgres")
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(0)

	return db, nil
}
