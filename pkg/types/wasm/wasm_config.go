package wasm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/pkg/errors"

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
	case enums.CONFIG_TYPE__PROJECT_MQTT:
		return &MqttClient{}, nil
	default:
		return nil, errors.Errorf("invalid config type: %d", t)
	}
}

type Configuration interface {
	ConfigType() enums.ConfigType
	WithContext(context.Context) context.Context
}

// ConfigurationWithInit support recursive initialize
type ConfigurationWithInit interface {
	Configuration
	types.ValidatedInitializerWith
}

// ConfigurationWithUninit support recursive uninitialize
type ConfigurationWithUninit interface {
	Configuration
	Uninit(context.Context) error
}

func InitConfiguration(ctx context.Context, c Configuration) error {
	if canBeInit, ok := c.(ConfigurationWithInit); ok {
		return canBeInit.Init(ctx)
	}
	return nil
}

func UninitConfiguration(ctx context.Context, c Configuration) error {
	if canBeUninit, ok := c.(ConfigurationWithUninit); ok {
		return canBeUninit.Uninit(ctx)
	}
	return nil
}
