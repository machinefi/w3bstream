// This is a generated source file. DO NOT EDIT
// Source: applet_mgr/client.go

package applet_mgr

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

type Interface interface {
	Context() context.Context
	WithContext(context.Context) Interface
	BatchRemoveApplet(req *BatchRemoveApplet, metas ...kit.Metadata) (kit.Metadata, error)
	BatchRemoveInstance(req *BatchRemoveInstance, metas ...kit.Metadata) (kit.Metadata, error)
	BatchRemovePublisher(req *BatchRemovePublisher, metas ...kit.Metadata) (kit.Metadata, error)
	BatchRemoveStrategy(req *BatchRemoveStrategy, metas ...kit.Metadata) (kit.Metadata, error)
	BatchRemoveTrafficLimit(req *BatchRemoveTrafficLimit, metas ...kit.Metadata) (kit.Metadata, error)
	ChainConfig(metas ...kit.Metadata) (*ChainConfigResp, kit.Metadata, error)
	ControlChainHeight(req *ControlChainHeight, metas ...kit.Metadata) (kit.Metadata, error)
	ControlChainTx(req *ControlChainTx, metas ...kit.Metadata) (kit.Metadata, error)
	ControlContractLog(req *ControlContractLog, metas ...kit.Metadata) (kit.Metadata, error)
	ControlInstance(req *ControlInstance, metas ...kit.Metadata) (kit.Metadata, error)
	CreateAccountAccessKey(req *CreateAccountAccessKey, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccessKeyCreateRsp, kit.Metadata, error)
	CreateAccountByUsernameAndPassword(req *CreateAccountByUsernameAndPassword, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccountCreateAccountByUsernameRsp, kit.Metadata, error)
	CreateAndStartInstance(req *CreateAndStartInstance, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsInstance, kit.Metadata, error)
	CreateApplet(req *CreateApplet, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAppletCreateRsp, kit.Metadata, error)
	CreateChainHeight(req *CreateChainHeight, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsChainHeight, kit.Metadata, error)
	CreateChainTx(req *CreateChainTx, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsChainTx, kit.Metadata, error)
	CreateContractLog(req *CreateContractLog, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsContractLog, kit.Metadata, error)
	CreateCronJob(req *CreateCronJob, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsCronJob, kit.Metadata, error)
	CreateOperator(req *CreateOperator, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsOperator, kit.Metadata, error)
	CreateOrUpdateProjectEnv(req *CreateOrUpdateProjectEnv, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsConfig, kit.Metadata, error)
	CreateOrUpdateProjectFlow(req *CreateOrUpdateProjectFlow, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsConfig, kit.Metadata, error)
	CreateProject(req *CreateProject, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesProjectCreateRsp, kit.Metadata, error)
	CreateProjectOperator(req *CreateProjectOperator, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsProjectOperator, kit.Metadata, error)
	CreateProjectSchema(req *CreateProjectSchema, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsConfig, kit.Metadata, error)
	CreatePublisher(req *CreatePublisher, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsPublisher, kit.Metadata, error)
	CreateStrategy(req *CreateStrategy, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsStrategy, kit.Metadata, error)
	CreateTrafficLimit(req *CreateTrafficLimit, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsTrafficLimit, kit.Metadata, error)
	DeleteAccountAccessKeyByName(req *DeleteAccountAccessKeyByName, metas ...kit.Metadata) (kit.Metadata, error)
	DownloadResource(req *DownloadResource, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxAttachment, kit.Metadata, error)
	EthClient(metas ...kit.Metadata) (*EthClientRsp, kit.Metadata, error)
	GetAccessKeyByName(req *GetAccessKeyByName, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccessKeyListData, kit.Metadata, error)
	GetApplet(req *GetApplet, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsApplet, kit.Metadata, error)
	GetDownloadResourceUrl(req *GetDownloadResourceUrl, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesResourceDownLoadResourceRsp, kit.Metadata, error)
	GetInstanceByAppletID(req *GetInstanceByAppletID, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsInstance, kit.Metadata, error)
	GetInstanceByInstanceID(req *GetInstanceByInstanceID, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsInstance, kit.Metadata, error)
	GetOperatorAddr(req *GetOperatorAddr, metas ...kit.Metadata) (*string, kit.Metadata, error)
	GetProject(req *GetProject, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsProject, kit.Metadata, error)
	GetProjectEnv(req *GetProjectEnv, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgTypesWasmEnv, kit.Metadata, error)
	GetProjectFlow(req *GetProjectFlow, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgTypesWasmFlow, kit.Metadata, error)
	GetProjectOperator(req *GetProjectOperator, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesOperatorDetail, kit.Metadata, error)
	GetProjectSchema(req *GetProjectSchema, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgTypesWasmDatabase, kit.Metadata, error)
	GetPublisher(req *GetPublisher, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsPublisher, kit.Metadata, error)
	GetStrategy(req *GetStrategy, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsStrategy, kit.Metadata, error)
	GetTrafficLimit(req *GetTrafficLimit, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsTrafficLimit, kit.Metadata, error)
	HandleEvent(req *HandleEvent, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesEventEventRsp, kit.Metadata, error)
	ListAccessGroupMetas(req *ListAccessGroupMetas, metas ...kit.Metadata) (*[]GithubComMachinefiW3BstreamPkgModulesAccessKeyGroupMetaBase, kit.Metadata, error)
	ListAccountAccessKey(req *ListAccountAccessKey, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccessKeyListRsp, kit.Metadata, error)
	ListApplet(req *ListApplet, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAppletListRsp, kit.Metadata, error)
	ListCronJob(req *ListCronJob, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesCronjobListRsp, kit.Metadata, error)
	ListOperator(req *ListOperator, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesOperatorListDetailRsp, kit.Metadata, error)
	ListProject(req *ListProject, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesProjectListRsp, kit.Metadata, error)
	ListProjectDetail(req *ListProjectDetail, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesProjectListDetailRsp, kit.Metadata, error)
	ListPublisher(req *ListPublisher, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesPublisherListRsp, kit.Metadata, error)
	ListResources(req *ListResources, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesResourceListRsp, kit.Metadata, error)
	ListStrategy(req *ListStrategy, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesStrategyListRsp, kit.Metadata, error)
	ListTrafficLimit(req *ListTrafficLimit, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesTrafficlimitListRsp, kit.Metadata, error)
	Liveness(metas ...kit.Metadata) (*map[string]string, kit.Metadata, error)
	LoginByEthAddress(req *LoginByEthAddress, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccountLoginRsp, kit.Metadata, error)
	LoginByUsername(req *LoginByUsername, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccountLoginRsp, kit.Metadata, error)
	RemoveApplet(req *RemoveApplet, metas ...kit.Metadata) (kit.Metadata, error)
	RemoveChainHeight(req *RemoveChainHeight, metas ...kit.Metadata) (kit.Metadata, error)
	RemoveChainTx(req *RemoveChainTx, metas ...kit.Metadata) (kit.Metadata, error)
	RemoveContractLog(req *RemoveContractLog, metas ...kit.Metadata) (kit.Metadata, error)
	RemoveCronJob(req *RemoveCronJob, metas ...kit.Metadata) (kit.Metadata, error)
	RemoveInstance(req *RemoveInstance, metas ...kit.Metadata) (kit.Metadata, error)
	RemoveOperator(req *RemoveOperator, metas ...kit.Metadata) (kit.Metadata, error)
	RemoveProject(req *RemoveProject, metas ...kit.Metadata) (kit.Metadata, error)
	RemoveProjectOperator(req *RemoveProjectOperator, metas ...kit.Metadata) (kit.Metadata, error)
	RemovePublisher(req *RemovePublisher, metas ...kit.Metadata) (kit.Metadata, error)
	RemoveResource(req *RemoveResource, metas ...kit.Metadata) (kit.Metadata, error)
	RemoveStrategy(req *RemoveStrategy, metas ...kit.Metadata) (kit.Metadata, error)
	RemoveTrafficLimit(req *RemoveTrafficLimit, metas ...kit.Metadata) (kit.Metadata, error)
	RemoveWasmLogByInstanceID(req *RemoveWasmLogByInstanceID, metas ...kit.Metadata) (kit.Metadata, error)
	UpdateAccountAccessKeyByName(req *UpdateAccountAccessKeyByName, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccessKeyUpdateRsp, kit.Metadata, error)
	UpdateApplet(req *UpdateApplet, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAppletCreateRsp, kit.Metadata, error)
	UpdatePasswordByAccountID(req *UpdatePasswordByAccountID, metas ...kit.Metadata) (kit.Metadata, error)
	UpdatePublisher(req *UpdatePublisher, metas ...kit.Metadata) (kit.Metadata, error)
	UpdateStrategy(req *UpdateStrategy, metas ...kit.Metadata) (kit.Metadata, error)
	UpdateTrafficLimit(req *UpdateTrafficLimit, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsTrafficLimit, kit.Metadata, error)
	VersionRouter(metas ...kit.Metadata) (*string, kit.Metadata, error)
}

