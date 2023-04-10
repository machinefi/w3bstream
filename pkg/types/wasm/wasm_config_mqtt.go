package wasm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
)

type MqttClient struct {
	*mqtt.Client
}

func (m *MqttClient) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__PROJECT_MQTT
}

// TODO impl MqttClient.Init

func (m *MqttClient) WithContext(ctx context.Context) context.Context {
	mq := types.MustMqttBrokerFromContext(ctx)
	log := types.MustLoggerFromContext(ctx)
	projectName := types.MustProjectFromContext(ctx).ProjectName.Name
	cli, err := mq.Client(projectName)
	if err != nil {
		log.Error(err)
		return ctx
	}
	client := &MqttClient{
		cli,
	}

	return WithMQTTClient(ctx, client)
}
