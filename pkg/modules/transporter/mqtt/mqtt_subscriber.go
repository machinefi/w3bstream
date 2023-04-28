package mqtt

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"

	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
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
		l := conflog.FromContext(ctx)
		if err := proto.Unmarshal(msg.Payload(), ev); err != nil {
			l.Error(err)
		}
		l.WithValues(
			"id", ev.Header.GetEventId(),
			"type", ev.Header.GetEventType(),
			"time", ev.Header.GetPubTime(),
			"data", len(ev.Payload),
		).Info("mqtt")
		if _, err := proxy.Forward(ctx, s.topic, ev); err != nil {
			// log error
			// how to notify publisher this error?
		}
	})
}
