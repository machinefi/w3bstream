package wasm

import (
	"context"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"

	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
)

type MqttClientScheme string

const (
	MqttClientScheme_TCP MqttClientScheme = "tcp"
)

var clients = mapx.New[string, bool]()

type MqttClient struct {
	Scheme   MqttClientScheme `json:"scheme,omitempty"` // Scheme support tcp only TODO support other protocol
	Host     string           `json:"host"`
	Port     uint16           `json:"port"`
	Username string           `json:"username,omitempty"`
	Password string           `json:"password,omitempty"`
	cli      mqtt.Client
	server   string
}

func (m *MqttClient) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__PROJECT_MQTT_CLIENT
}

func (m *MqttClient) WithContext(ctx context.Context) context.Context {
	return WithMqttClient(ctx, m)
}

func (m *MqttClient) Init(ctx context.Context) error {
	_, l := conflog.FromContext(ctx).Start(ctx)
	defer l.End()

	m.server = (&types.Endpoint{
		Scheme:   "tcp",
		Hostname: m.Host,
		Port:     m.Port,
		Username: m.Username,
		Password: types.Password(m.Password),
	}).String()

	l = l.WithValues("broker", m.server)
	if _, ok := clients.Load(m.server); ok {
		l.Warn(errors.New("broker's client already subscribing"))
		return nil
	}
	brokers.Store(m.server, true)

	prj := types.MustProjectFromContext(ctx)
	l = l.WithValues("project", prj.Name)

	cli := mqtt.NewClient(
		mqtt.NewClientOptions().
			AddBroker(m.server).
			SetClientID(prj.Name).
			SetKeepAlive(time.Minute).
			SetPingTimeout(time.Second),
	)

	if tok := cli.Connect(); tok.Wait() && tok.Error() != nil {
		l.Error(tok.Error())
		return tok.Error()
	}
	m.cli = cli

	return nil
}

func (m *MqttClient) UnInit() {
	m.cli.Disconnect(1)
	clients.Remove(m.server)
}

func (m *MqttClient) PublishWithTopic(ctx context.Context, topic string, payload interface{}) error {
	_, l := conflog.FromContext(ctx).Start(ctx)
	defer l.End()

	l = l.WithValues("broker", m.server, "topic", topic, "payload", payload)
	l.Info("start sending")
	if tok := m.cli.Publish(topic, byte(enums.MQTT_QOS__ONCE), false, payload); tok.Wait() && tok.Error() != nil {
		l.Error(tok.Error())
		return tok.Error()
	}
	return nil
}
