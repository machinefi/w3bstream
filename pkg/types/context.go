package types

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/conf/postgres"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/models"
)

type Context uint8

type (
	CtxDBExecutor        struct{} // CtxDBExecutor sqlx.DBExecutor
	CtxMonitorDBExecutor struct{} // CtxMonitorDBExecutor sqlx.DBExecutor
	CtxPgEndpoint        struct{} // CtxPgEndpoint postgres.Endpoint
	CtxLogger            struct{} // CtxLogger log.Logger
	CtxMqttBroker        struct{} // CtxMqttBroker mqtt.Broker
	CtxUploadConfig      struct{} // CtxUploadConfig UploadConfig
	CtxEthClient         struct{} // CtxEthClient ETHClientConfig
	CtxTaskWorker        struct{}
	CtxTaskBoard         struct{}
	CtxProject           struct{}
	CtxApplet            struct{}
)

func WithDBExecutor(ctx context.Context, v sqlx.DBExecutor) context.Context {
	return contextx.WithValue(ctx, CtxDBExecutor{}, v)
}

func WithDBExecutorContext(v sqlx.DBExecutor) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxDBExecutor{}, v)
	}
}

func DBExecutorFromContext(ctx context.Context) (sqlx.DBExecutor, bool) {
	v, ok := ctx.Value(CtxDBExecutor{}).(sqlx.DBExecutor)
	return v, ok
}

func MustDBExecutorFromContext(ctx context.Context) sqlx.DBExecutor {
	v, ok := DBExecutorFromContext(ctx)
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
