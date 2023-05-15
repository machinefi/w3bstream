package types

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/filesystem"
	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/conf/redis"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/client"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/models"
)

type (
	CtxMgrDBExecutor     struct{} // CtxMgrDBExecutor sqlx.DBExecutor
	CtxMonitorDBExecutor struct{} // CtxMonitorDBExecutor sqlx.DBExecutor
	CtxWasmDBEndpoint    struct{} // CtxWasmDBEndpoint sqlx.DBExecutor
	CtxLogger            struct{} // CtxLogger log.Logger
	CtxMqttBroker        struct{} // CtxMqttBroker mqtt.Broker
	CtxRedisEndpoint     struct{} // CtxRedisEndpoint redis.Redis
	CtxUploadConfig      struct{} // CtxUploadConfig UploadConfig
	CtxTaskWorker        struct{}
	CtxTaskBoard         struct{}
	CtxProject           struct{}
	CtxApplet            struct{}
	CtxResource          struct{}
	CtxInstance          struct{}
	CtxEthClient         struct{} // CtxEthClient ETHClientConfig
	CtxWhiteList         struct{}
	CtxFileSystem        struct{}
	CtxStrategy          struct{}
	CtxPublisher         struct{}
	CtxCronJob           struct{}
	CtxOperator          struct{}
	ContractLog          struct{}
	ChainHeight          struct{}
	ChainTx              struct{}
	CtxAccount           struct{}
	CtxStrategyResults   struct{}
	CtxFileSystemOp      struct{}
	CtxProxyClient       struct{}
	CtxResourceOwnership struct{}
)

func WithStrategyResults(ctx context.Context, v []*StrategyResult) context.Context {
	return contextx.WithValue(ctx, CtxStrategyResults{}, v)
}

func WithStrategyResultsContext(v []*StrategyResult) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxStrategyResults{}, v)
	}
}

func StrategyResultsFromContext(ctx context.Context) ([]*StrategyResult, bool) {
	v, ok := ctx.Value(CtxStrategyResults{}).([]*StrategyResult)
	return v, ok
}

func MustStrategyResultsFromContext(ctx context.Context) []*StrategyResult {
	v, ok := StrategyResultsFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithAccount(ctx context.Context, v *models.Account) context.Context {
	return contextx.WithValue(ctx, CtxAccount{}, v)
}

func WithAccountContext(v *models.Account) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxAccount{}, v)
	}
}

func AccountFromContext(ctx context.Context) (*models.Account, bool) {
	v, ok := ctx.Value(CtxAccount{}).(*models.Account)
	return v, ok
}

func MustAccountFromContext(ctx context.Context) *models.Account {
	v, ok := AccountFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithPublisher(ctx context.Context, v *models.Publisher) context.Context {
	return contextx.WithValue(ctx, CtxPublisher{}, v)
}

func WithPublisherContext(v *models.Publisher) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxPublisher{}, v)
	}
}

func PublisherFromContext(ctx context.Context) (*models.Publisher, bool) {
	v, ok := ctx.Value(CtxPublisher{}).(*models.Publisher)
	return v, ok
}

func MustPublisherFromContext(ctx context.Context) *models.Publisher {
	v, ok := PublisherFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithCronJob(ctx context.Context, v *models.CronJob) context.Context {
	return contextx.WithValue(ctx, CtxCronJob{}, v)
}

func WithCronJobContext(v *models.CronJob) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxCronJob{}, v)
	}
}

func CronJobFromContext(ctx context.Context) (*models.CronJob, bool) {
	v, ok := ctx.Value(CtxCronJob{}).(*models.CronJob)
	return v, ok
}

func MustCronJobFromContext(ctx context.Context) *models.CronJob {
	v, ok := CronJobFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithOperator(ctx context.Context, v *models.Operator) context.Context {
	return contextx.WithValue(ctx, CtxOperator{}, v)
}

func WithOperatorContext(v *models.Operator) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxOperator{}, v)
	}
}

func OperatorFromContext(ctx context.Context) (*models.Operator, bool) {
	v, ok := ctx.Value(CtxOperator{}).(*models.Operator)
	return v, ok
}

