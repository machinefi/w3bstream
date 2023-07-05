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
	case InvalidConfigType:
		return "InvalidConfigType"
	case DeprecatedProject:
		return "DeprecatedProject"
	case UnknownDeployCommand:
		return "UnknownDeployCommand"
	case InvalidCronExpressions:
		return "InvalidCronExpressions"
	case InvalidPrivateKey:
		return "InvalidPrivateKey"
	case InvalidDeleteCondition:
		return "InvalidDeleteCondition"
	case UnknownMonitorCommand:
		return "UnknownMonitorCommand"
	case InvalidContractLogIDs:
		return "InvalidContractLogIDs"
	case InvalidChainTxIDs:
		return "InvalidChainTxIDs"
	case InvalidChainHeightIDs:
		return "InvalidChainHeightIDs"
	case InvalidVMState:
		return "InvalidVMState"
	case InvalidAccessKeyIdentityType:
		return "InvalidAccessKeyIdentityType"
	case Unauthorized:
		return "Unauthorized"
	case InvalidAuthValue:
		return "InvalidAuthValue"
	case InvalidAuthAccountID:
		return "InvalidAuthAccountID"
	case NoProjectPermission:
		return "NoProjectPermission"
	case NoOperatorPermission:
		return "NoOperatorPermission"
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
	case InvalidAuthPublisherID:
		return "InvalidAuthPublisherID"
	case CurrentAccountAbsence:
		return "CurrentAccountAbsence"
	case InvalidEventChannel:
		return "InvalidEventChannel"
	case InvalidEventToken:
		return "InvalidEventToken"
	case InvalidAppletContext:
		return "InvalidAppletContext"
	case NoResourcePermission:
		return "NoResourcePermission"
	case InvalidAccessKey:
		return "InvalidAccessKey"
	case AccessKeyExpired:
		return "AccessKeyExpired"
	case Forbidden:
		return "Forbidden"
	case DisabledAccount:
		return "DisabledAccount"
	case WhiteListForbidden:
		return "WhiteListForbidden"
	case UploadFileSizeLimit:
		return "UploadFileSizeLimit"
	case UploadFileMd5Unmatched:
		return "UploadFileMd5Unmatched"
	case UploadFileDiskLimit:
		return "UploadFileDiskLimit"
	case TopicAlreadySubscribed:
		return "TopicAlreadySubscribed"
	case OccupiedOperator:
		return "OccupiedOperator"
	case UnsupportedFSOperator:
		return "UnsupportedFSOperator"
	case NotFound:
		return "NotFound"
	case ProjectNotFound:
		return "ProjectNotFound"
	case ConfigNotFound:
		return "ConfigNotFound"
	case ResourceNotFound:
		return "ResourceNotFound"
	case AppletNotFound:
		return "AppletNotFound"
	case InstanceNotFound:
		return "InstanceNotFound"
	case StrategyNotFound:
		return "StrategyNotFound"
	case PublisherNotFound:
		return "PublisherNotFound"
	case AccountIdentityNotFound:
		return "AccountIdentityNotFound"
	case ResourcePermNotFound:
		return "ResourcePermNotFound"
	case CronJobNotFound:
		return "CronJobNotFound"
	case InstanceNotRunning:
		return "InstanceNotRunning"
	case BlockchainNotFound:
		return "BlockchainNotFound"
	case ContractLogNotFound:
		return "ContractLogNotFound"
	case ChainTxNotFound:
		return "ChainTxNotFound"
	case ChainHeightNotFound:
		return "ChainHeightNotFound"
	case AccountNotFound:
		return "AccountNotFound"
	case AccountPasswordNotFound:
		return "AccountPasswordNotFound"
	case OperatorNotFound:
		return "OperatorNotFound"
	case TrafficLimitNotFound:
		return "TrafficLimitNotFound"
	case ProjectOperatorNotFound:
		return "ProjectOperatorNotFound"
	case AccessKeyNotFound:
		return "AccessKeyNotFound"
	case Conflict:
		return "Conflict"
	case ProjectNameConflict:
		return "ProjectNameConflict"
	case ResourceConflict:
		return "ResourceConflict"
	case ResourceOwnerConflict:
		return "ResourceOwnerConflict"
	case StrategyConflict:
		return "StrategyConflict"
	case ConfigConflict:
		return "ConfigConflict"
	case PublisherConflict:
		return "PublisherConflict"
	case MultiInstanceDeployed:
		return "MultiInstanceDeployed"
	case AppletNameConflict:
		return "AppletNameConflict"
	case CronJobConflict:
		return "CronJobConflict"
	case ContractLogConflict:
		return "ContractLogConflict"
	case ChainTxConflict:
		return "ChainTxConflict"
	case ChainHeightConflict:
		return "ChainHeightConflict"
	case AccountIdentityConflict:
		return "AccountIdentityConflict"
	case AccountConflict:
		return "AccountConflict"
	case AccountPasswordConflict:
		return "AccountPasswordConflict"
	case OperatorConflict:
		return "OperatorConflict"
	case TrafficLimitConflict:
		return "TrafficLimitConflict"
	case ProjectOperatorConflict:
		return "ProjectOperatorConflict"
	case AccessKeyNameConflict:
		return "AccessKeyNameConflict"
	case InternalServerError:
		return "InternalServerError"
	case DatabaseError:
		return "DatabaseError"
	case UploadFileFailed:
		return "UploadFileFailed"
	case FetchResourceFailed:
		return "FetchResourceFailed"
	case ConfigInitFailed:
		return "ConfigInitFailed"
	case ConfigUninitFailed:
		return "ConfigUninitFailed"
	case ConfigParseFailed:
		return "ConfigParseFailed"
	case CreateInstanceFailed:
		return "CreateInstanceFailed"
	case BatchRemoveAppletFailed:
		return "BatchRemoveAppletFailed"
	case MqttSubscribeFailed:
		return "MqttSubscribeFailed"
	case MqttConnectFailed:
		return "MqttConnectFailed"
	case BatchRemoveWasmLogFailed:
		return "BatchRemoveWasmLogFailed"
	case GenTokenFailed:
		return "GenTokenFailed"
	case TrafficLimitExceededFailed:
		return "TrafficLimitExceededFailed"
	case CreateTrafficSchedulerFailed:
		return "CreateTrafficSchedulerFailed"
	case UpdateTrafficSchedulerFailed:
		return "UpdateTrafficSchedulerFailed"
	}
	return "UNKNOWN"
}

