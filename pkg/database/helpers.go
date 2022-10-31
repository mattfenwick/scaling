package database

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
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
