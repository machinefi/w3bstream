package mqtt

import (
	"context"
	"encoding/json"
	"net/url"
	"path"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/modules/transporter/proxy"
)

type subscriber struct {
	cli   *confmqtt.Client
	topic string
}

func ParseInboundMessage(msg mqtt.Message) (*eventpb.Event, error) {
	topic := msg.Topic()

	parts := strings.Split(topic, "/")

	if len(parts) == 1 {
		ev := &eventpb.Event{}
		err := proto.Unmarshal(msg.Payload(), ev)
		return ev, err
	}

	if len(parts) != 4 && len(parts) != 5 {
		return nil, ErrInvalidTopicParts
	}

	if parts[1] != "push" {
		return nil, ErrNotInboundTopic
	}

	ev := &eventpb.Event{
		Header: &eventpb.Header{
			Token:     parts[2],
			EventType: parts[3],
		},
		Payload: msg.Payload(),
	}
	if len(parts) == 5 {
		values, err := url.ParseQuery(parts[4])
		if err != nil {
			return nil, err
		}
		if v := values.Get("ts"); v != "" {
			ev.Header.PubTime, _ = strconv.ParseInt(v, 10, 64)
		}
		if v := values.Get("id"); v != "" {
			ev.Header.EventId = v
		}
	}
	return ev, nil
}

func (s *subscriber) subscribing(ctx context.Context) error {
	return s.cli.WithTopic(s.topic + "/#").Subscribe(func(c mqtt.Client, msg mqtt.Message) {
		ctx, l := logger.NewSpanContext(ctx, "modules.transporter.mqtt.subscriber.handle")
		defer l.End()

		ev, err := ParseInboundMessage(msg)
		if err != nil {
			l.Error(err)
			return
		}

		l = l.WithValues(
			"topic", msg.Topic(),
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

		topic := path.Join("ack", ret.(*event.EventRsp).PublisherKey)
		cli := s.cli.WithTopic(topic)
		if err = cli.Publish(rsp); err != nil {
			l.Error(errors.Wrap(err, "publish rsp"))
		}
		l.Info("%s: %s", topic, string(rsp))
	})
}

var (
	ErrInvalidTopicParts = errors.New("invalid topic parts")
	ErrNotInboundTopic   = errors.New("not inbound topic")
)
