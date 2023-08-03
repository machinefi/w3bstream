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
	wasmapi "github.com/machinefi/w3bstream/pkg/modules/vm/wasmapi/types"
)

// global contexts
type (
	// CtxMgrDBExecutor type sqlx.DBExecutor for global manager server database
	CtxMgrDBExecutor struct{}
	// CtxMonitorDBExecutor type sqlx.DBExecutor for global monitor server database
	CtxMonitorDBExecutor struct{}
	// CtxWasmDBEndpoint type *types.Endpoint. for global wasm database endpoint
	CtxWasmDBEndpoint struct{}
	// CtxLogger type log.Logger. service logger
	CtxLogger struct{}
	// CtxMqttBroker *mqtt.Broker. mqtt broker
	CtxMqttBroker struct{}
	// CtxRedisEndpoint type *redis.Redis. redis endpoint
	CtxRedisEndpoint struct{}
	// CtxUploadConfig type *UploadConfig. resource upload configuration
	CtxUploadConfig struct{}
	// CtxTaskWorker type *mq.TaskWorker. service async task worker
	CtxTaskWorker struct{}
	// CtxTaskBoard type *mq.TaskBoard service async task manager
	CtxTaskBoard struct{}
	// CtxWhiteList type *EthAddressWhiteList global eth address white list
	CtxEthAddressWhiteList struct{}
	// CtxEthClient type *ETHClientConfig global eth chain endpoints
	CtxEthClient struct{}
	// CtxChainConfig type *ChainConfig global chain endpoints
	CtxChainConfig struct{}
	// CtxFileSystemOp type filesystem.FileSystemOp describe resource storing operation type
	CtxFileSystemOp struct{}
	// CtxProxyClient type *client.Client http client for forwarding mqtt event
	CtxProxyClient struct{}
	// CtxWasmDBConfig type *WasmDBConfig wasm database config TODO combine with WasmDBEndpoint
	CtxWasmDBConfig struct{}
	// CtxRobotNotifierConfig type *RobotNotifierConfig for notify service level message to maintainers.
	CtxRobotNotifierConfig struct{}
	// CtxMetricsCenterConfig *MetricsCenterConfig for metrics
	CtxMetricsCenterConfig struct{}
)

// model contexts
type (
	// CtxProject type *models.Project
	CtxProject struct{}
	// CtxApplet type *models.Applet
	CtxApplet struct{}
	// CtxResource type *models.Resource
	CtxResource struct{}
	// CtxInstance type *models.Instance
	CtxInstance struct{}
	// CtxStrategy type *models.Strategy
	CtxStrategy struct{}
	// CtxPublisher type *models.Publisher
	CtxPublisher struct{}
	// CtxCronJob type *models.CronJob
	CtxCronJob struct{}
	// CtxOperator type *models.Operator
	CtxOperator struct{}
	// CtxOperators type []models.Operator filtered operators
	CtxOperators struct{}
	// CtxProjectOperator type *models.ProjectOperator
	CtxProjectOperator struct{}
	// CtxContractLog type *models.ContractLog
	CtxContractLog struct{}
	// CtxChainHeight type *models.ChainHeight
	CtxChainHeight struct{}
	// CtxChainTx type *models.ChainTx
	CtxChainTx struct{}
	// CtxAccount type *models.Account
	CtxAccount struct{}
	// CtxResourceOwnership type *models.ResourceOwnership
	CtxResourceOwnership struct{}
	// CtxTrafficLimit type *models.TrafficLimit
	CtxTrafficLimit struct{}
)
type (
	// CtxStrategyResults type []*StrategyResult event strategies
	CtxStrategyResults struct{} // CtxStrategyResults
	// CtxEventID type string. current event id
	CtxEventID struct{}
	// CtxWasmApiServer type wasmapi/types.Server wasm global async server TODO move to wasm context package
	CtxWasmApiServer struct{}
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
	return contextx.WithValue(ctx, CtxContractLog{}, v)
}

func WithContractLogContext(v *models.ContractLog) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxContractLog{}, v)
	}
}

func ContractLogFromContext(ctx context.Context) (*models.ContractLog, bool) {
	v, ok := ctx.Value(CtxContractLog{}).(*models.ContractLog)
	return v, ok
}

func MustContractLogFromContext(ctx context.Context) *models.ContractLog {
	v, ok := ContractLogFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithChainHeight(ctx context.Context, v *models.ChainHeight) context.Context {
	return contextx.WithValue(ctx, CtxChainHeight{}, v)
}

func WithChainHeightContext(v *models.ChainHeight) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxChainHeight{}, v)
	}
}

func ChainHeightFromContext(ctx context.Context) (*models.ChainHeight, bool) {
	v, ok := ctx.Value(CtxChainHeight{}).(*models.ChainHeight)
	return v, ok
}

func MustChainHeightFromContext(ctx context.Context) *models.ChainHeight {
	v, ok := ChainHeightFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithChainTx(ctx context.Context, v *models.ChainTx) context.Context {
	return contextx.WithValue(ctx, CtxChainTx{}, v)
}

func WithChainTxContext(v *models.ChainTx) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxChainTx{}, v)
	}
}

func ChainTxFromContext(ctx context.Context) (*models.ChainTx, bool) {
	v, ok := ctx.Value(CtxChainTx{}).(*models.ChainTx)
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

func WithChainConfig(ctx context.Context, v *ChainConfig) context.Context {
	return contextx.WithValue(ctx, CtxChainConfig{}, v)
}

func WithChainConfigContext(v *ChainConfig) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxChainConfig{}, v)
	}
}

func ChainConfigFromContext(ctx context.Context) (*ChainConfig, bool) {
	v, ok := ctx.Value(CtxChainConfig{}).(*ChainConfig)
	return v, ok
}

