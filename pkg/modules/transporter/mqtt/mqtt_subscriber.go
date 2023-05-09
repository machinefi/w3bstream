package mqtt

import (
	"context"
	"encoding/json"
	"path"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

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
	return s.cli.WithTopic(s.topic).Subscribe(func(c mqtt.Client, msg mqtt.Message) {
		ev := &eventpb.Event{}
		l := conflog.FromContext(ctx)
		if err := proto.Unmarshal(msg.Payload(), ev); err != nil {
			l.Error(errors.Wrap(err, "unmarshal req"))
			return
		}
		l = l.WithValues(
			"topic", s.topic,
			"id", ev.Header.GetEventId(),
			"type", ev.Header.GetEventType(),
			"time", ev.Header.GetPubTime(),
			"data", len(ev.Payload),
		)
		ret, err := proxy.Forward(ctx, s.topic, ev)
		if err != nil {
			l.Error(errors.Wrap(err, "forward"))
			return
		}
		rsp, err := json.Marshal(ret)
		if err != nil {
			l.Error(errors.Wrap(err, "marshal rsp"))
			return
		}
		topic := path.Join(s.topic, ev.Header.GetPubId())
		cli := s.cli.WithTopic(topic)
		if err = cli.Publish(rsp); err != nil {
			l.Error(errors.Wrap(err, "publish rsp"))
		}
		l.Info("%s: %s", topic, string(rsp))
	})
}
