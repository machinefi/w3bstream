package status

import "net/http"

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
	// @errTalk MQTT Subscribe Failed
	MqttSubscribeFailed
	// @errTalk MQTT Connect Failed
	MqttConnectFailed
	// @errTalk Batch Remove WasmLog Failed
	BatchRemoveWasmLogFailed
	// @errTalk Gen Token Failed
	GenTokenFailed
	// @errTalk Traffic Limit Exceeded Failed
	TrafficLimitExceededFailed
	// @errTalk Create Traffic Scheduler Failed
	CreateTrafficSchedulerFailed
	// @errTalk Update Traffic Scheduler Failed
	UpdateTrafficSchedulerFailed
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
	// @errTalk No Operator Permission
	NoOperatorPermission
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
	// @errTalk Invalid Event Token
	InvalidEventToken
	// @errTalk Invalid Applet
	InvalidAppletContext
	// @errTalk No Resource Permission
	NoResourcePermission
	// @errTalk Invalid Account Access Key
	InvalidAccountAccessKey
	// @errTalk Account Access Key Expired
	AccountAccessKeyExpired
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
	// @errTalk Topic Already Subscribed
	TopicAlreadySubscribed
	// @errTalk Occupied Operator
	OccupiedOperator
	// @errTalk Unsupported FileSystem Operator
	UnsupportedFSOperator
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
	// @errTalk Contract Log Conflict
	ContractLogConflict
	// @errTalk Chain Tx Conflict
	ChainTxConflict
	// @errTalk Chain Height Conflict
	ChainHeightConflict
	// @errTalk Account Identity Conflict
	AccountIdentityConflict
	// @errTalk Account Conflict
	AccountConflict
	// @errTalk Account Password Conflict
	AccountPasswordConflict
	// @errTalk Operator Conflict
	OperatorConflict
	// @errTalk Traffic Limit Conflict
	TrafficLimitConflict
	// @errTalk Project Operator relationship Conflict
	ProjectOperatorConflict
	// @errTalk Account Key Name Conflict
	AccountKeyNameConflict
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
	// @errTalk Invalid Private Key
	InvalidPrivateKey
	// @errTalk Invalid Delete Condition
	InvalidDeleteCondition
	// @errTalk Unknown Deploy Command
	UnknownMonitorCommand
	// @errTalk Invalid Contract Log IDs
	InvalidContractLogIDs
	// @errTalk Invalid Chain Tx IDs
	InvalidChainTxIDs
	// @errTalk Invalid Chain Height IDs
	InvalidChainHeightIDs
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
	// @errTalk Instance Not Running
	InstanceNotRunning
	// @errTalk Blockchain Not Found
	BlockchainNotFound
	// @errTalk Contract Log Not Found
	ContractLogNotFound
	// @errTalk Chain Tx Not Found
	ChainTxNotFound
	// @errTalk Chain Height Not Found
	ChainHeightNotFound
	// @errTalk Account Not Found
	AccountNotFound
	// @errTalk Account Password Not Found
	AccountPasswordNotFound
	// @errTalk Operator Not Found
	OperatorNotFound
	// @errTalk Traffic Limit Not Found
	TrafficLimitNotFound
	// @errTalk Project Operator relationship Not Found
	ProjectOperatorNotFound
	// @errTalk Account Key Not Found
	AccountKeyNotFound
)