func (v Error) Msg() string {
	switch v {
	case BadRequest:
		return "BadRequest"
	case InvalidConfigType:
		return "Invalid Config Type"
	case DeprecatedProject:
		return "Deprecated Project"
	case UnknownDeployCommand:
		return "Unknown Deploy Command"
	case InvalidCronExpressions:
		return "Invalid Cron Expressions"
	case InvalidPrivateKey:
		return "Invalid Private Key"
	case InvalidDeleteCondition:
		return "Invalid Delete Condition"
	case UnknownMonitorCommand:
		return "Unknown Deploy Command"
	case InvalidContractLogIDs:
		return "Invalid Contract Log IDs"
	case InvalidChainTxIDs:
		return "Invalid Chain Tx IDs"
	case InvalidChainHeightIDs:
		return "Invalid Chain Height IDs"
	case InvalidVMState:
		return "Invalid VM State"
	case InvalidAccessKeyIdentityType:
		return "Invalid Access Key Identity Type"
	case Unauthorized:
		return "unauthorized"
	case InvalidAuthValue:
		return "Invalid Auth Value"
	case InvalidAuthAccountID:
		return "Invalid Auth Account ID"
	case NoProjectPermission:
		return "No Project Permission"
	case NoOperatorPermission:
		return "No Operator Permission"
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
	case InvalidAuthPublisherID:
		return "Invalid Auth Publisher ID"
	case CurrentAccountAbsence:
		return "Current Account Absence"
	case InvalidEventChannel:
		return "Invalid Event Channel"
	case InvalidEventToken:
		return "Invalid Event Token"
	case InvalidAppletContext:
		return "Invalid Applet"
	case NoResourcePermission:
		return "No Resource Permission"
	case InvalidAccessKey:
		return "Invalid Account Access Key"
	case AccessKeyExpired:
		return "Account Access Key Expired"
	case Forbidden:
		return "forbidden"
	case DisabledAccount:
		return "Disabled Account"
	case WhiteListForbidden:
		return "White List Forbidden"
	case UploadFileSizeLimit:
		return "Upload File Size Limit"
	case UploadFileMd5Unmatched:
		return "Upload File Md5 Unmatched"
	case UploadFileDiskLimit:
		return "Upload File Disk Limit"
	case TopicAlreadySubscribed:
		return "Topic Already Subscribed"
	case OccupiedOperator:
		return "Occupied Operator"
	case UnsupportedFSOperator:
		return "Unsupported FileSystem Operator"
	case NotFound:
		return "NotFound"
	case ProjectNotFound:
		return "Project Not Found"
	case ConfigNotFound:
		return "Config Not Found"
	case ResourceNotFound:
		return "Resource Not Found"
	case AppletNotFound:
		return "Applet Not Found"
	case InstanceNotFound:
		return "Instance Not Found"
	case StrategyNotFound:
		return "Strategy Not Found"
	case PublisherNotFound:
		return "Publisher Not Found"
	case AccountIdentityNotFound:
		return "Account Identity Not Found"
	case ResourcePermNotFound:
		return "Resource Perm Not Found"
	case CronJobNotFound:
		return "Cron Job Not Found"
	case InstanceNotRunning:
		return "Instance Not Running"
	case BlockchainNotFound:
		return "Blockchain Not Found"
	case ContractLogNotFound:
		return "Contract Log Not Found"
	case ChainTxNotFound:
		return "Chain Tx Not Found"
	case ChainHeightNotFound:
		return "Chain Height Not Found"
	case AccountNotFound:
		return "Account Not Found"
	case AccountPasswordNotFound:
		return "Account Password Not Found"
	case OperatorNotFound:
		return "Operator Not Found"
	case TrafficLimitNotFound:
		return "Traffic Limit Not Found"
	case ProjectOperatorNotFound:
		return "Project Operator relationship Not Found"
	case AccessKeyNotFound:
		return "Account Key Not Found"
	case Conflict:
		return "Conflict conflict error"
	case ProjectNameConflict:
		return "Project Name Conflict"
	case ResourceConflict:
		return "Resource Conflict"
	case ResourceOwnerConflict:
		return "Resource Owner Conflict"
	case StrategyConflict:
		return "Strategy Conflict"
	case ConfigConflict:
		return "Config Conflict"
	case PublisherConflict:
		return "Publisher Conflict"
	case MultiInstanceDeployed:
		return "Multi Instance Deployed"
	case AppletNameConflict:
		return "Applet Name Conflict"
	case CronJobConflict:
		return "Cron Job Conflict"
	case ContractLogConflict:
		return "Contract Log Conflict"
	case ChainTxConflict:
		return "Chain Tx Conflict"
	case ChainHeightConflict:
		return "Chain Height Conflict"
	case AccountIdentityConflict:
		return "Account Identity Conflict"
	case AccountConflict:
		return "Account Conflict"
	case AccountPasswordConflict:
		return "Account Password Conflict"
	case OperatorConflict:
		return "Operator Conflict"
	case TrafficLimitConflict:
		return "Traffic Limit Conflict"
	case ProjectOperatorConflict:
		return "Project Operator relationship Conflict"
	case AccessKeyNameConflict:
		return "Access Key Name Conflict"
	case InternalServerError:
		return "internal error"
	case DatabaseError:
		return "Database Error"
	case UploadFileFailed:
		return "Upload File Failed"
	case FetchResourceFailed:
		return "Fetch Resource Failed"
	case ConfigInitFailed:
		return "Config Init Failed"
	case ConfigUninitFailed:
		return "Config Uninit Failed"
	case ConfigParseFailed:
		return "Config Parse Failed"
	case CreateInstanceFailed:
		return "Create Instance Failed"
	case BatchRemoveAppletFailed:
		return "Batch Remove Applet Failed"
	case MqttSubscribeFailed:
		return "MQTT Subscribe Failed"
	case MqttConnectFailed:
		return "MQTT Connect Failed"
	case BatchRemoveWasmLogFailed:
		return "Batch Remove WasmLog Failed"
	case GenTokenFailed:
		return "Gen Token Failed"
	case TrafficLimitExceededFailed:
		return "Traffic Limit Exceeded Failed"
	case CreateTrafficSchedulerFailed:
		return "Create Traffic Scheduler Failed"
	case UpdateTrafficSchedulerFailed:
		return "Update Traffic Scheduler Failed"
	}
	return "-"
}

