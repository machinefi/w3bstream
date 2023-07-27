package wasm

import (
	"context"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/enums"
)

func NewUserConfigurationByType(t enums.ConfigType) (Configuration, error) {
	switch t {
	case enums.CONFIG_TYPE__PROJECT_DATABASE:
		return &Database{}, nil
	case enums.CONFIG_TYPE__INSTANCE_CACHE:
		return &Cache{}, nil
	case enums.CONFIG_TYPE__PROJECT_ENV:
		return &Env{}, nil
	case enums.CONFIG_TYPE__PROJECT_FLOW:
		return &Flow{}, nil
	default:
		return nil, errors.Errorf("invalid config type: %d", t)
	}
}

type Configuration interface {
	ConfigType() enums.ConfigType
	WithContext(context.Context) context.Context
}

type CanBeUninit interface {
	Uninit(context.Context) error
}

type CanBeInit interface {
	Init(context.Context) error
}

func InitConfiguration(parent context.Context, c Configuration) error {
	if canBeInit, ok := c.(CanBeInit); ok {
		return canBeInit.Init(parent)
	}
	return nil
}

func UninitConfiguration(parent context.Context, c Configuration) error {
	if canBeUninit, ok := c.(CanBeUninit); ok {
		return canBeUninit.Uninit(parent)
	}
	return nil
}

type ConfigType string

const (
	ConfigLogger     ConfigType = "LOGGER"
	ConfigMqttClient ConfigType = "MQTT_CLIENT"
	ConfigChains     ConfigType = "CHAINS"
	ConfigMetrics    ConfigType = "METRICS"
)

var ConfigTypes = []ConfigType{
	ConfigLogger,
	ConfigMqttClient,
	ConfigChains,
	ConfigMetrics,
}

func NewGlobalConfigurationByType(t ConfigType) (GlobalConfiguration, error) {
	switch t {
	case ConfigLogger:
		return &Logger{}, nil
	case ConfigMqttClient:
		return &MqttClient{}, nil
	case ConfigChains:
		return &ChainClient{}, nil
	default: // TODO case ConfigMetrics:
		return nil, nil // errors.Errorf("invalid global config type: %d", t)
	}
}

type GlobalConfiguration interface {
	GlobalConfigType() ConfigType
	WithContext(context.Context) context.Context
}

func InitGlobalConfiguration(parent context.Context, c GlobalConfiguration) error {
	if canBeInit, ok := c.(CanBeInit); ok {
		return canBeInit.Init(parent)
	}
	return nil
}
