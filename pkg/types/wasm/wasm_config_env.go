package wasm

import (
	"context"
	"os"

	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Env struct {
	prefix string
	Env    [][2]string `json:"env"`
}

func (env *Env) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__PROJECT_ENV
}

func (env *Env) WithContext(ctx context.Context) context.Context {
	return WithEnv(ctx, env)
}

func (env *Env) Key(k string) string { return env.prefix + "__" + k }

func (env *Env) Get(k string) (v string, exists bool) {
	return os.LookupEnv(env.Key(k))
}

func (env *Env) Init(ctx context.Context) (err error) {
	env.prefix = types.MustProjectFromContext(ctx).Name + "__"

	defer func() {
		if err != nil {
			_ = env.Uninit(nil)
		}
	}()

	for _, pair := range env.Env {
		if pair[0] == "" {
			continue
		}
		if err = os.Setenv(env.prefix+pair[0], pair[1]); err != nil {
			return
		}
	}
	return nil
}

func (env *Env) Uninit(_ context.Context) error {
	for _, pair := range env.Env {
		_ = os.Unsetenv(env.prefix + pair[0])
	}
	return nil
}
