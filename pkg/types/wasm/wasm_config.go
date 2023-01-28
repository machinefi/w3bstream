package wasm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/enums"
)

func NewConfigurationByType(t enums.ConfigType) Configuration {
	switch t {
	case enums.CONFIG_TYPE__PROJECT_SCHEMA:
		return &Schema{}
	case enums.CONFIG_TYPE__INSTANCE_CACHE:
		return &Cache{}
	case enums.CONFIG_TYPE__PROJECT_ENV:
		return &Env{}
	default:
		return nil
	}
}

type Configuration interface {
	ConfigType() enums.ConfigType
	WithContext(context.Context) context.Context
}
