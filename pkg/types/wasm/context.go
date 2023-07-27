package wasm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
)

type (
	CtxSQLStore          struct{}
	CtxKVStore           struct{}
	CtxLogger            struct{}
	CtxEnv               struct{}
	CtxRedisPrefix       struct{}
	CtxChainClient       struct{}
	CtxRuntimeResource   struct{}
	CtxRuntimeEventTypes struct{}
	CtxMqttClient        struct{}
	CtxCustomMetrics     struct{}
	CtxFlow              struct{}
	CtxTraceInfo         struct{}
)

func WithSQLStore(ctx context.Context, v *Database) context.Context {
	return contextx.WithValue(ctx, CtxSQLStore{}, v)
}

func WithSQLStoreContext(v *Database) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxSQLStore{}, v)
	}
}

func SQLStoreFromContext(ctx context.Context) (*Database, bool) {
	v, ok := ctx.Value(CtxSQLStore{}).(*Database)
	return v, ok
}

func MustSQLStoreFromContext(ctx context.Context) *Database {
	v, ok := SQLStoreFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithKVStore(ctx context.Context, v KVStore) context.Context {
	return contextx.WithValue(ctx, CtxKVStore{}, v)
}

func WithKVStoreContext(v KVStore) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxKVStore{}, v)
	}
}

func KVStoreFromContext(ctx context.Context) (KVStore, bool) {
	v, ok := ctx.Value(CtxKVStore{}).(KVStore)
	return v, ok
}

func MustKVStoreFromContext(ctx context.Context) KVStore {
	v, ok := KVStoreFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithLogger(ctx context.Context, v *Logger) context.Context {
	return contextx.WithValue(ctx, CtxLogger{}, v)
}

func WithLoggerContext(v *Logger) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxLogger{}, v)
	}
}

func LoggerFromContext(ctx context.Context) (*Logger, bool) {
	v, ok := ctx.Value(CtxLogger{}).(*Logger)
	return v, ok
}

func MustLoggerFromContext(ctx context.Context) *Logger {
	v, ok := LoggerFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithEnv(ctx context.Context, v *Env) context.Context {
	return contextx.WithValue(ctx, CtxEnv{}, v)
}

func WithEnvContext(v *Env) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxEnv{}, v)
	}
}

func EnvFromContext(ctx context.Context) (*Env, bool) {
	v, ok := ctx.Value(CtxEnv{}).(*Env)
	return v, ok
}

func MustEnvFromContext(ctx context.Context) *Env {
	v, ok := EnvFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithRedisPrefix(ctx context.Context, v string) context.Context {
	return contextx.WithValue(ctx, CtxRedisPrefix{}, v)
}

func WithRedisPrefixContext(v string) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxRedisPrefix{}, v)
	}
}

func RedisPrefixFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(CtxRedisPrefix{}).(string)
	return v, ok
}

func MustRedisPrefixFromContext(ctx context.Context) string {
	v, ok := RedisPrefixFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithChainClient(ctx context.Context, v *ChainClient) context.Context {
	return contextx.WithValue(ctx, CtxChainClient{}, v)
}

func WithChainClientContext(v *ChainClient) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxChainClient{}, v)
	}
}

func ChainClientFromContext(ctx context.Context) (*ChainClient, bool) {
	v, ok := ctx.Value(CtxChainClient{}).(*ChainClient)
	return v, ok
}

func MustChainClientFromContext(ctx context.Context) *ChainClient {
	v, ok := ChainClientFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithRuntimeResource(ctx context.Context, v *mapx.Map[uint32, []byte]) context.Context {
	return contextx.WithValue(ctx, CtxRuntimeResource{}, v)
}

func WithRuntimeResourceContext(v *mapx.Map[uint32, []byte]) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxRuntimeResource{}, v)
	}
}

func RuntimeResourceFromContext(ctx context.Context) (*mapx.Map[uint32, []byte], bool) {
	v, ok := ctx.Value(CtxRuntimeResource{}).(*mapx.Map[uint32, []byte])
	return v, ok
}

func MustRuntimeResourceFromContext(ctx context.Context) *mapx.Map[uint32, []byte] {
	v, ok := RuntimeResourceFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithRuntimeEventTypes(ctx context.Context, v *mapx.Map[uint32, []byte]) context.Context {
	return contextx.WithValue(ctx, CtxRuntimeEventTypes{}, v)
}

func WithRuntimeEventTypesContext(v *mapx.Map[uint32, []byte]) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxRuntimeEventTypes{}, v)
	}
}

func RuntimeEventTypesFromContext(ctx context.Context) (*mapx.Map[uint32, []byte], bool) {
	v, ok := ctx.Value(CtxRuntimeEventTypes{}).(*mapx.Map[uint32, []byte])
	return v, ok
}

func MustRuntimeEventTypesFromContext(ctx context.Context) *mapx.Map[uint32, []byte] {
	v, ok := RuntimeResourceFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithMQTTClient(ctx context.Context, mq *mqtt.Client) context.Context {
	return contextx.WithValue(ctx, CtxMqttClient{}, mq)
}

func WithMQTTClientContext(mq *mqtt.Client) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxMqttClient{}, mq)
	}
}

func MQTTClientFromContext(ctx context.Context) (*mqtt.Client, bool) {
	v, ok := ctx.Value(CtxMqttClient{}).(*mqtt.Client)
	return v, ok
}

func MustMQTTClientFromContext(ctx context.Context) *mqtt.Client {
	v, ok := MQTTClientFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithCustomMetrics(ctx context.Context, mt metrics.CustomMetrics) context.Context {
	return contextx.WithValue(ctx, CtxCustomMetrics{}, mt)
}

func WithCustomMetricsContext(mt metrics.CustomMetrics) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxCustomMetrics{}, mt)
	}
}

func CustomMetricsFromContext(ctx context.Context) (metrics.CustomMetrics, bool) {
	v, ok := ctx.Value(CtxCustomMetrics{}).(metrics.CustomMetrics)
	return v, ok
}

func MustCustomMetricsFromContext(ctx context.Context) metrics.CustomMetrics {
	v, ok := CustomMetricsFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithFlow(ctx context.Context, flow *Flow) context.Context {
	return contextx.WithValue(ctx, CtxFlow{}, flow)
}

func WithFlowContext(flow *Flow) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxFlow{}, flow)
	}
}

func FlowFromContext(ctx context.Context) (*Flow, bool) {
	v, ok := ctx.Value(CtxFlow{}).(*Flow)
	return v, ok
}

func MustFlowFromContext(ctx context.Context) *Flow {
	v, ok := FlowFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithTraceInfo(ctx context.Context, v *TraceInfo) context.Context {
	return contextx.WithValue(ctx, CtxTraceInfo{}, v)
}

func WithTraceInfoContext(flow *TraceInfo) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxTraceInfo{}, flow)
	}
}

func TraceInfoFromContext(ctx context.Context) (*TraceInfo, bool) {
	v, ok := ctx.Value(CtxTraceInfo{}).(*TraceInfo)
	return v, ok
}

func MustTraceInfoFromContext(ctx context.Context) *TraceInfo {
	v, ok := TraceInfoFromContext(ctx)
	must.BeTrue(ok)
	return v
}
