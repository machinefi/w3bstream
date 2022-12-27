package project

import (
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
)

type Env struct {
	prefix string
	values *mapx.Map[string, string]
	Values [][2]string `json:"values"`
}

func (env *Env) init(prefix string) {
	env.prefix = prefix
	env.values = mapx.New[string, string]()

	for _, pair := range env.Values {
		if pair[0] != "" {
			env.values.Store(prefix+"__"+pair[0], pair[1])
		}
	}
}

func (env *Env) Prefix() string { return env.prefix }

func (env *Env) Get(k string) (v string, exists bool) {
	return env.values.Load(env.prefix + "__" + k)
}

func (env *Env) Set(k, v string) {
	env.values.Store(env.prefix+"__"+k, v)
}
