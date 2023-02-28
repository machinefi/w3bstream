package wasm

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
)

type (
	CtxSQLStore           struct{}
	CtxKVStore            struct{}
	CtxLogger             struct{}
	CtxEnv                struct{}
	CtxEnvPrefix          struct{}
	CtxRedisPrefix        struct{}
	CtxChainClient        struct{}
	CtxRuntimeResource    struct{}
	CtxMqttBroker         struct{}
	CtxMqttClient         struct{}
	CtxMqttInboundMessage struct{}
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

func WithEnvPrefix(ctx context.Context, v string) context.Context {
	return contextx.WithValue(ctx, CtxEnvPrefix{}, v)
}

func WithEnvPrefixContext(v string) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxEnvPrefix{}, v)
	}
}

func EnvPrefixFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(CtxEnvPrefix{}).(string)
	return v, ok
}

func MustEnvPrefixFromContext(ctx context.Context) string {
	v, ok := EnvPrefixFromContext(ctx)
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

func WithRuntimeResource(ctx context.Context, v *mapx.Map[uint32, interface{}]) context.Context {
	return contextx.WithValue(ctx, CtxRuntimeResource{}, v)
}

func WithRuntimeResourceContext(v *mapx.Map[uint32, interface{}]) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxRuntimeResource{}, v)
	}
}

func RuntimeResourceFromContext(ctx context.Context) (*mapx.Map[uint32, interface{}], bool) {
	v, ok := ctx.Value(CtxRuntimeResource{}).(*mapx.Map[uint32, interface{}])
	return v, ok
}

func MustRuntimeResourceFromContext(ctx context.Context) *mapx.Map[uint32, interface{}] {
	v, ok := RuntimeResourceFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithMqttBroker(ctx context.Context, v *MqttBroker) context.Context {
	return contextx.WithValue(ctx, CtxMqttBroker{}, v)
}

func WithMqttBrokerContext(v *MqttBroker) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxMqttBroker{}, v)
	}
}

func MqttBrokerFromContext(ctx context.Context) (*MqttBroker, bool) {
	v, ok := ctx.Value(CtxMqttBroker{}).(*MqttBroker)
	return v, ok
}

func MustMqttBrokerFromContext(ctx context.Context) *MqttBroker {
	v, ok := MqttBrokerFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithMqttInboundMessage(ctx context.Context, v mqtt.Message) context.Context {
	return contextx.WithValue(ctx, CtxMqttInboundMessage{}, v)
}

func WithMqttInboundMessageContext(v mqtt.Message) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxMqttInboundMessage{}, v)
	}
}

func MqttInboundMessageFromContext(ctx context.Context) (mqtt.Message, bool) {
	v, ok := ctx.Value(CtxMqttInboundMessage{}).(mqtt.Message)
	return v, ok
}

func MustInboundMqttMessageFromContext(ctx context.Context) mqtt.Message {
	v, ok := MqttInboundMessageFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithMqttClient(ctx context.Context, mq *MqttClient) context.Context {
	return contextx.WithValue(ctx, CtxMqttClient{}, mq)
}

func WithMqttClientContext(mq *MqttClient) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxMqttClient{}, mq)
	}
}

func MqttClientFromContext(ctx context.Context) (*MqttClient, bool) {
	v, ok := ctx.Value(CtxMqttClient{}).(*MqttClient)
	return v, ok
}

func MustMQFromContext(ctx context.Context) *MqttClient {
	v, ok := MqttClientFromContext(ctx)
	must.BeTrue(ok)
	return v
}
