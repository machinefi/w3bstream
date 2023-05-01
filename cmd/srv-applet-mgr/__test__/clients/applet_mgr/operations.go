// This is a generated source file. DO NOT EDIT
// Source: applet_mgr/operations.go

package applet_mgr

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
)

type BatchRemoveApplet struct {
	ProjectName  string                                               `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string                                               `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AppletIDs    []GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"query" name:"appletID,omitempty"`
	AuthInQuery  string                                               `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	LNameLike    string                                               `in:"query" name:"lName,omitempty"`
	NameLike     string                                               `in:"query" name:"name,omitempty"`
	Names        []string                                             `in:"query" name:"names,omitempty"`
}

func (o *BatchRemoveApplet) Path() string {
	return "/srv-applet-mgr/v0/applet/x/:projectName"
}

func (o *BatchRemoveApplet) Method() string {
	return "DELETE"
}

func (o *BatchRemoveApplet) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.BatchRemoveApplet")
	return cli.Do(ctx, o, metas...)
}

func (o *BatchRemoveApplet) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *BatchRemoveApplet) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type BatchRemoveInstance struct {
	ProjectName  string                                               `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string                                               `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AppletIDs    []GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"query" name:"appletID,omitempty"`
	AuthInQuery  string                                               `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	InstanceIDs  []GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"query" name:"instanceID,omitempty"`
	States       []GithubComMachinefiW3BstreamPkgEnumsInstanceState   `in:"query" name:"state,omitempty"`
}

func (o *BatchRemoveInstance) Path() string {
	return "/srv-applet-mgr/v0/deploy/x/:projectName"
}

func (o *BatchRemoveInstance) Method() string {
	return "DELETE"
}

func (o *BatchRemoveInstance) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.BatchRemoveInstance")
	return cli.Do(ctx, o, metas...)
}

func (o *BatchRemoveInstance) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *BatchRemoveInstance) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type BatchRemovePublisher struct {
	ProjectName  string                                               `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string                                               `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                               `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	Keys         []string                                             `in:"query" name:"key,omitempty"`
	LNameLike    string                                               `in:"query" name:"lname,omitempty"`
	NameLike     string                                               `in:"query" name:"name,omitempty"`
	PublisherIDs []GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"query" name:"publisherIDs,omitempty"`
	RNameLike    string                                               `in:"query" name:"rname,omitempty"`
}

func (o *BatchRemovePublisher) Path() string {
	return "/srv-applet-mgr/v0/publisher/x/:projectName"
}

func (o *BatchRemovePublisher) Method() string {
	return "DELETE"
}

func (o *BatchRemovePublisher) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.BatchRemovePublisher")
	return cli.Do(ctx, o, metas...)
}

func (o *BatchRemovePublisher) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *BatchRemovePublisher) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type BatchRemoveStrategy struct {
	ProjectName  string                                               `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string                                               `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AppletIDs    []GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"query" name:"appletID,omitempty"`
	AuthInQuery  string                                               `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	EventTypes   []string                                             `in:"query" name:"eventType,omitempty"`
	Handlers     []string                                             `in:"query" name:"handler,omitempty"`
	StrategyIDs  []GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"query" name:"strategyID,omitempty"`
}

func (o *BatchRemoveStrategy) Path() string {
	return "/srv-applet-mgr/v0/strategy/x/:projectName"
}

func (o *BatchRemoveStrategy) Method() string {
	return "DELETE"
}

func (o *BatchRemoveStrategy) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.BatchRemoveStrategy")
	return cli.Do(ctx, o, metas...)
}

func (o *BatchRemoveStrategy) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *BatchRemoveStrategy) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type ControlInstance struct {
	Cmd          GithubComMachinefiW3BstreamPkgEnumsDeployCmd       `in:"path" name:"cmd"`
	InstanceID   GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"instanceID"`
	AuthInHeader string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *ControlInstance) Path() string {
	return "/srv-applet-mgr/v0/deploy/:instanceID/:cmd"
}

func (o *ControlInstance) Method() string {
	return "PUT"
}

func (o *ControlInstance) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.ControlInstance")
	return cli.Do(ctx, o, metas...)
}

func (o *ControlInstance) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *ControlInstance) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type CreateAccountByUsernameAndPassword struct {
	CreateAccountByUsernameReq GithubComMachinefiW3BstreamPkgModulesAccountCreateAccountByUsernameReq `in:"body"`
}

func (o *CreateAccountByUsernameAndPassword) Path() string {
	return "/srv-applet-mgr/v0/register/admin"
}

func (o *CreateAccountByUsernameAndPassword) Method() string {
	return "POST"
}

// @StatusErr[AccountConflict][409999015][Account Conflict]!
// @StatusErr[AccountIdentityConflict][409999014][Account Identity Conflict]!
// @StatusErr[AccountPasswordConflict][409999016][Account Password Conflict]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[NoAdminPermission][401999005][No Admin Permission]!

func (o *CreateAccountByUsernameAndPassword) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.CreateAccountByUsernameAndPassword")
	return cli.Do(ctx, o, metas...)
}

func (o *CreateAccountByUsernameAndPassword) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccountCreateAccountByUsernameRsp, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModulesAccountCreateAccountByUsernameRsp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *CreateAccountByUsernameAndPassword) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccountCreateAccountByUsernameRsp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type CreateApplet struct {
	ProjectName  string                                               `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string                                               `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                               `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	CreateReq    GithubComMachinefiW3BstreamPkgModulesAppletCreateReq `in:"body" mime:"multipart"`
}

func (o *CreateApplet) Path() string {
	return "/srv-applet-mgr/v0/applet/x/:projectName"
}

func (o *CreateApplet) Method() string {
	return "POST"
}

// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[AppletNameConflict][409999009][Applet Name Conflict]!
// @StatusErr[ConfigConflict][409999006][Config Conflict]!
// @StatusErr[ConfigInitFailed][500999006][Config Init Failed]!
// @StatusErr[ConfigParseFailed][500999008][Config Parse Failed]!
// @StatusErr[ConfigUninitFailed][500999007][Config Uninit Failed]!
// @StatusErr[CreateInstanceFailed][500999010][Create Instance Failed]!
// @StatusErr[CurrentAccountAbsence][401999012][Current Account Absence]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InternalServerError][500999001][internal error]
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidConfigType][400999002][Invalid Config Type]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[MD5ChecksumFailed][500999012][Md5 Checksum Failed]!
// @StatusErr[MultiInstanceDeployed][409999008][Multi Instance Deployed]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!
// @StatusErr[ResourceConflict][409999003][Resource Conflict]!
// @StatusErr[StrategyConflict][409999005][Strategy Conflict]!
// @StatusErr[UploadFileDiskLimit][403999006][Upload File Disk Limit]!
// @StatusErr[UploadFileFailed][500999003][Upload File Failed]!
// @StatusErr[UploadFileMd5Unmatched][403999005][Upload File Md5 Unmatched]!
// @StatusErr[UploadFileSizeLimit][403999004][Upload File Size Limit]!

func (o *CreateApplet) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.CreateApplet")
	return cli.Do(ctx, o, metas...)
}

func (o *CreateApplet) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAppletCreateRsp, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModulesAppletCreateRsp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *CreateApplet) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAppletCreateRsp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type CreateChainHeight struct {
	ProjectName          string                                                              `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader         string                                                              `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery          string                                                              `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	CreateChainHeightReq GithubComMachinefiW3BstreamPkgModulesBlockchainCreateChainHeightReq `in:"body"`
}

func (o *CreateChainHeight) Path() string {
	return "/srv-applet-mgr/v0/monitor/x/:projectName/chain_height"
}

func (o *CreateChainHeight) Method() string {
	return "POST"
}

// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[BlockchainNotFound][404999013][Blockchain Not Found]!
// @StatusErr[ChainHeightConflict][409999013][Chain Height Conflict]!
// @StatusErr[CurrentAccountAbsence][401999012][Current Account Absence]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!

func (o *CreateChainHeight) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.CreateChainHeight")
	return cli.Do(ctx, o, metas...)
}

func (o *CreateChainHeight) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsChainHeight, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModelsChainHeight)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *CreateChainHeight) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsChainHeight, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type CreateChainTx struct {
	ProjectName      string                                                          `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader     string                                                          `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery      string                                                          `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	CreateChainTxReq GithubComMachinefiW3BstreamPkgModulesBlockchainCreateChainTxReq `in:"body"`
}

func (o *CreateChainTx) Path() string {
	return "/srv-applet-mgr/v0/monitor/x/:projectName/chain_tx"
}

func (o *CreateChainTx) Method() string {
	return "POST"
}

// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[BlockchainNotFound][404999013][Blockchain Not Found]!
// @StatusErr[ChainTxConflict][409999012][Chain Tx Conflict]!
// @StatusErr[CurrentAccountAbsence][401999012][Current Account Absence]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!

func (o *CreateChainTx) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.CreateChainTx")
	return cli.Do(ctx, o, metas...)
}

func (o *CreateChainTx) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsChainTx, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModelsChainTx)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *CreateChainTx) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsChainTx, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type CreateContractLog struct {
	ProjectName          string                                                              `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader         string                                                              `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery          string                                                              `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	CreateContractLogReq GithubComMachinefiW3BstreamPkgModulesBlockchainCreateContractLogReq `in:"body"`
}

func (o *CreateContractLog) Path() string {
	return "/srv-applet-mgr/v0/monitor/x/:projectName/contract_log"
}

func (o *CreateContractLog) Method() string {
	return "POST"
}

// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[BlockchainNotFound][404999013][Blockchain Not Found]!
// @StatusErr[ContractLogConflict][409999011][Contract Log Conflict]!
// @StatusErr[CurrentAccountAbsence][401999012][Current Account Absence]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!

func (o *CreateContractLog) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.CreateContractLog")
	return cli.Do(ctx, o, metas...)
}

func (o *CreateContractLog) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsContractLog, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModelsContractLog)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *CreateContractLog) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsContractLog, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type CreateCronJob struct {
	ProjectID    GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID    `in:"path" name:"projectID"`
	AuthInHeader string                                                `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                                `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	CreateReq    GithubComMachinefiW3BstreamPkgModulesCronjobCreateReq `in:"body"`
}

func (o *CreateCronJob) Path() string {
	return "/srv-applet-mgr/v0/cronjob/:projectID"
}

func (o *CreateCronJob) Method() string {
	return "POST"
}

// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[CronJobConflict][409999010][Cron Job Conflict]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidCronExpressions][400999005][Invalid Cron Expressions]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!

func (o *CreateCronJob) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.CreateCronJob")
	return cli.Do(ctx, o, metas...)
}

func (o *CreateCronJob) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsCronJob, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModelsCronJob)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *CreateCronJob) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsCronJob, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type CreateInstance struct {
	AppletID     GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID   `in:"path" name:"appletID"`
	AuthInHeader string                                               `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                               `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	CreateReq    GithubComMachinefiW3BstreamPkgModulesDeployCreateReq `in:"body"`
}

func (o *CreateInstance) Path() string {
	return "/srv-applet-mgr/v0/deploy/applet/:appletID"
}

func (o *CreateInstance) Method() string {
	return "POST"
}

// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[AppletNotFound][404999005][Applet Not Found]!
// @StatusErr[ConfigConflict][409999006][Config Conflict]!
// @StatusErr[ConfigInitFailed][500999006][Config Init Failed]!
// @StatusErr[ConfigParseFailed][500999008][Config Parse Failed]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InstanceNotFound][404999006][Instance Not Found]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[MultiInstanceDeployed][409999008][Multi Instance Deployed]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!
// @StatusErr[ResourceNotFound][404999004][Resource Not Found]!

func (o *CreateInstance) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.CreateInstance")
	return cli.Do(ctx, o, metas...)
}

func (o *CreateInstance) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsInstance, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModelsInstance)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *CreateInstance) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsInstance, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type CreateProject struct {
	AuthInHeader string                                                `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                                `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	CreateReq    GithubComMachinefiW3BstreamPkgModulesProjectCreateReq `in:"body"`
}

func (o *CreateProject) Path() string {
	return "/srv-applet-mgr/v0/project"
}

func (o *CreateProject) Method() string {
	return "POST"
}

// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[ClientClosedRequest][499000000][ClientClosedRequest]
// @StatusErr[ConfigConflict][409999006][Config Conflict]!
// @StatusErr[ConfigInitFailed][500999006][Config Init Failed]!
// @StatusErr[ConfigParseFailed][500999008][Config Parse Failed]!
// @StatusErr[CurrentAccountAbsence][401999012][Current Account Absence]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidEventToken][401999014][Invalid Event Token]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[MqttConnectFailed][500999014][MQTT Connect Failed]!
// @StatusErr[MqttSubscribeFailed][500999013][MQTT Subscribe Failed]!
// @StatusErr[ProjectNameConflict][409999002][Project Name Conflict]!
// @StatusErr[RequestFailed][500000000][RequestFailed]
// @StatusErr[RequestTransformFailed][400000000][RequestTransformFailed]
// @StatusErr[TopicAlreadySubscribed][403999007][Topic Already Subscribed]!

func (o *CreateProject) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.CreateProject")
	return cli.Do(ctx, o, metas...)
}

func (o *CreateProject) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesProjectCreateRsp, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModulesProjectCreateRsp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *CreateProject) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesProjectCreateRsp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type CreateProjectSchema struct {
	ProjectName  string                                          `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string                                          `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                          `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	Database     GithubComMachinefiW3BstreamPkgTypesWasmDatabase `in:"body"`
}

