// This is a generated source file. DO NOT EDIT
// Source: applet_mgr/types.go

package applet_mgr

import (
	"bytes"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
	"github.com/machinefi/w3bstream/pkg/modules/account"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/modules/cronjob"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
	"github.com/machinefi/w3bstream/pkg/modules/trafficlimit"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

type BytesBuffer = bytes.Buffer

type ChainConfigResp struct {
	Chains string `json:"chains"`
}

type CurrentAccount struct {
	GithubComMachinefiW3BstreamPkgModelsAccount
}

type CurrentPublisher struct {
	GithubComMachinefiW3BstreamPkgModelsPublisher
}

type EthClientRsp struct {
	Clients string `json:"clients"`
}

type GithubComMachinefiW3BstreamPkgDependsBaseTypesDuration = types.Duration

type GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID = types.SFID

type GithubComMachinefiW3BstreamPkgDependsBaseTypesSFIDs = types.SFIDs

type GithubComMachinefiW3BstreamPkgDependsBaseTypesTimestamp = types.Timestamp

type GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxAttachment = httpx.Attachment

type GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesBool = datatypes.Bool

type GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimes = datatypes.OperationTimes

type GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimesWithDeleted struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimes
}

type GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID = datatypes.PrimaryID

type GithubComMachinefiW3BstreamPkgDependsKitStatusxErrorField = statusx.ErrorField

type GithubComMachinefiW3BstreamPkgDependsKitStatusxErrorFields = statusx.ErrorFields

type GithubComMachinefiW3BstreamPkgDependsKitStatusxStatusErr = statusx.StatusErr

type GithubComMachinefiW3BstreamPkgEnumsAccessKeyIdentityType = enums.AccessKeyIdentityType

type GithubComMachinefiW3BstreamPkgEnumsAccessPermission = enums.AccessPermission

type GithubComMachinefiW3BstreamPkgEnumsAccountRole = enums.AccountRole

type GithubComMachinefiW3BstreamPkgEnumsAccountState = enums.AccountState

type GithubComMachinefiW3BstreamPkgEnumsCacheMode = enums.CacheMode

type GithubComMachinefiW3BstreamPkgEnumsChainName = enums.ChainName

type GithubComMachinefiW3BstreamPkgEnumsConfigType = enums.ConfigType

type GithubComMachinefiW3BstreamPkgEnumsDeployCmd = enums.DeployCmd

type GithubComMachinefiW3BstreamPkgEnumsFlowOperator = enums.FlowOperator

type GithubComMachinefiW3BstreamPkgEnumsFlowSink = enums.FlowSink

type GithubComMachinefiW3BstreamPkgEnumsInstanceState = enums.InstanceState

type GithubComMachinefiW3BstreamPkgEnumsMonitorCmd = enums.MonitorCmd

type GithubComMachinefiW3BstreamPkgEnumsOperatorKeyType = enums.OperatorKeyType

type GithubComMachinefiW3BstreamPkgEnumsProtocol = enums.Protocol

type GithubComMachinefiW3BstreamPkgEnumsTrafficLimitType = enums.TrafficLimitType

type GithubComMachinefiW3BstreamPkgEnumsWasmDBDatatype = enums.WasmDBDatatype

type GithubComMachinefiW3BstreamPkgEnumsWasmDBDialect = enums.WasmDBDialect

type GithubComMachinefiW3BstreamPkgModelsAccount struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelAccount
	GithubComMachinefiW3BstreamPkgModelsAccountInfo
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimesWithDeleted
}

type GithubComMachinefiW3BstreamPkgModelsAccountInfo = models.AccountInfo

type GithubComMachinefiW3BstreamPkgModelsApplet struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelProject
	GithubComMachinefiW3BstreamPkgModelsRelApplet
	GithubComMachinefiW3BstreamPkgModelsRelResource
	GithubComMachinefiW3BstreamPkgModelsAppletInfo
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimes
}

type GithubComMachinefiW3BstreamPkgModelsAppletInfo = models.AppletInfo

type GithubComMachinefiW3BstreamPkgModelsChainHeight struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelChainHeight
	GithubComMachinefiW3BstreamPkgModelsChainHeightData
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimes
}

