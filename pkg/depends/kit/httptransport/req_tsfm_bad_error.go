package httptransport

import (
	"net/http"

	pkgerr "github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	vldterr "github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
)

type BadRequestError interface {
	EnableErrTalk()
	SetMsg(msg string)
	AddErr(err error, nameOrIdx ...interface{})
}

type badRequest struct {
	errorFields statusx.ErrorFields
	errTalk     bool
	msg         string
}

func (e *badRequest) EnableErrTalk() { e.errTalk = true }

func (e *badRequest) SetMsg(msg string) { e.msg = msg }

func (e *badRequest) AddErr(err error, nameOrIdx ...interface{}) {
	if len(nameOrIdx) > 1 {
		e.errorFields = append(e.errorFields, &statusx.ErrorField{
			In:    nameOrIdx[0].(string),
			Field: vldterr.KeyPath(nameOrIdx[1:]).String(),
			Msg:   err.Error(),
		})
	}
}

func (e *badRequest) Err() error {
	if len(e.errorFields) == 0 {
		return nil
	}

	msg := e.msg
	if msg == "" {
		msg = "invalid parameters"
	}

	err := statusx.Wrap(pkgerr.New(""), http.StatusBadRequest, "badRequest").
		WithMsg(msg).
		AppendErrorFields(e.errorFields...)

	if e.errTalk {
		err = err.EnableErrTalk()
	}

	return err
}
