package types

import (
	"context"

	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
)

type (
	CtxMqttBroker struct{}
	CtxChannels   struct{}
)

func WithMqttBroker(ctx context.Context, v *confmqtt.Broker) context.Context {
	return contextx.WithValue(ctx, CtxMqttBroker{}, v)
}

func WithMqttBrokerContext(v *confmqtt.Broker) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxMqttBroker{}, v)
	}
}

func MqttBrokerFromContext(ctx context.Context) (*confmqtt.Broker, bool) {
	v, ok := ctx.Value(CtxMqttBroker{}).(*confmqtt.Broker)
	return v, ok
}

func MustMqttBrokerFromContext(ctx context.Context) *confmqtt.Broker {
	v, ok := MqttBrokerFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithChannels(ctx context.Context, v []*Channel) context.Context {
	return contextx.WithValue(ctx, CtxChannels{}, v)
}

func WithChannelsContext(v []*Channel) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxChannels{}, v)
	}
}

func ChannelsFromContext(ctx context.Context) ([]*Channel, bool) {
	v, ok := ctx.Value(CtxChannels{}).([]*Channel)
	return v, ok
}

func MustChannelsFromContext(ctx context.Context) []*Channel {
	v, ok := ChannelsFromContext(ctx)
	must.BeTrue(ok)
	return v
}