func (o *CreateProjectSchema) Path() string {
	return "/srv-applet-mgr/v0/project_config/x/:projectName"
}

func (o *CreateProjectSchema) Method() string {
	return "POST"
}

// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[ConfigConflict][409999006][Config Conflict]!
// @StatusErr[ConfigInitFailed][500999006][Config Init Failed]!
// @StatusErr[ConfigParseFailed][500999008][Config Parse Failed]!
// @StatusErr[CurrentAccountAbsence][401999012][Current Account Absence]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!

func (o *CreateProjectSchema) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.CreateProjectSchema")
	return cli.Do(ctx, o, metas...)
}

func (o *CreateProjectSchema) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsConfig, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModelsConfig)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *CreateProjectSchema) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsConfig, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type CreatePublisher struct {
	ProjectName  string                                                  `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string                                                  `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                                  `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	CreateReq    GithubComMachinefiW3BstreamPkgModulesPublisherCreateReq `in:"body"`
}

func (o *CreatePublisher) Path() string {
	return "/srv-applet-mgr/v0/publisher/x/:projectName"
}

func (o *CreatePublisher) Method() string {
	return "POST"
}

// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[CurrentAccountAbsence][401999012][Current Account Absence]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[GenPublisherTokenFailed][500999009][Gen Publisher Token Failed]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!
// @StatusErr[PublisherConflict][409999007][Publisher Conflict]!

func (o *CreatePublisher) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.CreatePublisher")
	return cli.Do(ctx, o, metas...)
}

func (o *CreatePublisher) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsPublisher, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModelsPublisher)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *CreatePublisher) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsPublisher, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type CreateStrategy struct {
	ProjectName  string                                                 `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string                                                 `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                                 `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	CreateReq    GithubComMachinefiW3BstreamPkgModulesStrategyCreateReq `in:"body"`
}

func (o *CreateStrategy) Path() string {
	return "/srv-applet-mgr/v0/strategy/x/:projectName"
}

func (o *CreateStrategy) Method() string {
	return "POST"
}

// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[AppletNotFound][404999005][Applet Not Found]!
// @StatusErr[CurrentAccountAbsence][401999012][Current Account Absence]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InstanceNotFound][404999006][Instance Not Found]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!
// @StatusErr[ResourceNotFound][404999004][Resource Not Found]!
// @StatusErr[StrategyConflict][409999005][Strategy Conflict]!

func (o *CreateStrategy) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.CreateStrategy")
	return cli.Do(ctx, o, metas...)
}

func (o *CreateStrategy) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsStrategy, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModelsStrategy)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *CreateStrategy) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsStrategy, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type GetApplet struct {
	AppletID     GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"appletID"`
	AuthInHeader string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *GetApplet) Path() string {
	return "/srv-applet-mgr/v0/applet/data/:appletID"
}

func (o *GetApplet) Method() string {
	return "GET"
}

// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[AppletNotFound][404999005][Applet Not Found]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InstanceNotFound][404999006][Instance Not Found]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!
// @StatusErr[ResourceNotFound][404999004][Resource Not Found]!

func (o *GetApplet) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.GetApplet")
	return cli.Do(ctx, o, metas...)
}

func (o *GetApplet) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsApplet, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModelsApplet)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *GetApplet) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsApplet, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type GetInstanceByAppletID struct {
	AppletID     GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"appletID"`
	AuthInHeader string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *GetInstanceByAppletID) Path() string {
	return "/srv-applet-mgr/v0/deploy/applet/:appletID"
}

func (o *GetInstanceByAppletID) Method() string {
	return "GET"
}

// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[AppletNotFound][404999005][Applet Not Found]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InstanceNotFound][404999006][Instance Not Found]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!
// @StatusErr[ResourceNotFound][404999004][Resource Not Found]!

func (o *GetInstanceByAppletID) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.GetInstanceByAppletID")
	return cli.Do(ctx, o, metas...)
}

func (o *GetInstanceByAppletID) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsInstance, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModelsInstance)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *GetInstanceByAppletID) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsInstance, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type GetInstanceByInstanceID struct {
	InstanceID   GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"instanceID"`
	AuthInHeader string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *GetInstanceByInstanceID) Path() string {
	return "/srv-applet-mgr/v0/deploy/instance/:instanceID"
}

func (o *GetInstanceByInstanceID) Method() string {
	return "GET"
}

// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[AppletNotFound][404999005][Applet Not Found]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InstanceNotFound][404999006][Instance Not Found]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!
// @StatusErr[ResourceNotFound][404999004][Resource Not Found]!

func (o *GetInstanceByInstanceID) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.GetInstanceByInstanceID")
	return cli.Do(ctx, o, metas...)
}

func (o *GetInstanceByInstanceID) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsInstance, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModelsInstance)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *GetInstanceByInstanceID) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsInstance, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type GetOperatorAddr struct {
	AuthInHeader string `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *GetOperatorAddr) Path() string {
	return "/srv-applet-mgr/v0/account/operatoraddr"
}

