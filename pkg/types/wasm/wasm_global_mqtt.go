package wasm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/types"
)

type MqttClient struct {
	cli *mqtt.Client
}

func (m *MqttClient) GlobalConfigType() ConfigType { return ConfigMqttClient }

func (m *MqttClient) Init(parent context.Context) error {
	var (
		broker = types.MustMqttBrokerFromContext(parent)
		prj    = types.MustProjectFromContext(parent)
	)
	cli, err := broker.Client(prj.Name)
	if err != nil {
		return err
	}
	m.cli = cli
	return nil
}

func (m *MqttClient) WithContext(ctx context.Context) context.Context {
	return WithMQTTClient(ctx, m.cli)
}
