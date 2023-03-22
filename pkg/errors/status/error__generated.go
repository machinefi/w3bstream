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
	case NoAdminPermission:
		return "NoAdminPermission"
	case InvalidOldPassword:
		return "InvalidOldPassword"
	case InvalidNewPassword:
		return "InvalidNewPassword"
	case InvalidPassword:
		return "InvalidPassword"
	case InvalidEthLoginSignature:
		return "InvalidEthLoginSignature"
	case InvalidEthLoginMessage:
		return "InvalidEthLoginMessage"
	case Forbidden:
		return "Forbidden"
	case InstanceLimit:
		return "InstanceLimit"
	case DisabledAccount:
		return "DisabledAccount"
	case NotFound:
		return "NotFound"
	case Conflict:
		return "Conflict"
	case InternalServerError:
		return "InternalServerError"
	case UploadFileFailed:
		return "UploadFileFailed"
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
	case NoAdminPermission:
		return ""
	case InvalidOldPassword:
		return ""
	case InvalidNewPassword:
		return ""
	case InvalidPassword:
		return ""
	case InvalidEthLoginSignature:
		return ""
	case InvalidEthLoginMessage:
		return ""
	case Forbidden:
		return ""
	case InstanceLimit:
		return "deployed instance limit"
	case DisabledAccount:
		return ""
	case NotFound:
		return "NotFound"
	case Conflict:
		return "Conflict conflict error"
	case InternalServerError:
		return "InternalServerError internal error"
	case UploadFileFailed:
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
	case NoAdminPermission:
		return false
	case InvalidOldPassword:
		return false
	case InvalidNewPassword:
		return false
	case InvalidPassword:
		return false
	case InvalidEthLoginSignature:
		return false
	case InvalidEthLoginMessage:
		return false
	case Forbidden:
		return false
	case InstanceLimit:
		return true
	case DisabledAccount:
		return false
	case NotFound:
		return false
	case Conflict:
		return false
	case InternalServerError:
		return false
	case UploadFileFailed:
		return false
	}
	return false
}
