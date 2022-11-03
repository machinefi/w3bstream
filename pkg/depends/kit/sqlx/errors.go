package sqlx

import (
	"fmt"

	"github.com/pkg/errors"
)

var ErrNotTx = errors.New("db is not *sql.Tx")
var ErrNotDB = errors.New("db is not *sql.DB")

func NewSqlError(tpe sqlErrType, msg string) *SqlError {
	return &SqlError{
		Type: tpe,
		Msg:  msg,
	}
}

type SqlError struct {
	Type sqlErrType
	Msg  string
}

func (e *SqlError) Error() string {
	return fmt.Sprintf("Sqlx [%s] %s", e.Type, e.Msg)
}

type sqlErrType string

var (
	sqlErrTypeNotFound sqlErrType = "NotFound"
	sqlErrTypeConflict sqlErrType = "Conflict"
)

var DuplicateEntryErrNumber uint16 = 1062

func DBErr(err error) *dbErr {
	return &dbErr{
		err: err,
	}
}

type dbErr struct {
	err error

	errDefault  error
	errNotFound error
	errConflict error
}

func (r dbErr) WithNotFound(err error) *dbErr {
	r.errNotFound = err
	return &r
}

func (r dbErr) WithDefault(err error) *dbErr {
	r.errDefault = err
	return &r
}

func (r dbErr) WithConflict(err error) *dbErr {
	r.errConflict = err
	return &r
}

func (r *dbErr) IsNotFound() bool {
	if sqlErr, ok := UnwrapAll(r.err).(*SqlError); ok {
		return sqlErr.Type == sqlErrTypeNotFound
	}
	return false
}

func (r *dbErr) IsConflict() bool {
	if sqlErr, ok := UnwrapAll(r.err).(*SqlError); ok {
		return sqlErr.Type == sqlErrTypeConflict
	}
	return false
}

func (r *dbErr) Err() error {
	if r.err == nil {
		return nil
	}
	if e, ok := UnwrapAll(r.err).(*SqlError); ok {
		switch e.Type {
		case sqlErrTypeNotFound:
			if r.errNotFound != nil {
				return r.errNotFound
			}
		case sqlErrTypeConflict:
			if r.errConflict != nil {
				return r.errConflict
			}
		}
		if r.errDefault != nil {
			return r.errDefault
		}
	}
	return r.err
}

func UnwrapAll(err error) error {
	for {
		if cause := UnwrapOnce(err); cause != nil {
			err = cause
			continue
		}
		break
	}
	return err
}

func UnwrapOnce(err error) (cause error) {
	switch e := err.(type) {
	case interface{ Cause() error }:
		return e.Cause()
	case interface{ Unwrap() error }:
		return e.Unwrap()
	}
	return nil
}

// TODO should import some driver kit error list
