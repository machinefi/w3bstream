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
	case WhiteListForbidden:
		return "WhiteListForbidden"
	case NotFound:
		return "NotFound"
	case ProjectNotFound:
		return "ProjectNotFound"
	case Conflict:
		return "Conflict"
	case ProjectNameConflict:
		return "ProjectNameConflict"
	case InternalServerError:
		return "InternalServerError"
	case DatabaseError:
		return "DatabaseError"
	case UploadFileFailed:
		return "UploadFileFailed"
	case CreateChannelFailed:
		return "CreateChannelFailed"
	}
	return "UNKNOWN"
}

func (v Error) Msg() string {
	switch v {
	case BadRequest:
		return "BadRequest"
	case MD5ChecksumFailed:
		return "Md5 Checksum Failed"
	case InvalidChainClient:
		return "Invalid Chain Client"
	case Unauthorized:
		return "Unauthorized unauthorized"
	case InvalidAuthValue:
		return "Invalid Auth Value"
	case InvalidAuthAccountID:
		return "Invalid Auth Account ID"
	case NoProjectPermission:
		return "No Project Permission"
	case NoAdminPermission:
		return "No Admin Permission"
	case InvalidOldPassword:
		return "Invalid Old Password"
	case InvalidNewPassword:
		return "Invalid New Password"
	case InvalidPassword:
		return "Invalid Password"
	case InvalidEthLoginSignature:
		return "Invalid Siwe Signature"
	case InvalidEthLoginMessage:
		return "Invalid Siwe Message"
	case Forbidden:
		return "Forbidden"
	case InstanceLimit:
		return "deployed instance limit"
	case DisabledAccount:
		return "Disabled Account"
	case WhiteListForbidden:
		return "White List Forbidden"
	case NotFound:
		return "NotFound"
	case ProjectNotFound:
		return "Project Not Found"
	case Conflict:
		return "Conflict conflict error"
	case ProjectNameConflict:
		return "Project Name Conflict"
	case InternalServerError:
		return "InternalServerError internal error"
	case DatabaseError:
		return "Database Error"
	case UploadFileFailed:
		return "Upload File Failed"
	case CreateChannelFailed:
		return "Create Message Channel Failed"
	}
	return "-"
}

func (v Error) CanBeTalk() bool {
	switch v {
	case BadRequest:
		return true
	case MD5ChecksumFailed:
		return true
	case InvalidChainClient:
		return true
	case Unauthorized:
		return true
	case InvalidAuthValue:
		return true
	case InvalidAuthAccountID:
		return true
	case NoProjectPermission:
		return true
	case NoAdminPermission:
		return true
	case InvalidOldPassword:
		return true
	case InvalidNewPassword:
		return true
	case InvalidPassword:
		return true
	case InvalidEthLoginSignature:
		return true
	case InvalidEthLoginMessage:
		return true
	case Forbidden:
		return true
	case InstanceLimit:
		return true
	case DisabledAccount:
		return true
	case WhiteListForbidden:
		return true
	case NotFound:
		return true
	case ProjectNotFound:
		return true
	case Conflict:
		return true
	case ProjectNameConflict:
		return true
	case InternalServerError:
		return true
	case DatabaseError:
		return true
	case UploadFileFailed:
		return true
	case CreateChannelFailed:
		return true
	}
	return false
}
