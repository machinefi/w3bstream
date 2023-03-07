package wasm

import (
	"context"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/conf/tls"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/pkg/errors"
)

type MqttBrokerScheme string

const (
	MqttBrokerScheme_TCP   MqttBrokerScheme = "tcp"
	MqttBrokerScheme_MQTT  MqttBrokerScheme = "mqtt"
	MqttBrokerScheme_MQTTs MqttBrokerScheme = "mqtts"
)

var brokers = mapx.New[string, *MqttBroker]()

type MqttBroker struct {
	Scheme   MqttBrokerScheme `json:"scheme"` // Scheme support tcp only TODO support other protocol
	Host     string           `json:"host"`
	Port     uint16           `json:"port"`
	Username string           `json:"username,omitempty"`
	Password string           `json:"password,omitempty"`
	Topics   []string         `json:"topics"`
	TLS      *tls.X509KeyPair `json:"tls,omitempty"`
	server   string
	cli      mqtt.Client
}

func (b *MqttBroker) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__PROJECT_MQTT_BROKER
}

func (b *MqttBroker) WithContext(ctx context.Context) context.Context {
	return WithMqttBroker(ctx, b)
}

func (b *MqttBroker) Init(ctx context.Context) error {
	_, l := conflog.FromContext(ctx).Start(ctx)
	defer l.End()

	b.server = (&types.Endpoint{
		Scheme:   string(b.Scheme),
		Hostname: b.Host,
		Port:     b.Port,
		Username: b.Username,
		Password: types.Password(b.Password),
	}).String()
	if err := b.TLS.Init(); err != nil {
		return err
	}

	l = l.WithValues("broker", b.server)
	if _b, ok := brokers.Load(b.server); ok {
		l.Warn(errors.New("broker already subscribing"))
		*b = *_b
		return nil
	}
	brokers.Store(b.server, b)

	prj := types.MustProjectFromContext(ctx)
	hdl := types.MustMqttMsgHandlerFromContext(ctx)
	l = l.WithValues("prj", prj.Name)

	cli := mqtt.NewClient(
		mqtt.NewClientOptions().
			AddBroker(b.server).
			SetClientID(prj.Name).
			SetKeepAlive(time.Minute).
			SetDefaultPublishHandler(hdl).
			SetPingTimeout(time.Second).
			SetWriteTimeout(time.Second).
			SetTLSConfig(b.TLS.TLSConfig()),
	)

	if tok := cli.Connect(); tok.Wait() && tok.Error() != nil {
		l.Error(tok.Error())
		return tok.Error()
	}
	b.cli = cli

	topics := map[string]struct{}{}
	for _, topic := range b.Topics {
		topics[topic] = struct{}{}
	}

	for topic := range topics {
		l := l.WithValues("topic", topic)
		if tok := cli.Subscribe(topic, 0, nil); tok.Wait() && tok.Error() != nil {
			l.Error(tok.Error())
			return tok.Error()
		}
		l.WithValues("prj", prj.Name, "topic", topic).Info("start subscribing")
	}
	return nil
}

func (b *MqttBroker) Uninit() {
	for _, topic := range b.Topics {
		b.cli.Unsubscribe(topic)
	}
	b.cli.Disconnect(1)
	brokers.Remove(b.server)
}

func (b *MqttBroker) PublishWithTopic(ctx context.Context, topic string, payload []byte) error {
	_, l := conflog.FromContext(ctx).Start(ctx)
	defer l.End()

	l = l.WithValues(
		"broker", b.server,
		"topic", topic,
		"payload", string(payload),
		"app", types.MustAppletFromContext(ctx).Name,
	)
	l.Info("start send")
	if tok := b.cli.Publish(topic, byte(enums.MQTT_QOS__ONCE), false, payload); tok.Wait() && tok.Error() != nil {
		l.Error(tok.Error())
		return tok.Error()
	}
	l.Info("sent")
	return nil
}