func MustOperatorFromContext(ctx context.Context) *models.Operator {
	v, ok := OperatorFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithContractLog(ctx context.Context, v *models.ContractLog) context.Context {
	return contextx.WithValue(ctx, ContractLog{}, v)
}

func WithContractLogContext(v *models.ContractLog) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, ContractLog{}, v)
	}
}

func ContractLogFromContext(ctx context.Context) (*models.ContractLog, bool) {
	v, ok := ctx.Value(ContractLog{}).(*models.ContractLog)
	return v, ok
}

func MustContractLogFromContext(ctx context.Context) *models.ContractLog {
	v, ok := ContractLogFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithChainHeight(ctx context.Context, v *models.ChainHeight) context.Context {
	return contextx.WithValue(ctx, ChainHeight{}, v)
}

func WithChainHeightContext(v *models.ChainHeight) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, ChainHeight{}, v)
	}
}

func ChainHeightFromContext(ctx context.Context) (*models.ChainHeight, bool) {
	v, ok := ctx.Value(ChainHeight{}).(*models.ChainHeight)
	return v, ok
}

func MustChainHeightFromContext(ctx context.Context) *models.ChainHeight {
	v, ok := ChainHeightFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithChainTx(ctx context.Context, v *models.ChainTx) context.Context {
	return contextx.WithValue(ctx, ChainTx{}, v)
}

func WithChainTxContext(v *models.ChainTx) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, ChainTx{}, v)
	}
}

func ChainTxFromContext(ctx context.Context) (*models.ChainTx, bool) {
	v, ok := ctx.Value(ChainTx{}).(*models.ChainTx)
	return v, ok
}

func MustChainTxFromContext(ctx context.Context) *models.ChainTx {
	v, ok := ChainTxFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithStrategy(ctx context.Context, v *models.Strategy) context.Context {
	return contextx.WithValue(ctx, CtxStrategy{}, v)
}

func WithStrategyContext(v *models.Strategy) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxStrategy{}, v)
	}
}

func StrategyFromContext(ctx context.Context) (*models.Strategy, bool) {
	v, ok := ctx.Value(CtxStrategy{}).(*models.Strategy)
	return v, ok
}

func MustStrategyFromContext(ctx context.Context) *models.Strategy {
	v, ok := StrategyFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithMgrDBExecutor(ctx context.Context, v sqlx.DBExecutor) context.Context {
	return contextx.WithValue(ctx, CtxMgrDBExecutor{}, v)
}

func WithMgrDBExecutorContext(v sqlx.DBExecutor) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxMgrDBExecutor{}, v)
	}
}

func MgrDBExecutorFromContext(ctx context.Context) (sqlx.DBExecutor, bool) {
	v, ok := ctx.Value(CtxMgrDBExecutor{}).(sqlx.DBExecutor)
	return v, ok
}

func MustMgrDBExecutorFromContext(ctx context.Context) sqlx.DBExecutor {
	v, ok := MgrDBExecutorFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithMonitorDBExecutor(ctx context.Context, v sqlx.DBExecutor) context.Context {
	return contextx.WithValue(ctx, CtxMonitorDBExecutor{}, v)
}

func WithMonitorDBExecutorContext(v sqlx.DBExecutor) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxMonitorDBExecutor{}, v)
	}
}

func MonitorDBExecutorFromContext(ctx context.Context) (sqlx.DBExecutor, bool) {
	v, ok := ctx.Value(CtxMonitorDBExecutor{}).(sqlx.DBExecutor)
	return v, ok
}

func MustMonitorDBExecutorFromContext(ctx context.Context) sqlx.DBExecutor {
	v, ok := MonitorDBExecutorFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithWasmDBEndpoint(ctx context.Context, v *types.Endpoint) context.Context {
	return contextx.WithValue(ctx, CtxWasmDBEndpoint{}, v)
}

func WithWasmDBEndpointContext(v *types.Endpoint) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxWasmDBEndpoint{}, v)
	}
}

func WasmDBEndpointFromContext(ctx context.Context) (*types.Endpoint, bool) {
	v, ok := ctx.Value(CtxWasmDBEndpoint{}).(*types.Endpoint)
	return v, ok
}

func MustWasmDBEndpointFromContext(ctx context.Context) *types.Endpoint {
	v, ok := WasmDBEndpointFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithRedis(ctx context.Context, v *redis.Redis) context.Context {
	return contextx.WithValue(ctx, CtxRedisEndpoint{}, v)
}

func WithRedisEndpointContext(v *redis.Redis) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxRedisEndpoint{}, v)
	}
}

