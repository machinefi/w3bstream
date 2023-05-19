package mqtt

import (
	"context"
	"strconv"

	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/types"
)

var topics = mapx.New[string, *subscriber]()

func Subscribe(ctx context.Context, topic string) error {
	broker := types.MustMqttBrokerFromContext(ctx)

	cli, err := broker.Client(strconv.Itoa(topics.Len() + 1))
	if err != nil {
		return status.MqttConnectFailed.StatusErr().WithDesc(err.Error())
	}
	cli.WithTopic(topic)

	sub := &subscriber{topic: topic, cli: cli}
	if err = sub.subscribing(ctx); err != nil {
		broker.Close(cli)
		return status.MqttSubscribeFailed.StatusErr().WithDesc(err.Error())
	}
	if !topics.StoreNX(topic, sub) {
		broker.Close(cli)
		return status.TopicAlreadySubscribed.StatusErr().WithDesc(topic)
	}
	return nil
}

func Stop(ctx context.Context, topic string) {
	broker := types.MustMqttBrokerFromContext(ctx)

	sub, ok := topics.LoadAndRemove(topic)
	if ok {
		broker.Close(sub.cli)
	}
}
