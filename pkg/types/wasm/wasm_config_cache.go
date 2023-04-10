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
}

func (c *Cache) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__INSTANCE_CACHE
}

func (c *Cache) WithContext(ctx context.Context) context.Context {
	return WithKVStore(ctx, c.NewKVStore(ctx))
}

// TODO use KVStore as a member and impl Cache.Init
func (c *Cache) NewKVStore(ctx context.Context) KVStore {
	switch c.Mode {
	case enums.CACHE_MODE__REDIS:
		return kvdb.NewRedisDB(types.MustRedisEndpointFromContext(ctx))
	case enums.CACHE_MODE__MEMORY:
		return kvdb.NewMemDB()
	}
	return nil
}
