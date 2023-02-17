package wasm

import (
	"context"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/enums"
)

func NewConfigurationByType(t enums.ConfigType) (Configuration, error) {
	switch t {
	case enums.CONFIG_TYPE__PROJECT_SCHEMA:
		return &Schema{}, nil
	case enums.CONFIG_TYPE__INSTANCE_CACHE:
		return &Cache{}, nil
	case enums.CONFIG_TYPE__PROJECT_ENV:
		return &Env{}, nil
	case enums.CONFIG_TYPE__PROJECT_MQTT_BROKER:
		return &MqttBroker{}, nil
	default:
		return nil, errors.Errorf("invalid config type: %d", t)
	}
}

type Configuration interface {
	ConfigType() enums.ConfigType
	WithContext(context.Context) context.Context
}

// Init: init wasm configuration
func Init(ctx context.Context, c Configuration) error {
	switch v := c.(type) {
	case types.Initializer:
		v.Init()
		return nil
	case types.ValidatedInitializer:
		return v.Init()
	case types.InitializerWithContext:
		v.Init(ctx)
		return nil
	case types.ValidatedInitializerWithContext:
		return v.Init(ctx)
	default:
		return nil
	}
}
