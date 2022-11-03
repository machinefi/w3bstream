package scanner

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func Scan(ctx context.Context, rows *sql.Rows, v interface{}) error {
	if rows == nil {
		return nil
	}
	defer rows.Close()

	iter, err := ScanIteratorFor(v)
	if err != nil {
		return err
	}

	for rows.Next() {
		item := iter.New()
		if err := ScanRows(ctx, rows, item); err != nil {
			return err
		}
		if err := iter.Next(item); err != nil {
			return err
		}
	}

	if hasRecord, ok := iter.(interface{ HasRecord() bool }); ok {
		if !hasRecord.HasRecord() {
			return RecordNotFound
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	return nil
}

var RecordNotFound = errors.Errorf("record is not found")