func NewClient(c kit.Client) *Client {
	return &(Client{
		Client: c,
	})
}

type Client struct {
	Client kit.Client
	ctx    context.Context
}

func (c *Client) Context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

func (c *Client) WithContext(ctx context.Context) Interface {
	cc := new(Client)
	cc.Client, cc.ctx = c.Client, ctx
	return cc
}

func (c *Client) BatchRemoveApplet(req *BatchRemoveApplet, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) BatchRemoveInstance(req *BatchRemoveInstance, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) BatchRemovePublisher(req *BatchRemovePublisher, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) BatchRemoveStrategy(req *BatchRemoveStrategy, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) BatchRemoveTrafficLimit(req *BatchRemoveTrafficLimit, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ChainConfig(metas ...kit.Metadata) (*ChainConfigResp, kit.Metadata, error) {
	return (&ChainConfig{}).InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ControlChainHeight(req *ControlChainHeight, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ControlChainTx(req *ControlChainTx, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ControlContractLog(req *ControlContractLog, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ControlInstance(req *ControlInstance, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateAccountAccessKey(req *CreateAccountAccessKey, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccessKeyCreateRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateAccountByUsernameAndPassword(req *CreateAccountByUsernameAndPassword, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccountCreateAccountByUsernameRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateAndStartInstance(req *CreateAndStartInstance, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsInstance, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateApplet(req *CreateApplet, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAppletCreateRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateChainHeight(req *CreateChainHeight, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsChainHeight, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateChainTx(req *CreateChainTx, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsChainTx, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateContractLog(req *CreateContractLog, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsContractLog, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateCronJob(req *CreateCronJob, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsCronJob, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateOperator(req *CreateOperator, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsOperator, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateOrUpdateProjectEnv(req *CreateOrUpdateProjectEnv, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsConfig, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateOrUpdateProjectFlow(req *CreateOrUpdateProjectFlow, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsConfig, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateProject(req *CreateProject, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesProjectCreateRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateProjectOperator(req *CreateProjectOperator, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsProjectOperator, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateProjectSchema(req *CreateProjectSchema, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsConfig, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreatePublisher(req *CreatePublisher, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsPublisher, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateStrategy(req *CreateStrategy, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsStrategy, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) CreateTrafficLimit(req *CreateTrafficLimit, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsTrafficLimit, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) DeleteAccountAccessKeyByName(req *DeleteAccountAccessKeyByName, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) DownloadResource(req *DownloadResource, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxAttachment, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) EthClient(metas ...kit.Metadata) (*EthClientRsp, kit.Metadata, error) {
	return (&EthClient{}).InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetAccessKeyByName(req *GetAccessKeyByName, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccessKeyListData, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetApplet(req *GetApplet, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsApplet, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetDownloadResourceUrl(req *GetDownloadResourceUrl, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesResourceDownLoadResourceRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetInstanceByAppletID(req *GetInstanceByAppletID, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsInstance, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetInstanceByInstanceID(req *GetInstanceByInstanceID, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsInstance, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetOperatorAddr(req *GetOperatorAddr, metas ...kit.Metadata) (*string, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetProject(req *GetProject, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsProject, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetProjectEnv(req *GetProjectEnv, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgTypesWasmEnv, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetProjectFlow(req *GetProjectFlow, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgTypesWasmFlow, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetProjectOperator(req *GetProjectOperator, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesOperatorDetail, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetProjectSchema(req *GetProjectSchema, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgTypesWasmDatabase, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetPublisher(req *GetPublisher, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsPublisher, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetStrategy(req *GetStrategy, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsStrategy, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetTrafficLimit(req *GetTrafficLimit, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsTrafficLimit, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) HandleEvent(req *HandleEvent, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesEventEventRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ListAccessGroupMetas(req *ListAccessGroupMetas, metas ...kit.Metadata) (*[]GithubComMachinefiW3BstreamPkgModulesAccessKeyGroupMetaBase, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ListAccountAccessKey(req *ListAccountAccessKey, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccessKeyListRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ListApplet(req *ListApplet, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAppletListRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ListCronJob(req *ListCronJob, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesCronjobListRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ListOperator(req *ListOperator, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesOperatorListDetailRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ListProject(req *ListProject, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesProjectListRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ListProjectDetail(req *ListProjectDetail, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesProjectListDetailRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ListPublisher(req *ListPublisher, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesPublisherListRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ListResources(req *ListResources, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesResourceListRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ListStrategy(req *ListStrategy, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesStrategyListRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ListTrafficLimit(req *ListTrafficLimit, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesTrafficlimitListRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) Liveness(metas ...kit.Metadata) (*map[string]string, kit.Metadata, error) {
	return (&Liveness{}).InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) LoginByEthAddress(req *LoginByEthAddress, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccountLoginRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) LoginByUsername(req *LoginByUsername, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccountLoginRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemoveApplet(req *RemoveApplet, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemoveChainHeight(req *RemoveChainHeight, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemoveChainTx(req *RemoveChainTx, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemoveContractLog(req *RemoveContractLog, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemoveCronJob(req *RemoveCronJob, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemoveInstance(req *RemoveInstance, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemoveOperator(req *RemoveOperator, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemoveProject(req *RemoveProject, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemoveProjectOperator(req *RemoveProjectOperator, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemovePublisher(req *RemovePublisher, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemoveResource(req *RemoveResource, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemoveStrategy(req *RemoveStrategy, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemoveTrafficLimit(req *RemoveTrafficLimit, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemoveWasmLogByInstanceID(req *RemoveWasmLogByInstanceID, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) UpdateAccountAccessKeyByName(req *UpdateAccountAccessKeyByName, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAccessKeyUpdateRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) UpdateApplet(req *UpdateApplet, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModulesAppletCreateRsp, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) UpdatePasswordByAccountID(req *UpdatePasswordByAccountID, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) UpdatePublisher(req *UpdatePublisher, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) UpdateStrategy(req *UpdateStrategy, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) UpdateTrafficLimit(req *UpdateTrafficLimit, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgModelsTrafficLimit, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) VersionRouter(metas ...kit.Metadata) (*string, kit.Metadata, error) {
	return (&VersionRouter{}).InvokeContext(c.Context(), c.Client, metas...)
}
