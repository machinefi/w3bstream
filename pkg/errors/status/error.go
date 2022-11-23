package status

import (
	"net/http"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
)

//go:generate toolkit gen status Error
type Error int

func (Error) ServiceCode() int {
	return 999 * 1e3
}

const (
	// InternalServerError internal error
	InternalServerError Error = http.StatusInternalServerError*1e6 + iota + 1
	UploadFileFailed
	ExtractFileFailed
	LoadVMFailed
)

const (
	// @errTalk Unauthorized unauthorized
	Unauthorized Error = http.StatusUnauthorized*1e6 + iota + 1
	InvalidAuthValue
	InvalidAuthAccountID
	NoProjectPermission
)

const (
	Forbidden Error = http.StatusForbidden*1e6 + iota + 1
	// @errTalk deployed instance limit
	InstanceLimit
)

const (
	// Conflict conflict error
	Conflict Error = http.StatusConflict*1e6 + iota + 1
)

const (
	// BadRequest
	BadRequest Error = http.StatusBadRequest*1e6 + iota + 1
	MD5ChecksumFailed
)

const (
	// NotFound
	NotFound Error = http.StatusNotFound*1e6 + iota + 1
)

func CheckDatabaseError(err error, msg ...string) error {
	desc := ""
	if len(msg) > 0 {
		desc = msg[0]
	}
	if err != nil {
		desc = desc + ":" + err.Error()
		e := sqlx.DBErr(err)
		if e.IsNotFound() {
			return NotFound.StatusErr().WithDesc(desc)
		} else if e.IsConflict() {
			return Conflict.StatusErr().WithDesc(desc)
		} else {
			desc = desc + " " + err.Error()
			return InternalServerError.StatusErr().WithDesc(desc)
		}
	}
	return nil
}
