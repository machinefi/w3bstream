package wasm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/pkg/errors"
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
		panic(errors.Errorf("unknown config type %d", t))
	}
}

type Configuration interface {
	ConfigType() enums.ConfigType
	WithContext(context.Context) context.Context
}