func (o *GetOperatorAddr) Method() string {
	return "GET"
}

// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!

func (o *GetOperatorAddr) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.GetOperatorAddr")
	return cli.Do(ctx, o, metas...)
}

func (o *GetOperatorAddr) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*string, kit.Metadata, error) {
	rsp := new(string)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *GetOperatorAddr) Invoke(cli kit.Client, metas ...kit.Metadata) (*string, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type GetProject struct {
	ProjectName  string `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *GetProject) Path() string {
	return "/srv-applet-mgr/v0/project/x/:projectName/data"
}

func (o *GetProject) Method() string {
	return "GET"
}

// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[CurrentAccountAbsence][401999012][Current Account Absence]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DeprecatedProject][400999003][Deprecated Project]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!

func (o *GetProject) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.GetProject")
	return cli.Do(ctx, o, metas...)
}

func (o *GetProject) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsProject, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModelsProject)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *GetProject) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsProject, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type GetProjectSchema struct {
	ProjectName  string `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *GetProjectSchema) Path() string {
	return "/srv-applet-mgr/v0/project_config/x/:projectName"
}

func (o *GetProjectSchema) Method() string {
	return "GET"
}

// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[ConfigNotFound][404999003][Config Not Found]!
// @StatusErr[ConfigParseFailed][500999008][Config Parse Failed]!
// @StatusErr[CurrentAccountAbsence][401999012][Current Account Absence]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidConfigType][400999002][Invalid Config Type]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!

func (o *GetProjectSchema) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.GetProjectSchema")
	return cli.Do(ctx, o, metas...)
}

func (o *GetProjectSchema) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgTypesWasmConfiguration, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgTypesWasmConfiguration)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *GetProjectSchema) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgTypesWasmConfiguration, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type GetPublisher struct {
	PublisherID  GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"publisherID"`
	AuthInHeader string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *GetPublisher) Path() string {
	return "/srv-applet-mgr/v0/publisher/data/:publisherID"
}

func (o *GetPublisher) Method() string {
	return "GET"
}

// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!
// @StatusErr[PublisherNotFound][404999008][Publisher Not Found]!

func (o *GetPublisher) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.GetPublisher")
	return cli.Do(ctx, o, metas...)
}

func (o *GetPublisher) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsPublisher, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModelsPublisher)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *GetPublisher) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsPublisher, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type GetStrategy struct {
	StrategyID   GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"strategyID"`
	AuthInHeader string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *GetStrategy) Path() string {
	return "/srv-applet-mgr/v0/strategy/data/:strategyID"
}

func (o *GetStrategy) Method() string {
	return "GET"
}

// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!
// @StatusErr[StrategyNotFound][404999007][Strategy Not Found]!

func (o *GetStrategy) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.GetStrategy")
	return cli.Do(ctx, o, metas...)
}

