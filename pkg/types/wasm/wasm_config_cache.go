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

	kv KVStore
}

func (c *Cache) Init(parent context.Context) error {
	switch c.Mode {
	case enums.CACHE_MODE__REDIS:
		c.kv = kvdb.NewRedisDB(types.MustRedisEndpointFromContext(parent))
	default:
		c.kv = kvdb.NewMemDB()
	}
	return nil
}

func (c *Cache) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__INSTANCE_CACHE
}

func (c *Cache) WithContext(ctx context.Context) context.Context {
	return WithKVStore(ctx, c.kv)
}
