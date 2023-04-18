package wasm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm/kvdb"
)

func DefaultCache() *Cache {
	return &Cache{Mode: enums.CACHE_MODE__MEMORY}
}

type Cache struct {
	Mode   enums.CacheMode `json:"mode"`
	Prefix string          `json:"prefix,omitempty"`
	KVStore
}

func (c *Cache) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__INSTANCE_CACHE
}

func (c *Cache) WithContext(ctx context.Context) context.Context {
	return WithKVStore(ctx, c)
}

func (c *Cache) Init(ctx context.Context) error {
	switch c.Mode {
	case enums.CACHE_MODE__REDIS:
		c.KVStore = kvdb.NewRedisDB(types.MustRedisEndpointFromContext(ctx))
	default:
		c.KVStore = kvdb.NewMemDB()
	}
	return nil
}