func RedisEndpointFromContext(ctx context.Context) (*redis.Redis, bool) {
	v, ok := ctx.Value(CtxRedisEndpoint{}).(*redis.Redis)
	return v, ok
}

func MustRedisEndpointFromContext(ctx context.Context) *redis.Redis {
	v, ok := RedisEndpointFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithLogger(ctx context.Context, v log.Logger) context.Context {
	return contextx.WithValue(ctx, CtxLogger{}, v)
}

func WithLoggerContext(v log.Logger) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxLogger{}, v)
	}
}

func LoggerFromContext(ctx context.Context) (log.Logger, bool) {
	v, ok := ctx.Value(CtxLogger{}).(log.Logger)
	return v, ok
}

func MustLoggerFromContext(ctx context.Context) log.Logger {
	v, ok := LoggerFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithMqttBroker(ctx context.Context, v *mqtt.Broker) context.Context {
	return contextx.WithValue(ctx, CtxMqttBroker{}, v)
}

func WithMqttBrokerContext(v *mqtt.Broker) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxMqttBroker{}, v)
	}
}

func MqttBrokerFromContext(ctx context.Context) (*mqtt.Broker, bool) {
	v, ok := ctx.Value(CtxMqttBroker{}).(*mqtt.Broker)
	return v, ok
}

func MustMqttBrokerFromContext(ctx context.Context) *mqtt.Broker {
	v, ok := MqttBrokerFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithTaskBoard(ctx context.Context, tb *mq.TaskBoard) context.Context {
	return contextx.WithValue(ctx, CtxTaskBoard{}, tb)
}

func WithTaskBoardContext(tb *mq.TaskBoard) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return WithTaskBoard(ctx, tb)
	}
}

func TaskBoardFromContext(ctx context.Context) (*mq.TaskBoard, bool) {
	v, ok := ctx.Value(CtxTaskBoard{}).(*mq.TaskBoard)
	return v, ok
}

func MustTaskBoardFromContext(ctx context.Context) *mq.TaskBoard {
	v, ok := TaskBoardFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithTaskWorker(ctx context.Context, tw *mq.TaskWorker) context.Context {
	return contextx.WithValue(ctx, CtxTaskWorker{}, tw)
}

func WithTaskWorkerContext(tw *mq.TaskWorker) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return WithTaskWorker(ctx, tw)
	}
}

func TaskWorkerFromContext(ctx context.Context) (*mq.TaskWorker, bool) {
	v, ok := ctx.Value(CtxTaskWorker{}).(*mq.TaskWorker)
	return v, ok
}

func MustTaskWorkerFromContext(ctx context.Context) *mq.TaskWorker {
	v, ok := TaskWorkerFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithProject(ctx context.Context, p *models.Project) context.Context {
	_p := *p
	return contextx.WithValue(ctx, CtxProject{}, &_p)
}

func WithProjectContext(p *models.Project) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return WithProject(ctx, p)
	}
}

func ProjectFromContext(ctx context.Context) (*models.Project, bool) {
	v, ok := ctx.Value(CtxProject{}).(*models.Project)
	return v, ok
}

func MustProjectFromContext(ctx context.Context) *models.Project {
	v, ok := ProjectFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithApplet(ctx context.Context, a *models.Applet) context.Context {
	_a := *a
	return contextx.WithValue(ctx, CtxApplet{}, &_a)
}

func WithAppletContext(a *models.Applet) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return WithApplet(ctx, a)
	}
}

func AppletFromContext(ctx context.Context) (*models.Applet, bool) {
	v, ok := ctx.Value(CtxApplet{}).(*models.Applet)
	return v, ok
}

func MustAppletFromContext(ctx context.Context) *models.Applet {
	v, ok := AppletFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithResource(ctx context.Context, r *models.Resource) context.Context {
	_r := *r
	return contextx.WithValue(ctx, CtxResource{}, &_r)
}

func WithResourceContext(r *models.Resource) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return WithResource(ctx, r)
	}
}

func ResourceFromContext(ctx context.Context) (*models.Resource, bool) {
	v, ok := ctx.Value(CtxResource{}).(*models.Resource)
	return v, ok
}

