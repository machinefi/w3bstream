package types

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/conf/filesystem"
	"github.com/machinefi/w3bstream/pkg/depends/conf/filesystem/amazonS3"
	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/conf/postgres"
	"github.com/machinefi/w3bstream/pkg/depends/conf/redis"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/models"
)

type (
	CtxMgrDBExecutor     struct{} // CtxMgrDBExecutor sqlx.DBExecutor
	CtxMonitorDBExecutor struct{} // CtxMonitorDBExecutor sqlx.DBExecutor
	CtxWasmDBExecutor    struct{} // CtxWasmDBExecutor sqlx.DBExecutor
	CtxPgEndpoint        struct{} // CtxPgEndpoint postgres.Endpoint
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
	CtxAmazonS3          struct{}
	CtxFileSystemOp      struct{}
)

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

func WithWasmDBExecutor(ctx context.Context, v sqlx.DBExecutor) context.Context {
	return contextx.WithValue(ctx, CtxWasmDBExecutor{}, v)
}

func WithWasmDBExecutorContext(v sqlx.DBExecutor) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxWasmDBExecutor{}, v)
	}
}

func WasmDBExecutorFromContext(ctx context.Context) (sqlx.DBExecutor, bool) {
	v, ok := ctx.Value(CtxWasmDBExecutor{}).(sqlx.DBExecutor)
	return v, ok
}

func MustWasmDBExecutorFromContext(ctx context.Context) sqlx.DBExecutor {
	v, ok := WasmDBExecutorFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithPgEndpoint(ctx context.Context, v postgres.Endpoint) context.Context {
	return contextx.WithValue(ctx, CtxPgEndpoint{}, v)
}

func WithPgEndpointContext(v *postgres.Endpoint) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxPgEndpoint{}, v)
	}
}

func PgEndpointFromContext(ctx context.Context) (*postgres.Endpoint, bool) {
	v, ok := ctx.Value(CtxPgEndpoint{}).(*postgres.Endpoint)
	return v, ok
}

func MustPgEndpointFromContext(ctx context.Context) *postgres.Endpoint {
	v, ok := PgEndpointFromContext(ctx)
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

func WithAmazonS3(ctx context.Context, v *amazonS3.AmazonS3) context.Context {
	return contextx.WithValue(ctx, CtxAmazonS3{}, v)
}

func WithAmazonS3Context(v *amazonS3.AmazonS3) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxAmazonS3{}, v)
	}
}

func AmazonS3FromContext(ctx context.Context) (*amazonS3.AmazonS3, bool) {
	v, ok := ctx.Value(CtxAmazonS3{}).(*amazonS3.AmazonS3)
	return v, ok
}

func MustS3FromContext(ctx context.Context) *amazonS3.AmazonS3 {
	v, ok := AmazonS3FromContext(ctx)
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
