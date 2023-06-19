package mock_sqlx

import (
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
)

var (
	ErrConflict = sqlx.NewSqlError(sqlx.SqlErrTypeConflict, "")
	ErrNotFound = sqlx.NewSqlError(sqlx.SqlErrTypeNotFound, "")
	ErrDatabase = errors.New("database error")
)
