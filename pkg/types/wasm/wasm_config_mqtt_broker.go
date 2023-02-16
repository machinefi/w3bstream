package wasm

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
)

type MqttBrokerScheme string

const (
	MqttBrokerScheme_TCP MqttBrokerScheme = "tcp"
)

type MqttBroker struct {
	Scheme   MqttBrokerScheme `json:"scheme,omitempty"` // Scheme support tcp only TODO support other protocol
	Host     string           `json:"host"`
	Port     uint16           `json:"port"`
	Username string           `json:"username,omitempty"`
	Password string           `json:"password,omitempty"`
	Topics   []string         `json:"topics"`
	broker   *confmqtt.Broker
}

func (b *MqttBroker) Broker() *confmqtt.Broker {
	return b.broker
}

func (b *MqttBroker) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__PROJECT_MQTT_BROKER
}

func (b *MqttBroker) WithContext(ctx context.Context) context.Context {
	return WithMqttBroker(ctx, b.broker)
}

func (b *MqttBroker) Init(ctx context.Context) error {
	_, l := conflog.FromContext(ctx).Start(ctx, "wasm.MqttBroker.Init")
	defer l.End()

	server := types.Endpoint{
		Scheme:   "tcp",
		Hostname: b.Host,
		Port:     b.Port,
		Username: b.Username,
		Password: types.Password(b.Password),
	}
	broker := &confmqtt.Broker{Server: server}
	broker.SetDefault()

	err := broker.Init()
	if err != nil {
		l.Error(err)
		return err
	}
	b.broker = broker

	prj := types.MustProjectFromContext(ctx)

	cli, err := broker.Client(prj.Name)
	if err != nil {
		l.Error(err)
		return err
	}

	hdl := types.MustMqttMsgHandlerFromContext(ctx)
	if hdl == nil {
		hdl = func(_ mqtt.Client, msg mqtt.Message) {
			l.WithValues(
				"cid", cli.Cid(),
				"msg_id", msg.MessageID(),
				"topic", msg.Topic(),
				"payload", msg.Payload(),
			).Debug("default handler")
		}
	}

	topics := map[string]struct{}{}
	for _, topic := range b.Topics {
		topics[topic] = struct{}{}
	}

	for topic := range topics {
		go func(topic string) {
			// 	c := cli.WithTopic(topic)
			// 	l := l.WithValues("cid", cli.Cid(), "topic", topic)
			// 	defer c.Disconnect()
			// 	if err := c.Subscribe(hdl); err != nil {
			// 		l.Error(err)
			// 	}
			// 	l.Info("stop subscribing")
		}(topic)
		l.Info("start mqtt subscribing")
	}

	return nil
}
