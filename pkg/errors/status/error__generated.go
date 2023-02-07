// This is a generated source file. DO NOT EDIT
// Source: status/error__generated.go

package status

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
	case BadRequest:
		return "BadRequest"
	case MD5ChecksumFailed:
		return "MD5ChecksumFailed"
	case InvalidChainClient:
		return "InvalidChainClient"
	case Unauthorized:
		return "Unauthorized"
	case InvalidAuthValue:
		return "InvalidAuthValue"
	case InvalidAuthAccountID:
		return "InvalidAuthAccountID"
	case NoProjectPermission:
		return "NoProjectPermission"
	case Forbidden:
		return "Forbidden"
	case InstanceLimit:
		return "InstanceLimit"
	case NotFound:
		return "NotFound"
	case Conflict:
		return "Conflict"
	case InternalServerError:
		return "InternalServerError"
	case UploadFileFailed:
		return "UploadFileFailed"
	case ExtractFileFailed:
		return "ExtractFileFailed"
	case LoadVMFailed:
		return "LoadVMFailed"
	}
	return "UNKNOWN"
}

func (v Error) Msg() string {
	switch v {
	case BadRequest:
		return "BadRequest"
	case MD5ChecksumFailed:
		return ""
	case InvalidChainClient:
		return ""
	case Unauthorized:
		return "Unauthorized unauthorized"
	case InvalidAuthValue:
		return ""
	case InvalidAuthAccountID:
		return ""
	case NoProjectPermission:
		return ""
	case Forbidden:
		return ""
	case InstanceLimit:
		return "deployed instance limit"
	case NotFound:
		return "NotFound"
	case Conflict:
		return "Conflict conflict error"
	case InternalServerError:
		return "InternalServerError internal error"
	case UploadFileFailed:
		return ""
	case ExtractFileFailed:
		return ""
	case LoadVMFailed:
		return ""
	}
	return "-"
}

func (v Error) CanBeTalk() bool {
	switch v {
	case BadRequest:
		return false
	case MD5ChecksumFailed:
		return false
	case InvalidChainClient:
		return false
	case Unauthorized:
		return true
	case InvalidAuthValue:
		return false
	case InvalidAuthAccountID:
		return false
	case NoProjectPermission:
		return false
	case Forbidden:
		return false
	case InstanceLimit:
		return true
	case NotFound:
		return false
	case Conflict:
		return false
	case InternalServerError:
		return false
	case UploadFileFailed:
		return false
	case ExtractFileFailed:
		return false
	case LoadVMFailed:
		return false
	}
	return false
}