func (v Error) CanBeTalk() bool {
	switch v {
	case BadRequest:
		return true
	case InvalidConfigType:
		return true
	case DeprecatedProject:
		return true
	case UnknownDeployCommand:
		return true
	case InvalidCronExpressions:
		return true
	case InvalidPrivateKey:
		return true
	case InvalidDeleteCondition:
		return true
	case UnknownMonitorCommand:
		return true
	case InvalidContractLogIDs:
		return true
	case InvalidChainTxIDs:
		return true
	case InvalidChainHeightIDs:
		return true
	case InvalidVMState:
		return true
	case InvalidAccessKeyIdentityType:
		return true
	case Unauthorized:
		return false
	case InvalidAuthValue:
		return true
	case InvalidAuthAccountID:
		return true
	case NoProjectPermission:
		return true
	case NoOperatorPermission:
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
	case InvalidAuthPublisherID:
		return true
	case CurrentAccountAbsence:
		return true
	case InvalidEventChannel:
		return true
	case InvalidEventToken:
		return true
	case InvalidAppletContext:
		return true
	case NoResourcePermission:
		return true
	case InvalidAccessKey:
		return true
	case AccessKeyExpired:
		return true
	case Forbidden:
		return false
	case DisabledAccount:
		return true
	case WhiteListForbidden:
		return true
	case UploadFileSizeLimit:
		return true
	case UploadFileMd5Unmatched:
		return true
	case UploadFileDiskLimit:
		return true
	case TopicAlreadySubscribed:
		return true
	case OccupiedOperator:
		return true
	case UnsupportedFSOperator:
		return true
	case NotFound:
		return true
	case ProjectNotFound:
		return true
	case ConfigNotFound:
		return true
	case ResourceNotFound:
		return true
	case AppletNotFound:
		return true
	case InstanceNotFound:
		return true
	case StrategyNotFound:
		return true
	case PublisherNotFound:
		return true
	case AccountIdentityNotFound:
		return true
	case ResourcePermNotFound:
		return true
	case CronJobNotFound:
		return true
	case InstanceNotRunning:
		return true
	case BlockchainNotFound:
		return true
	case ContractLogNotFound:
		return true
	case ChainTxNotFound:
		return true
	case ChainHeightNotFound:
		return true
	case AccountNotFound:
		return true
	case AccountPasswordNotFound:
		return true
	case OperatorNotFound:
		return true
	case TrafficLimitNotFound:
		return true
	case ProjectOperatorNotFound:
		return true
	case AccessKeyNotFound:
		return true
	case Conflict:
		return true
	case ProjectNameConflict:
		return true
	case ResourceConflict:
		return true
	case ResourceOwnerConflict:
		return true
	case StrategyConflict:
		return true
	case ConfigConflict:
		return true
	case PublisherConflict:
		return true
	case MultiInstanceDeployed:
		return true
	case AppletNameConflict:
		return true
	case CronJobConflict:
		return true
	case ContractLogConflict:
		return true
	case ChainTxConflict:
		return true
	case ChainHeightConflict:
		return true
	case AccountIdentityConflict:
		return true
	case AccountConflict:
		return true
	case AccountPasswordConflict:
		return true
	case OperatorConflict:
		return true
	case TrafficLimitConflict:
		return true
	case ProjectOperatorConflict:
		return true
	case AccessKeyNameConflict:
		return true
	case InternalServerError:
		return false
	case DatabaseError:
		return true
	case UploadFileFailed:
		return true
	case FetchResourceFailed:
		return true
	case ConfigInitFailed:
		return true
	case ConfigUninitFailed:
		return true
	case ConfigParseFailed:
		return true
	case CreateInstanceFailed:
		return true
	case BatchRemoveAppletFailed:
		return true
	case MqttSubscribeFailed:
		return true
	case MqttConnectFailed:
		return true
	case BatchRemoveWasmLogFailed:
		return true
	case GenTokenFailed:
		return true
	case TrafficLimitExceededFailed:
		return true
	case CreateTrafficSchedulerFailed:
		return true
	case UpdateTrafficSchedulerFailed:
		return true
	}
	return false
}