type GithubComMachinefiW3BstreamPkgModelsChainHeightData struct {
	GithubComMachinefiW3BstreamPkgModelsChainHeightInfo
	ProjectName string `json:"projectName"`
}

type GithubComMachinefiW3BstreamPkgModelsChainHeightInfo = models.ChainHeightInfo

type GithubComMachinefiW3BstreamPkgModelsChainTx struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelChainTx
	GithubComMachinefiW3BstreamPkgModelsChainTxData
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimes
}

type GithubComMachinefiW3BstreamPkgModelsChainTxData struct {
	GithubComMachinefiW3BstreamPkgModelsChainTxInfo
	ProjectName string `json:"projectName"`
}

type GithubComMachinefiW3BstreamPkgModelsChainTxInfo = models.ChainTxInfo

type GithubComMachinefiW3BstreamPkgModelsConfig struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelConfig
	GithubComMachinefiW3BstreamPkgModelsConfigBase
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimes
}

type GithubComMachinefiW3BstreamPkgModelsConfigBase = models.ConfigBase

type GithubComMachinefiW3BstreamPkgModelsContractLog struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelContractLog
	GithubComMachinefiW3BstreamPkgModelsContractLogData
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimes
}

type GithubComMachinefiW3BstreamPkgModelsContractLogData struct {
	GithubComMachinefiW3BstreamPkgModelsContractLogInfo
	ProjectName string `json:"projectName"`
}

type GithubComMachinefiW3BstreamPkgModelsContractLogInfo = models.ContractLogInfo

type GithubComMachinefiW3BstreamPkgModelsCronJob struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelCronJob
	GithubComMachinefiW3BstreamPkgModelsRelProject
	GithubComMachinefiW3BstreamPkgModelsCronJobInfo
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimesWithDeleted
}

type GithubComMachinefiW3BstreamPkgModelsCronJobInfo = models.CronJobInfo

type GithubComMachinefiW3BstreamPkgModelsInstance struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelInstance
	GithubComMachinefiW3BstreamPkgModelsRelApplet
	GithubComMachinefiW3BstreamPkgModelsInstanceInfo
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimes
}

type GithubComMachinefiW3BstreamPkgModelsInstanceInfo = models.InstanceInfo

type GithubComMachinefiW3BstreamPkgModelsMeta = models.Meta

type GithubComMachinefiW3BstreamPkgModelsOperator struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelAccount
	GithubComMachinefiW3BstreamPkgModelsRelOperator
	GithubComMachinefiW3BstreamPkgModelsOperatorInfo
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimesWithDeleted
}

type GithubComMachinefiW3BstreamPkgModelsOperatorInfo = models.OperatorInfo

type GithubComMachinefiW3BstreamPkgModelsProject struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelProject
	GithubComMachinefiW3BstreamPkgModelsRelAccount
	GithubComMachinefiW3BstreamPkgModelsProjectName
	GithubComMachinefiW3BstreamPkgModelsProjectBase
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimesWithDeleted
}

type GithubComMachinefiW3BstreamPkgModelsProjectBase = models.ProjectBase

type GithubComMachinefiW3BstreamPkgModelsProjectName = models.ProjectName

type GithubComMachinefiW3BstreamPkgModelsProjectOperator struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelProject
	GithubComMachinefiW3BstreamPkgModelsRelOperator
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimesWithDeleted
}

type GithubComMachinefiW3BstreamPkgModelsPublisher struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelProject
	GithubComMachinefiW3BstreamPkgModelsRelPublisher
	GithubComMachinefiW3BstreamPkgModelsPublisherInfo
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimes
}

type GithubComMachinefiW3BstreamPkgModelsPublisherInfo = models.PublisherInfo

type GithubComMachinefiW3BstreamPkgModelsRelAccount = models.RelAccount

type GithubComMachinefiW3BstreamPkgModelsRelApplet = models.RelApplet

type GithubComMachinefiW3BstreamPkgModelsRelChainHeight = models.RelChainHeight

type GithubComMachinefiW3BstreamPkgModelsRelChainTx = models.RelChainTx

type GithubComMachinefiW3BstreamPkgModelsRelConfig = models.RelConfig