func (o *GetStrategy) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsStrategy, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModelsStrategy)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *GetStrategy) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsStrategy, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type HandleEvent struct {
	// Channel message channel named (intact project name)
	Channel      string `in:"path" name:"channel"`
	AuthInHeader string `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	// EventID unique id for tracing event under channel
	EventID string `in:"query" name:"eventID,omitempty"`
	// EventType used for filter strategies created in w3b before
	EventType string `in:"query" name:"eventType,omitempty"`
	// Timestamp event time when publisher do send
	Timestamp int64 `in:"query" name:"timestamp,omitempty"`
	// Payload event payload (binary only)
	Payload BytesBuffer `in:"body" mime:"stream"`
}

func (o *HandleEvent) Path() string {
	return "/srv-applet-mgr/v0/event/:channel"
}

func (o *HandleEvent) Method() string {
	return "POST"
}

// @StatusErr[AccountConflict][409999015][Account Conflict]!
// @StatusErr[AccountIdentityConflict][409999014][Account Identity Conflict]!
// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[AccountPasswordConflict][409999016][Account Password Conflict]!
// @StatusErr[AccountPasswordNotFound][404999018][Account Password Not Found]!
// @StatusErr[AppletNameConflict][409999009][Applet Name Conflict]!
// @StatusErr[AppletNotFound][404999005][Applet Not Found]!
// @StatusErr[BadRequest][400999001][BadRequest]!
// @StatusErr[BatchRemoveAppletFailed][500999011][Batch Remove Applet Failed]!
// @StatusErr[BlockchainNotFound][404999013][Blockchain Not Found]!
// @StatusErr[ChainHeightConflict][409999013][Chain Height Conflict]!
// @StatusErr[ChainHeightNotFound][404999016][Chain Height Not Found]!
// @StatusErr[ChainTxConflict][409999012][Chain Tx Conflict]!
// @StatusErr[ChainTxNotFound][404999015][Chain Tx Not Found]!
// @StatusErr[ConfigConflict][409999006][Config Conflict]!
// @StatusErr[ConfigInitFailed][500999006][Config Init Failed]!
// @StatusErr[ConfigNotFound][404999003][Config Not Found]!
// @StatusErr[ConfigParseFailed][500999008][Config Parse Failed]!
// @StatusErr[ConfigUninitFailed][500999007][Config Uninit Failed]!
// @StatusErr[Conflict][409999001][Conflict conflict error]!
// @StatusErr[ContractLogConflict][409999011][Contract Log Conflict]!
// @StatusErr[ContractLogNotFound][404999014][Contract Log Not Found]!
// @StatusErr[CreateChannelFailed][500999004][Create Message Channel Failed]!
// @StatusErr[CreateInstanceFailed][500999010][Create Instance Failed]!
// @StatusErr[CronJobConflict][409999010][Cron Job Conflict]!
// @StatusErr[CronJobNotFound][404999011][Cron Job Not Found]!
// @StatusErr[CurrentAccountAbsence][401999012][Current Account Absence]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DeprecatedProject][400999003][Deprecated Project]!
// @StatusErr[DisabledAccount][403999002][Disabled Account]!
// @StatusErr[FetchResourceFailed][500999005][Fetch Resource Failed]!
// @StatusErr[Forbidden][403999001][forbidden]
// @StatusErr[GenPublisherTokenFailed][500999009][Gen Publisher Token Failed]!
// @StatusErr[InstanceNotFound][404999006][Instance Not Found]!
// @StatusErr[InstanceNotRunning][404999012][Instance Not Running]!
// @StatusErr[InternalServerError][500999001][internal error]
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthPublisherID][401999011][Invalid Auth Publisher ID]!
// @StatusErr[InvalidAuthPublisherID][401999011][Invalid Auth Publisher ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidConfigType][400999002][Invalid Config Type]!
// @StatusErr[InvalidCronExpressions][400999005][Invalid Cron Expressions]!
// @StatusErr[InvalidEthLoginMessage][401999010][Invalid Siwe Message]!
// @StatusErr[InvalidEthLoginSignature][401999009][Invalid Siwe Signature]!
// @StatusErr[InvalidEventChannel][401999013][Invalid Event Channel]!
// @StatusErr[InvalidEventToken][401999014][Invalid Event Token]!
// @StatusErr[InvalidNewPassword][401999007][Invalid New Password]!
// @StatusErr[InvalidOldPassword][401999006][Invalid Old Password]!
// @StatusErr[InvalidPassword][401999008][Invalid Password]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[MD5ChecksumFailed][500999012][Md5 Checksum Failed]!
// @StatusErr[MqttConnectFailed][500999014][MQTT Connect Failed]!
// @StatusErr[MqttSubscribeFailed][500999013][MQTT Subscribe Failed]!
// @StatusErr[MultiInstanceDeployed][409999008][Multi Instance Deployed]!
// @StatusErr[NoAdminPermission][401999005][No Admin Permission]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[NotFound][404999001][NotFound]!
// @StatusErr[ProjectNameConflict][409999002][Project Name Conflict]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!
// @StatusErr[PublisherConflict][409999007][Publisher Conflict]!
// @StatusErr[PublisherNotFound][404999008][Publisher Not Found]!
// @StatusErr[PublisherNotFound][404999008][Publisher Not Found]!
// @StatusErr[ResourceConflict][409999003][Resource Conflict]!
// @StatusErr[ResourceNotFound][404999004][Resource Not Found]!
// @StatusErr[ResourceOwnerConflict][409999004][Resource Owner Conflict]!
// @StatusErr[ResourcePermNotFound][404999010][Resource Perm Not Found]!
// @StatusErr[StrategyConflict][409999005][Strategy Conflict]!
// @StatusErr[StrategyNotFound][404999007][Strategy Not Found]!
// @StatusErr[TopicAlreadySubscribed][403999007][Topic Already Subscribed]!
// @StatusErr[Unauthorized][401999001][unauthorized]
// @StatusErr[UnknownDeployCommand][400999004][Unknown Deploy Command]!
// @StatusErr[UploadFileDiskLimit][403999006][Upload File Disk Limit]!
// @StatusErr[UploadFileFailed][500999003][Upload File Failed]!
// @StatusErr[UploadFileMd5Unmatched][403999005][Upload File Md5 Unmatched]!
// @StatusErr[UploadFileSizeLimit][403999004][Upload File Size Limit]!
// @StatusErr[WhiteListForbidden][403999003][White List Forbidden]!

func (o *HandleEvent) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.HandleEvent")
	return cli.Do(ctx, o, metas...)
}

func (o *HandleEvent) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesEventEventRsp, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModulesEventEventRsp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *HandleEvent) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesEventEventRsp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type ListApplet struct {
	ProjectName  string                                               `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string                                               `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AppletIDs    []GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"query" name:"appletID,omitempty"`
	AuthInQuery  string                                               `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	LNameLike    string                                               `in:"query" name:"lName,omitempty"`
	NameLike     string                                               `in:"query" name:"name,omitempty"`
	Names        []string                                             `in:"query" name:"names,omitempty"`
	Offset       int64                                                `in:"query" default:"0" name:"offset,omitempty" validate:"@int64[0,]"`
	Size         int64                                                `in:"query" default:"10" name:"size,omitempty" validate:"@int64[-1,]"`
}

func (o *ListApplet) Path() string {
	return "/srv-applet-mgr/v0/applet/x/:projectName/datalist"
}

func (o *ListApplet) Method() string {
	return "GET"
}

// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[CurrentAccountAbsence][401999012][Current Account Absence]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!

func (o *ListApplet) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.ListApplet")
	return cli.Do(ctx, o, metas...)
}

func (o *ListApplet) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAppletListRsp, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModulesAppletListRsp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *ListApplet) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAppletListRsp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type ListCronJob struct {
	ProjectID    GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID   `in:"path" name:"projectID"`
	AuthInHeader string                                               `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                               `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	CronJobIDs   []GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"query" name:"cronJobID,omitempty"`
	EventTypes   []string                                             `in:"query" name:"eventType,omitempty"`
	Offset       int64                                                `in:"query" default:"0" name:"offset,omitempty" validate:"@int64[0,]"`
	Size         int64                                                `in:"query" default:"10" name:"size,omitempty" validate:"@int64[-1,]"`
}

func (o *ListCronJob) Path() string {
	return "/srv-applet-mgr/v0/cronjob"
}

func (o *ListCronJob) Method() string {
	return "GET"
}

// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!

func (o *ListCronJob) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.ListCronJob")
	return cli.Do(ctx, o, metas...)
}