func MustChainConfigFromContext(ctx context.Context) *ChainConfig {
	v, ok := ChainConfigFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithEthAddressWhiteList(ctx context.Context, v *EthAddressWhiteList) context.Context {
	return contextx.WithValue(ctx, CtxEthAddressWhiteList{}, v)
}

func WithEthAddressWhiteListContext(v *EthAddressWhiteList) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxEthAddressWhiteList{}, v)
	}
}

func EthAddressWhiteListFromContext(ctx context.Context) (*EthAddressWhiteList, bool) {
	v, ok := ctx.Value(CtxEthAddressWhiteList{}).(*EthAddressWhiteList)
	return v, ok
}

func MustEthAddressWhiteListFromContext(ctx context.Context) *EthAddressWhiteList {
	v, ok := EthAddressWhiteListFromContext(ctx)
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

func WithWasmDBConfig(ctx context.Context, v *WasmDBConfig) context.Context {
	return contextx.WithValue(ctx, CtxWasmDBConfig{}, v)
}

func WithWasmDBConfigContext(v *WasmDBConfig) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxWasmDBConfig{}, v)
	}
}

func WasmDBConfigFromContext(ctx context.Context) (*WasmDBConfig, bool) {
	v, ok := ctx.Value(CtxWasmDBConfig{}).(*WasmDBConfig)
	return v, ok
}

func MustWasmDBConfigFromContext(ctx context.Context) *WasmDBConfig {
	v, ok := WasmDBConfigFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithEventID(ctx context.Context, v string) context.Context {
	return contextx.WithValue(ctx, CtxEventID{}, v)
}

func WithEventIDContext(v string) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxEventID{}, v)
	}
}

func EventIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(CtxEventID{}).(string)
	return v, ok
}

func MustEventIDFromContext(ctx context.Context) string {
	v, ok := EventIDFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithTrafficLimit(ctx context.Context, r *models.TrafficLimit) context.Context {
	_r := *r
	return contextx.WithValue(ctx, CtxTrafficLimit{}, &_r)
}

func WithTrafficLimitContext(r *models.TrafficLimit) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return WithTrafficLimit(ctx, r)
	}
}

func TrafficLimitFromContext(ctx context.Context) (*models.TrafficLimit, bool) {
	v, ok := ctx.Value(CtxTrafficLimit{}).(*models.TrafficLimit)
	return v, ok
}

func MustTrafficLimitFromContext(ctx context.Context) *models.TrafficLimit {
	v, ok := TrafficLimitFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithMetricsCenterConfig(ctx context.Context, v *MetricsCenterConfig) context.Context {
	return contextx.WithValue(ctx, CtxMetricsCenterConfig{}, v)
}

func WithMetricsCenterConfigContext(v *MetricsCenterConfig) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxMetricsCenterConfig{}, v)
	}
}

func MetricsCenterConfigFromContext(ctx context.Context) (*MetricsCenterConfig, bool) {
	v, ok := ctx.Value(CtxMetricsCenterConfig{}).(*MetricsCenterConfig)
	return v, ok
}

func MustMetricsCenterConfigFromContext(ctx context.Context) *MetricsCenterConfig {
	v, ok := MetricsCenterConfigFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithRobotNotifierConfig(ctx context.Context, v *RobotNotifierConfig) context.Context {
	return contextx.WithValue(ctx, CtxRobotNotifierConfig{}, v)
}

func WithRobotNotifierConfigContext(v *RobotNotifierConfig) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxRobotNotifierConfig{}, v)
	}
}

func RobotNotifierConfigFromContext(ctx context.Context) (*RobotNotifierConfig, bool) {
	v, ok := ctx.Value(CtxRobotNotifierConfig{}).(*RobotNotifierConfig)
	return v, ok
}

func MustRobotNotifierConfigFromContext(ctx context.Context) *RobotNotifierConfig {
	v, ok := RobotNotifierConfigFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithWasmApiServer(ctx context.Context, v wasmapi.Server) context.Context {
	return contextx.WithValue(ctx, CtxWasmApiServer{}, v)
}

func WithWasmApiServerContext(v wasmapi.Server) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxWasmApiServer{}, v)
	}
}

func WasmApiServerFromContext(ctx context.Context) (wasmapi.Server, bool) {
	v, ok := ctx.Value(CtxWasmApiServer{}).(wasmapi.Server)
	return v, ok
}

func MustWasmApiServerFromContext(ctx context.Context) wasmapi.Server {
	v, ok := WasmApiServerFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithOperators(ctx context.Context, v []models.Operator) context.Context {
	return contextx.WithValue(ctx, CtxOperators{}, v)
}

func WithOperatorsContext(v []models.Operator) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxOperators{}, v)
	}
}

func OperatorsFromContext(ctx context.Context) ([]models.Operator, bool) {
	v, ok := ctx.Value(CtxOperators{}).([]models.Operator)
	return v, ok
}

func MustOperatorsFromContext(ctx context.Context) []models.Operator {
	v, ok := OperatorsFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithProjectOperator(ctx context.Context, v *models.ProjectOperator) context.Context {
	return contextx.WithValue(ctx, CtxProjectOperator{}, v)
}

func WithProjectOperatorContext(v *models.ProjectOperator) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxProjectOperator{}, v)
	}
}

func ProjectOperatorFromContext(ctx context.Context) (*models.ProjectOperator, bool) {
	v, ok := ctx.Value(CtxProjectOperator{}).(*models.ProjectOperator)
	return v, ok
}

func MustProjectOperatorFromContext(ctx context.Context) *models.ProjectOperator {
	v, ok := ProjectOperatorFromContext(ctx)
	must.BeTrue(ok)
	return v
}
