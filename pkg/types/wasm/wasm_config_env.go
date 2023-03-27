package wasm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/enums"
)

func NewEvn(prefix string, kvs ...[2]string) *Env {
	return &Env{
		prefix: prefix,
		values: mapx.New[string, string](),
	}
}

type Env struct {
	prefix string
	values *mapx.Map[string, string]
	Env    [][2]string `json:"env"`
}

func (env *Env) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__PROJECT_ENV
}

func (env *Env) WithContext(ctx context.Context) context.Context {
	prefix, _ := EnvPrefixFromContext(ctx)
	if prefix != "" {
		prefix = prefix + "__"
	}
	env.prefix = prefix

	if env.values == nil {
		env.values = mapx.New[string, string]()
	}
	for _, pair := range env.Env {
		if pair[0] != "" {
			env.Set(pair[0], pair[1])
		}
	}
	return WithEnv(ctx, env)
}

func (env *Env) Prefix() string { return env.prefix }

func (env *Env) Get(k string) (v string, exists bool) {
	return env.values.Load(env.prefix + k)
}

func (env *Env) Set(k, v string) {
	env.values.Store(env.prefix+k, v)
}
