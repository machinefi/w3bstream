package wasm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
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
)

func WithSQLStore(ctx context.Context, v SQLStore) context.Context {
	return contextx.WithValue(ctx, CtxSQLStore{}, v)
}

func WithSQLStoreContext(v SQLStore) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxSQLStore{}, v)
	}
}

func SQLStoreFromContext(ctx context.Context) (SQLStore, bool) {
	v, ok := ctx.Value(CtxSQLStore{}).(SQLStore)
	return v, ok
}

func MustSQLStoreFromContext(ctx context.Context) SQLStore {
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

func WithMQTTClient(ctx context.Context, mq *MqttClient) context.Context {
	return contextx.WithValue(ctx, CtxMqttClient{}, mq)
}

func WithMQTTClientContext(mq *MqttClient) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxMqttClient{}, mq)
	}
}

func MQTTClientFromContext(ctx context.Context) (*MqttClient, bool) {
	v, ok := ctx.Value(CtxMqttClient{}).(*MqttClient)
	return v, ok
}

func MustMQTTClientFromContext(ctx context.Context) *MqttClient {
	v, ok := MQTTClientFromContext(ctx)
	must.BeTrue(ok)
	return v
}