func MustResourceFromContext(ctx context.Context) *models.Resource {
	v, ok := ResourceFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithInstance(ctx context.Context, i *models.Instance) context.Context {
	_i := *i
	return contextx.WithValue(ctx, CtxInstance{}, &_i)
}

func WithInstanceContext(i *models.Instance) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return WithInstance(ctx, i)
	}
}

func InstanceFromContext(ctx context.Context) (*models.Instance, bool) {
	v, ok := ctx.Value(CtxInstance{}).(*models.Instance)
	return v, ok
}

func MustInstanceFromContext(ctx context.Context) *models.Instance {
	v, ok := InstanceFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithETHClientConfig(ctx context.Context, v *ETHClientConfig) context.Context {
	return contextx.WithValue(ctx, CtxEthClient{}, v)
}

func WithETHClientConfigContext(v *ETHClientConfig) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxEthClient{}, v)
	}
}

func ETHClientConfigFromContext(ctx context.Context) (*ETHClientConfig, bool) {
	v, ok := ctx.Value(CtxEthClient{}).(*ETHClientConfig)
	return v, ok
}

func MustETHClientConfigFromContext(ctx context.Context) *ETHClientConfig {
	v, ok := ETHClientConfigFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithWhiteList(ctx context.Context, v *WhiteList) context.Context {
	return contextx.WithValue(ctx, CtxWhiteList{}, v)
}

func WithWhiteListContext(v *WhiteList) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxWhiteList{}, v)
	}
}

func WhiteListFromContext(ctx context.Context) (*WhiteList, bool) {
	v, ok := ctx.Value(CtxWhiteList{}).(*WhiteList)
	return v, ok
}

func MustWhiteListFromContext(ctx context.Context) *WhiteList {
	v, ok := WhiteListFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithFileSystemOp(ctx context.Context, v filesystem.FileSystemOp) context.Context {
	return contextx.WithValue(ctx, CtxFileSystemOp{}, v)
}

func WithFileSystemOpContext(v filesystem.FileSystemOp) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxFileSystemOp{}, v)
	}
}

func FileSystemOpFromContext(ctx context.Context) (filesystem.FileSystemOp, bool) {
	v, ok := ctx.Value(CtxFileSystemOp{}).(filesystem.FileSystemOp)
	return v, ok
}

func MustFileSystemOpFromContext(ctx context.Context) filesystem.FileSystemOp {
	v, ok := FileSystemOpFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithProxyClient(ctx context.Context, v *client.Client) context.Context {
	return contextx.WithValue(ctx, CtxProxyClient{}, v)
}

func WithProxyClientContext(v *client.Client) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxProxyClient{}, v)
	}
}

func ProxyClientFromContext(ctx context.Context) (*client.Client, bool) {
	v, ok := ctx.Value(CtxProxyClient{}).(*client.Client)
	return v, ok
}

func MustProxyClientFromContext(ctx context.Context) *client.Client {
	v, ok := ProxyClientFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithResourceOwnership(ctx context.Context, o *models.ResourceOwnership) context.Context {
	return contextx.WithValue(ctx, CtxResourceOwnership{}, o)
}

func WithResourceOwnershipContext(o *models.ResourceOwnership) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxResourceOwnership{}, o)
	}
}

func ResourceOwnershipFromContext(ctx context.Context) (*models.ResourceOwnership, bool) {
	v, ok := ctx.Value(CtxResourceOwnership{}).(*models.ResourceOwnership)
	return v, ok
}

func MustResourceOwnershipFromContext(ctx context.Context) *models.ResourceOwnership {
	v, ok := ResourceOwnershipFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithUploadConfig(ctx context.Context, v *UploadConfig) context.Context {
	return contextx.WithValue(ctx, CtxUploadConfig{}, v)
}

func WithUploadConfigContext(v *UploadConfig) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxUploadConfig{}, v)
	}
}

func UploadConfigFromContext(ctx context.Context) (*UploadConfig, bool) {
	v, ok := ctx.Value(CtxUploadConfig{}).(*UploadConfig)
	return v, ok
}

func MustUploadConfigFromContext(ctx context.Context) *UploadConfig {
	v, ok := UploadConfigFromContext(ctx)
	must.BeTrue(ok)
	return v
}
