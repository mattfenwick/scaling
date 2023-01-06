package database

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func ReadMany[A any](ctx context.Context, db *sql.DB, process func(*sql.Rows, *A) error, query string, args ...any) ([]*A, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to issue query")
	}

	var records []*A

	for rows.Next() {
		var record A
		err = process(rows, &record)
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

func ReadSingle[A any](ctx context.Context, db *sql.DB, process func(*sql.Row, *A) error, query string, args ...any) (*A, error) {
	var record A
	if err := process(db.QueryRowContext(ctx, query, args...), &record); err != nil {
		// TODO does ErrNoRows specifically matter?
		// if err == sql.ErrNoRows {
		// 	return ???
		// }
		return nil, errors.Wrapf(err, "unable to run query '%s' with args '%+v'", query, args)
	}
	logrus.Tracef("ReadSingle result: %+v", record)
	return &record, nil
}

func RunStatement(ctx context.Context, db *sql.DB, query string, args ...any) (sql.Result, error) {
	logrus.Tracef("running SQL query: '%s' with args '%+v'", query, args)
	result, err := db.ExecContext(ctx, query, args...)
	return result, errors.Wrapf(err, "unable to run query '%s' with args '%+v'", query, args)
}