type GithubComMachinefiW3BstreamPkgModelsRelContractLog = models.RelContractLog

type GithubComMachinefiW3BstreamPkgModelsRelCronJob = models.RelCronJob

type GithubComMachinefiW3BstreamPkgModelsRelInstance = models.RelInstance

type GithubComMachinefiW3BstreamPkgModelsRelOperator = models.RelOperator

type GithubComMachinefiW3BstreamPkgModelsRelProject = models.RelProject

type GithubComMachinefiW3BstreamPkgModelsRelPublisher = models.RelPublisher

type GithubComMachinefiW3BstreamPkgModelsRelResource = models.RelResource

type GithubComMachinefiW3BstreamPkgModelsRelStrategy = models.RelStrategy

type GithubComMachinefiW3BstreamPkgModelsRelTrafficLimit = models.RelTrafficLimit

type GithubComMachinefiW3BstreamPkgModelsResource struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelResource
	GithubComMachinefiW3BstreamPkgModelsResourceInfo
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimes
}

type GithubComMachinefiW3BstreamPkgModelsResourceInfo = models.ResourceInfo

type GithubComMachinefiW3BstreamPkgModelsResourceOwnerInfo = models.ResourceOwnerInfo

type GithubComMachinefiW3BstreamPkgModelsStrategy struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelStrategy
	GithubComMachinefiW3BstreamPkgModelsRelProject
	GithubComMachinefiW3BstreamPkgModelsRelApplet
	GithubComMachinefiW3BstreamPkgModelsStrategyInfo
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimesWithDeleted
}

type GithubComMachinefiW3BstreamPkgModelsStrategyInfo = models.StrategyInfo

type GithubComMachinefiW3BstreamPkgModelsTrafficLimit struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesPrimaryID
	GithubComMachinefiW3BstreamPkgModelsRelTrafficLimit
	GithubComMachinefiW3BstreamPkgModelsRelProject
	GithubComMachinefiW3BstreamPkgModelsTrafficLimitInfo
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimesWithDeleted
}

type GithubComMachinefiW3BstreamPkgModelsTrafficLimitInfo = models.TrafficLimitInfo

type GithubComMachinefiW3BstreamPkgModulesAccessKeyCreateReqBase = access_key.CreateReqBase

type GithubComMachinefiW3BstreamPkgModulesAccessKeyCreateRsp = access_key.CreateRsp

type GithubComMachinefiW3BstreamPkgModulesAccessKeyGroupAccessPrivilege = access_key.GroupAccessPrivilege

type GithubComMachinefiW3BstreamPkgModulesAccessKeyGroupAccessPrivileges = access_key.GroupAccessPrivileges

type GithubComMachinefiW3BstreamPkgModulesAccessKeyGroupMetaBase = access_key.GroupMetaBase

type GithubComMachinefiW3BstreamPkgModulesAccessKeyGroupMetaWithPrivilege struct {
	GithubComMachinefiW3BstreamPkgModulesAccessKeyGroupMetaBase
	Perm GithubComMachinefiW3BstreamPkgEnumsAccessPermission `json:"perm"`
}

type GithubComMachinefiW3BstreamPkgModulesAccessKeyListData struct {
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimes
	Desc       string                                                                 `json:"desc,omitempty"`
	ExpiredAt  *GithubComMachinefiW3BstreamPkgDependsBaseTypesTimestamp               `json:"expiredAt,omitempty"`
	LastUsed   *GithubComMachinefiW3BstreamPkgDependsBaseTypesTimestamp               `json:"lastUsed,omitempty"`
	Name       string                                                                 `json:"name"`
	Privileges []GithubComMachinefiW3BstreamPkgModulesAccessKeyGroupMetaWithPrivilege `json:"privileges,omitempty"`
}

type GithubComMachinefiW3BstreamPkgModulesAccessKeyListRsp = access_key.ListRsp

type GithubComMachinefiW3BstreamPkgModulesAccessKeyUpdateReq = access_key.UpdateReq

type GithubComMachinefiW3BstreamPkgModulesAccessKeyUpdateRsp = access_key.UpdateRsp