func (o *ListCronJob) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesCronjobListRsp, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModulesCronjobListRsp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *ListCronJob) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesCronjobListRsp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type ListProject struct {
	AuthInHeader string                                               `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                               `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	Names        []string                                             `in:"query" name:"name,omitempty"`
	Offset       int64                                                `in:"query" default:"0" name:"offset,omitempty" validate:"@int64[0,]"`
	ProjectIDs   []GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"query" name:"projectID,omitempty"`
	Size         int64                                                `in:"query" default:"10" name:"size,omitempty" validate:"@int64[-1,]"`
	Versions     []string                                             `in:"query" name:"version,omitempty"`
}

func (o *ListProject) Path() string {
	return "/srv-applet-mgr/v0/project/datalist"
}

func (o *ListProject) Method() string {
	return "GET"
}

// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!

func (o *ListProject) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.ListProject")
	return cli.Do(ctx, o, metas...)
}

func (o *ListProject) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesProjectListRsp, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModulesProjectListRsp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *ListProject) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesProjectListRsp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type ListProjectDetail struct {
	AuthInHeader string                                               `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                               `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	Names        []string                                             `in:"query" name:"name,omitempty"`
	Offset       int64                                                `in:"query" default:"0" name:"offset,omitempty" validate:"@int64[0,]"`
	ProjectIDs   []GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"query" name:"projectID,omitempty"`
	Size         int64                                                `in:"query" default:"10" name:"size,omitempty" validate:"@int64[-1,]"`
	Versions     []string                                             `in:"query" name:"version,omitempty"`
}

func (o *ListProjectDetail) Path() string {
	return "/srv-applet-mgr/v0/project/detail_list"
}

func (o *ListProjectDetail) Method() string {
	return "GET"
}

// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InstanceNotFound][404999006][Instance Not Found]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[ResourceNotFound][404999004][Resource Not Found]!

func (o *ListProjectDetail) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.ListProjectDetail")
	return cli.Do(ctx, o, metas...)
}

func (o *ListProjectDetail) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesProjectListDetailRsp, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModulesProjectListDetailRsp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *ListProjectDetail) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesProjectListDetailRsp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type ListPublisher struct {
	ProjectName  string                                               `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string                                               `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                               `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	Keys         []string                                             `in:"query" name:"key,omitempty"`
	LNameLike    string                                               `in:"query" name:"lname,omitempty"`
	NameLike     string                                               `in:"query" name:"name,omitempty"`
	Offset       int64                                                `in:"query" default:"0" name:"offset,omitempty" validate:"@int64[0,]"`
	PublisherIDs []GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"query" name:"publisherIDs,omitempty"`
	RNameLike    string                                               `in:"query" name:"rname,omitempty"`
	Size         int64                                                `in:"query" default:"10" name:"size,omitempty" validate:"@int64[-1,]"`
}

func (o *ListPublisher) Path() string {
	return "/srv-applet-mgr/v0/publisher/x/:projectName"
}

func (o *ListPublisher) Method() string {
	return "GET"
}

// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[CurrentAccountAbsence][401999012][Current Account Absence]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!

func (o *ListPublisher) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.ListPublisher")
	return cli.Do(ctx, o, metas...)
}

func (o *ListPublisher) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesPublisherListRsp, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModulesPublisherListRsp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *ListPublisher) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesPublisherListRsp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type ListResources struct {
	AuthInHeader   string                                                  `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery    string                                                  `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	ExpireAfter    GithubComMachinefiW3BstreamPkgDependsBaseTypesTimestamp `in:"query" name:"expireAfter,omitempty"`
	ExpireBefore   GithubComMachinefiW3BstreamPkgDependsBaseTypesTimestamp `in:"query" name:"expireBefore,omitempty"`
	Filenames      []string                                                `in:"query" name:"filename,omitempty"`
	FilenameLike   string                                                  `in:"query" name:"filenameLike,omitempty"`
	Md5            string                                                  `in:"query" name:"md5,omitempty"`
	Offset         int64                                                   `in:"query" default:"0" name:"offset,omitempty" validate:"@int64[0,]"`
	ResourceIDs    []GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID    `in:"query" name:"resourceID,omitempty"`
	Size           int64                                                   `in:"query" default:"10" name:"size,omitempty" validate:"@int64[-1,]"`
	UploadedAfter  GithubComMachinefiW3BstreamPkgDependsBaseTypesTimestamp `in:"query" name:"uploadedAfter,omitempty"`
	UploadedBefore GithubComMachinefiW3BstreamPkgDependsBaseTypesTimestamp `in:"query" name:"uploadedBefore,omitempty"`
}

func (o *ListResources) Path() string {
	return "/srv-applet-mgr/v0/resource"
}

func (o *ListResources) Method() string {
	return "GET"
}

// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!

func (o *ListResources) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.ListResources")
	return cli.Do(ctx, o, metas...)
}

