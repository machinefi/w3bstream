package wasm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
)

func DefaultMQClient() *MqttClient {
	return &MqttClient{}
}

type MqttClient struct {
	*mqtt.Client
}

func (m *MqttClient) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__PROJECT_MQTT
}

func (m *MqttClient) Init(ctx context.Context) error {
	if m.Client != nil {
		return nil
	}
	var (
		err error
		brk = types.MustMqttBrokerFromContext(ctx)
		prj = types.MustProjectFromContext(ctx)
	)
	m.Client, err = brk.Client(prj.Name)
	return err
}

func (m *MqttClient) WithContext(ctx context.Context) context.Context {
	return WithMQTTClient(ctx, m)
}