type GithubComMachinefiW3BstreamPkgModulesAccountCreateAccountByUsernameReq = account.CreateAccountByUsernameReq

type GithubComMachinefiW3BstreamPkgModulesAccountCreateAccountByUsernameRsp struct {
	GithubComMachinefiW3BstreamPkgModelsAccount
	Password string `json:"password"`
}

type GithubComMachinefiW3BstreamPkgModulesAccountLoginByEthAddressReq = account.LoginByEthAddressReq

type GithubComMachinefiW3BstreamPkgModulesAccountLoginByUsernameReq = account.LoginByUsernameReq

type GithubComMachinefiW3BstreamPkgModulesAccountLoginRsp = account.LoginRsp

type GithubComMachinefiW3BstreamPkgModulesAccountUpdatePasswordReq = account.UpdatePasswordReq

type GithubComMachinefiW3BstreamPkgModulesAppletCreateReq = applet.CreateReq

type GithubComMachinefiW3BstreamPkgModulesAppletCreateRsp struct {
	GithubComMachinefiW3BstreamPkgModelsApplet
	Instance   *GithubComMachinefiW3BstreamPkgModelsInstance  `json:"instance"`
	Resource   *GithubComMachinefiW3BstreamPkgModelsResource  `json:"resource,omitempty"`
	Strategies []GithubComMachinefiW3BstreamPkgModelsStrategy `json:"strategies,omitempty"`
}

type GithubComMachinefiW3BstreamPkgModulesAppletDetail struct {
	GithubComMachinefiW3BstreamPkgModelsApplet
	GithubComMachinefiW3BstreamPkgModelsResourceInfo
	GithubComMachinefiW3BstreamPkgModelsRelInstance
	GithubComMachinefiW3BstreamPkgModelsInstanceInfo
}

type GithubComMachinefiW3BstreamPkgModulesAppletInfo = applet.Info

type GithubComMachinefiW3BstreamPkgModulesAppletListRsp = applet.ListRsp

type GithubComMachinefiW3BstreamPkgModulesAppletUpdateReq = applet.UpdateReq

type GithubComMachinefiW3BstreamPkgModulesBlockchainBatchUpdateMonitorReq = blockchain.BatchUpdateMonitorReq

type GithubComMachinefiW3BstreamPkgModulesBlockchainCreateChainHeightReq struct {
	GithubComMachinefiW3BstreamPkgModelsChainHeightInfo
}

type GithubComMachinefiW3BstreamPkgModulesBlockchainCreateChainTxReq struct {
	GithubComMachinefiW3BstreamPkgModelsChainTxInfo
}

type GithubComMachinefiW3BstreamPkgModulesBlockchainCreateContractLogReq struct {
	GithubComMachinefiW3BstreamPkgModelsContractLogInfo
}

type GithubComMachinefiW3BstreamPkgModulesCronjobCreateReq struct {
	GithubComMachinefiW3BstreamPkgModelsCronJobInfo
}

type GithubComMachinefiW3BstreamPkgModulesCronjobListRsp = cronjob.ListRsp

type GithubComMachinefiW3BstreamPkgModulesDeployCreateReq = deploy.CreateReq

type GithubComMachinefiW3BstreamPkgModulesEventEventRsp = event.EventRsp

type GithubComMachinefiW3BstreamPkgModulesEventResult = event.Result

type GithubComMachinefiW3BstreamPkgModulesOperatorCreateReq = operator.CreateReq

type GithubComMachinefiW3BstreamPkgModulesOperatorDetail struct {
	GithubComMachinefiW3BstreamPkgModelsOperator
	Address string `json:"address"`
}

type GithubComMachinefiW3BstreamPkgModulesOperatorListDetailRsp = operator.ListDetailRsp

type GithubComMachinefiW3BstreamPkgModulesProjectCreateReq struct {
	GithubComMachinefiW3BstreamPkgModelsProjectName
	GithubComMachinefiW3BstreamPkgModelsProjectBase
	Database *GithubComMachinefiW3BstreamPkgTypesWasmDatabase `json:"database,omitempty"`
	Env      *GithubComMachinefiW3BstreamPkgTypesWasmEnv      `json:"envs,omitempty"`
	Flow     *GithubComMachinefiW3BstreamPkgTypesWasmFlow     `json:"flow,omitempty"`
}