func (o *ListResources) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesResourceListRsp, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModulesResourceListRsp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *ListResources) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesResourceListRsp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type ListStrategy struct {
	ProjectName  string                                               `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string                                               `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AppletIDs    []GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"query" name:"appletID,omitempty"`
	AuthInQuery  string                                               `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	EventTypes   []string                                             `in:"query" name:"eventType,omitempty"`
	Handlers     []string                                             `in:"query" name:"handler,omitempty"`
	Offset       int64                                                `in:"query" default:"0" name:"offset,omitempty" validate:"@int64[0,]"`
	Size         int64                                                `in:"query" default:"10" name:"size,omitempty" validate:"@int64[-1,]"`
	StrategyIDs  []GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"query" name:"strategyID,omitempty"`
}

func (o *ListStrategy) Path() string {
	return "/srv-applet-mgr/v0/strategy/x/:projectName/datalist"
}

func (o *ListStrategy) Method() string {
	return "GET"
}

// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[CurrentAccountAbsence][401999012][Current Account Absence]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!

func (o *ListStrategy) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.ListStrategy")
	return cli.Do(ctx, o, metas...)
}

func (o *ListStrategy) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesStrategyListRsp, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModulesStrategyListRsp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *ListStrategy) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesStrategyListRsp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type Liveness struct {
}

func (o *Liveness) Path() string {
	return "/liveness"
}

func (o *Liveness) Method() string {
	return "GET"
}

func (o *Liveness) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.Liveness")
	return cli.Do(ctx, o, metas...)
}

func (o *Liveness) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*map[string]string, kit.Metadata, error) {
	rsp := new(map[string]string)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *Liveness) Invoke(cli kit.Client, metas ...kit.Metadata) (*map[string]string, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type LoginByEthAddress struct {
	LoginByEthAddressReq GithubComMachinefiW3BstreamPkgModulesAccountLoginByEthAddressReq `in:"body"`
}

func (o *LoginByEthAddress) Path() string {
	return "/srv-applet-mgr/v0/login/wallet"
}

func (o *LoginByEthAddress) Method() string {
	return "PUT"
}

// @StatusErr[AccountIdentityConflict][409999014][Account Identity Conflict]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DisabledAccount][403999002][Disabled Account]!
// @StatusErr[InternalServerError][500999001][internal error]
// @StatusErr[InvalidEthLoginMessage][401999010][Invalid Siwe Message]!
// @StatusErr[InvalidEthLoginSignature][401999009][Invalid Siwe Signature]!
// @StatusErr[WhiteListForbidden][403999003][White List Forbidden]!

func (o *LoginByEthAddress) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.LoginByEthAddress")
	return cli.Do(ctx, o, metas...)
}

func (o *LoginByEthAddress) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccountLoginRsp, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModulesAccountLoginRsp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *LoginByEthAddress) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccountLoginRsp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type LoginByUsername struct {
	LoginByUsernameReq GithubComMachinefiW3BstreamPkgModulesAccountLoginByUsernameReq `in:"body"`
}

func (o *LoginByUsername) Path() string {
	return "/srv-applet-mgr/v0/login"
}

func (o *LoginByUsername) Method() string {
	return "PUT"
}

// @StatusErr[AccountIdentityNotFound][404999009][Account Identity Not Found]!
// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[AccountPasswordNotFound][404999018][Account Password Not Found]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DisabledAccount][403999002][Disabled Account]!
// @StatusErr[InternalServerError][500999001][internal error]
// @StatusErr[InvalidPassword][401999008][Invalid Password]!

func (o *LoginByUsername) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.LoginByUsername")
	return cli.Do(ctx, o, metas...)
}

func (o *LoginByUsername) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccountLoginRsp, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModulesAccountLoginRsp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *LoginByUsername) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccountLoginRsp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type RemoveApplet struct {
	AppletID     GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"appletID"`
	AuthInHeader string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *RemoveApplet) Path() string {
	return "/srv-applet-mgr/v0/applet/data/:appletID"
}

func (o *RemoveApplet) Method() string {
	return "DELETE"
}

func (o *RemoveApplet) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.RemoveApplet")
	return cli.Do(ctx, o, metas...)
}

func (o *RemoveApplet) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *RemoveApplet) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type RemoveChainHeight struct {
	ChainHeightID GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"chainHeightID"`
	ProjectName   string                                             `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader  string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery   string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *RemoveChainHeight) Path() string {
	return "/srv-applet-mgr/v0/monitor/x/:projectName/chain_height/:chainHeightID"
}

func (o *RemoveChainHeight) Method() string {
	return "DELETE"
}

func (o *RemoveChainHeight) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.RemoveChainHeight")
	return cli.Do(ctx, o, metas...)
}

func (o *RemoveChainHeight) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *RemoveChainHeight) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type RemoveChainTx struct {
	ChainTxID    GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"chainTxID"`
	ProjectName  string                                             `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *RemoveChainTx) Path() string {
	return "/srv-applet-mgr/v0/monitor/x/:projectName/chain_tx/:chainTxID"
}

func (o *RemoveChainTx) Method() string {
	return "DELETE"
}

func (o *RemoveChainTx) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.RemoveChainTx")
	return cli.Do(ctx, o, metas...)
}

func (o *RemoveChainTx) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *RemoveChainTx) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type RemoveContractLog struct {
	ContractLogID GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"contractLogID"`
	ProjectName   string                                             `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader  string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery   string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *RemoveContractLog) Path() string {
	return "/srv-applet-mgr/v0/monitor/x/:projectName/contract_log/:contractLogID"
}

func (o *RemoveContractLog) Method() string {
	return "DELETE"
}

func (o *RemoveContractLog) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.RemoveContractLog")
	return cli.Do(ctx, o, metas...)
}

func (o *RemoveContractLog) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *RemoveContractLog) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type RemoveCronJob struct {
	CronJobID    GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"cronJobID"`
	AuthInHeader string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *RemoveCronJob) Path() string {
	return "/srv-applet-mgr/v0/cronjob/data/:cronJobID"
}

func (o *RemoveCronJob) Method() string {
	return "DELETE"
}

func (o *RemoveCronJob) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.RemoveCronJob")
	return cli.Do(ctx, o, metas...)
}

func (o *RemoveCronJob) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *RemoveCronJob) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type RemoveInstance struct {
	InstanceID   GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"instanceID"`
	AuthInHeader string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *RemoveInstance) Path() string {
	return "/srv-applet-mgr/v0/deploy/data/:instanceID"
}

func (o *RemoveInstance) Method() string {
	return "DELETE"
}

func (o *RemoveInstance) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.RemoveInstance")
	return cli.Do(ctx, o, metas...)
}

func (o *RemoveInstance) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *RemoveInstance) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type RemoveProject struct {
	ProjectName  string `in:"path" name:"projectName" validate:"@projectName"`
	AuthInHeader string `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *RemoveProject) Path() string {
	return "/srv-applet-mgr/v0/project/x/:projectName"
}

func (o *RemoveProject) Method() string {
	return "DELETE"
}

func (o *RemoveProject) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.RemoveProject")
	return cli.Do(ctx, o, metas...)
}

func (o *RemoveProject) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *RemoveProject) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type RemovePublisher struct {
	PublisherID  GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"publisherID"`
	AuthInHeader string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *RemovePublisher) Path() string {
	return "/srv-applet-mgr/v0/publisher/data/:publisherID"
}

func (o *RemovePublisher) Method() string {
	return "DELETE"
}

func (o *RemovePublisher) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.RemovePublisher")
	return cli.Do(ctx, o, metas...)
}

