// This is a generated source file. DO NOT EDIT
// Source: jwt/error__generated.go

package jwt

import "github.com/machinefi/w3bstream/pkg/depends/kit/statusx"

var _ statusx.Error = (*Error)(nil)

func (v Error) StatusErr() *statusx.StatusErr {
	return &statusx.StatusErr{
		Key:       v.Key(),
		Code:      v.Code(),
		Msg:       v.Msg(),
		CanBeTalk: v.CanBeTalk(),
	}
}

func (v Error) Unwrap() error {
	return v.StatusErr()
}

func (v Error) Error() string {
	return v.StatusErr().Error()
}

func (v Error) StatusCode() int {
	return statusx.StatusCodeFromCode(int(v))
}

func (v Error) Code() int {
	if with, ok := (interface{})(v).(statusx.ServiceCode); ok {
		return with.ServiceCode() + int(v)
	}
	return int(v)
}

func (v Error) Key() string {
	switch v {
	case Unauthorized:
		return "Unauthorized"
	case InvalidToken:
		return "InvalidToken"
	case InvalidClaim:
		return "InvalidClaim"
	}
	return "UNKNOWN"
}

func (v Error) Msg() string {
	switch v {
	case Unauthorized:
		return ""
	case InvalidToken:
		return "Invalid Token"
	case InvalidClaim:
		return "Invalid Claim"
	}
	return "-"
}

func (v Error) CanBeTalk() bool {
	switch v {
	case Unauthorized:
		return false
	case InvalidToken:
		return true
	case InvalidClaim:
		return true
	}
	return false
}
