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
	// internal error
	InternalServerError Error = http.StatusInternalServerError*1e6 + iota + 1
	// @errTalk Database Error
	DatabaseError
	// @errTalk Upload File Failed
	UploadFileFailed
	// @errTalk Create Message Channel Failed
	CreateChannelFailed
	// @errTalk Fetch Resource Failed
	FetchResourceFailed
	// @errTalk Config Init Failed
	ConfigInitFailed
	// @errTalk Config Uninit Failed
	ConfigUninitFailed
	// @errTalk Config Parse Failed
	ConfigParseFailed
	// @errTalk Gen Publisher Token Failed
	GenPublisherTokenFailed
	// @errTalk Create Instance Failed
	CreateInstanceFailed
	// @errTalk Batch Remove Applet Failed
	BatchRemoveAppletFailed
	// @errTalk Md5 Checksum Failed
	MD5ChecksumFailed
)

const (
	// unauthorized
	Unauthorized Error = http.StatusUnauthorized*1e6 + iota + 1
	// @errTalk Invalid Auth Value
	InvalidAuthValue
	// @errTalk Invalid Auth Account ID
	InvalidAuthAccountID
	// @errTalk No Project Permission
	NoProjectPermission
	// @errTalk No Admin Permission
	NoAdminPermission
	// @errTalk Invalid Old Password
	InvalidOldPassword
	// @errTalk Invalid New Password
	InvalidNewPassword
	// @errTalk Invalid Password
	InvalidPassword
	// @errTalk Invalid Siwe Signature
	InvalidEthLoginSignature
	// @errTalk Invalid Siwe Message
	InvalidEthLoginMessage
	// @errTalk Invalid Auth Publisher ID
	InvalidAuthPublisherID
	// @errTalk Current Account Absence
	CurrentAccountAbsence
	// @errTalk Invalid Event Channel
	InvalidEventChannel
)

const (
	// forbidden
	Forbidden Error = http.StatusForbidden*1e6 + iota + 1
	// @errTalk Disabled Account
	DisabledAccount
	// @errTalk White List Forbidden
	WhiteListForbidden
	// @errTalk Upload File Size Limit
	UploadFileSizeLimit
	// @errTalk Upload File Md5 Unmatched
	UploadFileMd5Unmatched
	// @errTalk Upload File Disk Limit
	UploadFileDiskLimit
)

const (
	// @errTalk Conflict conflict error
	Conflict Error = http.StatusConflict*1e6 + iota + 1
	// @errTalk Project Name Conflict
	ProjectNameConflict
	// @errTalk Resource Conflict
	ResourceConflict
	// @errTalk Resource Owner Conflict
	ResourceOwnerConflict
	// @errTalk Strategy Conflict
	StrategyConflict
	// @errTalk Config Conflict
	ConfigConflict
	// @errTalk Publisher Conflict
	PublisherConflict
	// @errTalk Multi Instance Deployed
	MultiInstanceDeployed
	// @errTalk Applet Name Conflict
	AppletNameConflict
	// @errTalk Cron Job Conflict
	CronJobConflict
)

const (
	// @errTalk BadRequest
	BadRequest Error = http.StatusBadRequest*1e6 + iota + 1
	// @errTalk Invalid Config Type
	InvalidConfigType
	// @errTalk Deprecated Project
	DeprecatedProject
	// @errTalk Unknown Deploy Command
	UnknownDeployCommand
	// @errTalk Invalid Cron Expressions
	InvalidCronExpressions
)

const (
	// @errTalk NotFound
	NotFound Error = http.StatusNotFound*1e6 + iota + 1
	// @errTalk Project Not Found
	ProjectNotFound
	// @errTalk Config Not Found
	ConfigNotFound
	// @errTalk Resource Not Found
	ResourceNotFound
	// @errTalk Applet Not Found
	AppletNotFound
	// @errTalk Instance Not Found
	InstanceNotFound
	// @errTalk Strategy Not Found
	StrategyNotFound
	// @errTalk Publisher Not Found
	PublisherNotFound
	// @errTalk Account Identity Not Found
	AccountIdentityNotFound
	// @errTalk Resource Perm Not Found
	ResourcePermNotFound
	// @errTalk Cron Job Not Found
	CronJobNotFound
)

// Deprecated: pls check database error and return defined status error
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