func (o *RemovePublisher) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *RemovePublisher) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type RemoveResource struct {
	ResourceID   GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"resourceID"`
	AuthInHeader string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *RemoveResource) Path() string {
	return "/srv-applet-mgr/v0/resource/:resourceID"
}

func (o *RemoveResource) Method() string {
	return "DELETE"
}

func (o *RemoveResource) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.RemoveResource")
	return cli.Do(ctx, o, metas...)
}

func (o *RemoveResource) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *RemoveResource) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type RemoveStrategy struct {
	StrategyID   GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID `in:"path" name:"strategyID"`
	AuthInHeader string                                             `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                             `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
}

func (o *RemoveStrategy) Path() string {
	return "/srv-applet-mgr/v0/strategy/data/:strategyID"
}

func (o *RemoveStrategy) Method() string {
	return "DELETE"
}

func (o *RemoveStrategy) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.RemoveStrategy")
	return cli.Do(ctx, o, metas...)
}

func (o *RemoveStrategy) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *RemoveStrategy) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type UpdateApplet struct {
	AppletID     GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID   `in:"path" name:"appletID"`
	AuthInHeader string                                               `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                               `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	UpdateReq    GithubComMachinefiW3BstreamPkgModulesAppletUpdateReq `in:"body" mime:"multipart"`
}

func (o *UpdateApplet) Path() string {
	return "/srv-applet-mgr/v0/applet/:appletID"
}

func (o *UpdateApplet) Method() string {
	return "PUT"
}

// @StatusErr[AccountNotFound][404999017][Account Not Found]!
// @StatusErr[AppletNameConflict][409999009][Applet Name Conflict]!
// @StatusErr[AppletNotFound][404999005][Applet Not Found]!
// @StatusErr[ConfigConflict][409999006][Config Conflict]!
// @StatusErr[ConfigInitFailed][500999006][Config Init Failed]!
// @StatusErr[ConfigParseFailed][500999008][Config Parse Failed]!
// @StatusErr[ConfigUninitFailed][500999007][Config Uninit Failed]!
// @StatusErr[CreateInstanceFailed][500999010][Create Instance Failed]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[DatabaseError][500999002][Database Error]!
// @StatusErr[InstanceNotFound][404999006][Instance Not Found]!
// @StatusErr[InternalServerError][500999001][internal error]
// @StatusErr[InvalidAuthAccountID][401999003][Invalid Auth Account ID]!
// @StatusErr[InvalidAuthValue][401999002][Invalid Auth Value]!
// @StatusErr[InvalidClaim][401999003][Invalid Claim]!
// @StatusErr[InvalidConfigType][400999002][Invalid Config Type]!
// @StatusErr[InvalidToken][401999002][Invalid Token]!
// @StatusErr[MD5ChecksumFailed][500999012][Md5 Checksum Failed]!
// @StatusErr[MultiInstanceDeployed][409999008][Multi Instance Deployed]!
// @StatusErr[NoProjectPermission][401999004][No Project Permission]!
// @StatusErr[ProjectNotFound][404999002][Project Not Found]!
// @StatusErr[ResourceConflict][409999003][Resource Conflict]!
// @StatusErr[ResourceNotFound][404999004][Resource Not Found]!
// @StatusErr[StrategyConflict][409999005][Strategy Conflict]!
// @StatusErr[UploadFileDiskLimit][403999006][Upload File Disk Limit]!
// @StatusErr[UploadFileFailed][500999003][Upload File Failed]!
// @StatusErr[UploadFileMd5Unmatched][403999005][Upload File Md5 Unmatched]!
// @StatusErr[UploadFileSizeLimit][403999004][Upload File Size Limit]!

func (o *UpdateApplet) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.UpdateApplet")
	return cli.Do(ctx, o, metas...)
}

func (o *UpdateApplet) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAppletCreateRsp, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgModulesAppletCreateRsp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *UpdateApplet) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAppletCreateRsp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type UpdatePasswordByAccountID struct {
	AuthInHeader      string                                                        `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery       string                                                        `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	UpdatePasswordReq GithubComMachinefiW3BstreamPkgModulesAccountUpdatePasswordReq `in:"body"`
}

func (o *UpdatePasswordByAccountID) Path() string {
	return "/srv-applet-mgr/v0/account"
}

func (o *UpdatePasswordByAccountID) Method() string {
	return "PUT"
}

func (o *UpdatePasswordByAccountID) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.UpdatePasswordByAccountID")
	return cli.Do(ctx, o, metas...)
}

func (o *UpdatePasswordByAccountID) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *UpdatePasswordByAccountID) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type UpdatePublisher struct {
	ProjectName  string                                                  `in:"path" name:"projectName" validate:"@projectName"`
	PublisherID  GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID      `in:"path" name:"publisherID"`
	AuthInHeader string                                                  `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                                  `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	UpdateReq    GithubComMachinefiW3BstreamPkgModulesPublisherUpdateReq `in:"body"`
}

func (o *UpdatePublisher) Path() string {
	return "/srv-applet-mgr/v0/publisher/x/:projectName/:publisherID"
}

func (o *UpdatePublisher) Method() string {
	return "PUT"
}

func (o *UpdatePublisher) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.UpdatePublisher")
	return cli.Do(ctx, o, metas...)
}

func (o *UpdatePublisher) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *UpdatePublisher) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type UpdateStrategy struct {
	StrategyID   GithubComMachinefiW3BstreamPkgDependsBaseTypesSFID     `in:"path" name:"strategyID"`
	AuthInHeader string                                                 `in:"header" name:"Authorization,omitempty" validate:"@string[1,]"`
	AuthInQuery  string                                                 `in:"query" name:"authorization,omitempty" validate:"@string[1,]"`
	UpdateReq    GithubComMachinefiW3BstreamPkgModulesStrategyCreateReq `in:"body"`
}

func (o *UpdateStrategy) Path() string {
	return "/srv-applet-mgr/v0/strategy/:strategyID"
}

func (o *UpdateStrategy) Method() string {
	return "PUT"
}

func (o *UpdateStrategy) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.UpdateStrategy")
	return cli.Do(ctx, o, metas...)
}

func (o *UpdateStrategy) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *UpdateStrategy) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type VersionRouter struct {
}

func (o *VersionRouter) Path() string {
	return "/version"
}

func (o *VersionRouter) Method() string {
	return "GET"
}

func (o *VersionRouter) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "applet-mgr.VersionRouter")
	return cli.Do(ctx, o, metas...)
}

func (o *VersionRouter) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*string, kit.Metadata, error) {
	rsp := new(string)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *VersionRouter) Invoke(cli kit.Client, metas ...kit.Metadata) (*string, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}