type GithubComMachinefiW3BstreamPkgModulesProjectCreateRsp struct {
	GithubComMachinefiW3BstreamPkgModelsProject
	ChannelState GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesBool `json:"channelState"`
	Database     *GithubComMachinefiW3BstreamPkgTypesWasmDatabase          `json:"database,omitempty"`
	Env          *GithubComMachinefiW3BstreamPkgTypesWasmEnv               `json:"envs,omitempty"`
	Flow         *GithubComMachinefiW3BstreamPkgTypesWasmFlow              `json:"flow,omitempty"`
}

type GithubComMachinefiW3BstreamPkgModulesProjectDetail = project.Detail

type GithubComMachinefiW3BstreamPkgModulesProjectListDetailRsp = project.ListDetailRsp

type GithubComMachinefiW3BstreamPkgModulesProjectListRsp = project.ListRsp

type GithubComMachinefiW3BstreamPkgModulesPublisherCreateReq = publisher.CreateReq

type GithubComMachinefiW3BstreamPkgModulesPublisherListRsp = publisher.ListRsp

type GithubComMachinefiW3BstreamPkgModulesPublisherUpdateReq = publisher.UpdateReq

type GithubComMachinefiW3BstreamPkgModulesResourceDownLoadResourceRsp = resource.DownLoadResourceRsp

type GithubComMachinefiW3BstreamPkgModulesResourceListRsp = resource.ListRsp

type GithubComMachinefiW3BstreamPkgModulesResourceResourceInfo struct {
	GithubComMachinefiW3BstreamPkgModelsRelResource
	GithubComMachinefiW3BstreamPkgModelsResourceInfo
	GithubComMachinefiW3BstreamPkgModelsResourceOwnerInfo
	GithubComMachinefiW3BstreamPkgDependsKitSqlxDatatypesOperationTimes
}

type GithubComMachinefiW3BstreamPkgModulesStrategyCreateReq struct {
	GithubComMachinefiW3BstreamPkgModelsRelApplet
	GithubComMachinefiW3BstreamPkgModelsStrategyInfo
}

type GithubComMachinefiW3BstreamPkgModulesStrategyListRsp = strategy.ListRsp

type GithubComMachinefiW3BstreamPkgModulesTrafficlimitCreateReq = trafficlimit.CreateReq

type GithubComMachinefiW3BstreamPkgModulesTrafficlimitListRsp = trafficlimit.ListRsp

type GithubComMachinefiW3BstreamPkgModulesTrafficlimitUpdateReq = trafficlimit.UpdateReq

type GithubComMachinefiW3BstreamPkgTypesWasmCBInfo = wasm.CBInfo

type GithubComMachinefiW3BstreamPkgTypesWasmCache = wasm.Cache

type GithubComMachinefiW3BstreamPkgTypesWasmColumn = wasm.Column

type GithubComMachinefiW3BstreamPkgTypesWasmConstrains = wasm.Constrains

type GithubComMachinefiW3BstreamPkgTypesWasmDBInfo = wasm.DBInfo

type GithubComMachinefiW3BstreamPkgTypesWasmDatabase = wasm.Database

type GithubComMachinefiW3BstreamPkgTypesWasmEnv = wasm.Env

type GithubComMachinefiW3BstreamPkgTypesWasmFlow = wasm.Flow

type GithubComMachinefiW3BstreamPkgTypesWasmKey = wasm.Key

type GithubComMachinefiW3BstreamPkgTypesWasmOperator = wasm.Operator

type GithubComMachinefiW3BstreamPkgTypesWasmSchema = wasm.Schema

type GithubComMachinefiW3BstreamPkgTypesWasmSink = wasm.Sink

type GithubComMachinefiW3BstreamPkgTypesWasmSinkInfo = wasm.SinkInfo

type GithubComMachinefiW3BstreamPkgTypesWasmSource = wasm.Source

type GithubComMachinefiW3BstreamPkgTypesWasmTable = wasm.Table

type ProjectProvider struct {
	ProjectName string `name:"projectName" validate:"@projectName"`
}
