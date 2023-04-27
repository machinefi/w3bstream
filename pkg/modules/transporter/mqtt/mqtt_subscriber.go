package mqtt

import (
	"context"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/modules/transporter/proxy"
)

type subscriber struct {
	cli   *confmqtt.Client
	topic string
}

func (s *subscriber) subscribing(ctx context.Context) error {
	return s.cli.Subscribe(func(c mqtt.Client, msg mqtt.Message) {
		ev := &eventpb.Event{}
		if err := json.Unmarshal(msg.Payload(), ev); err != nil {
			// log error
			// how to notify publisher this error?
		}
		if _, err := proxy.Forward(ctx, s.topic, ev); err != nil {
			// log error
			// how to notify publisher this error?
		}
	})
}
