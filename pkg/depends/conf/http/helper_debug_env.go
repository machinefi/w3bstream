package http

import (
	"context"
	"os"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var (
	EnvQueryRouter = kit.NewRouter(&EnvQuery{})
	EnvSetRouter   = kit.NewRouter(&EnvSet{})
)

func RegisterEnvRouters(r *kit.Router) {
	r.Register(EnvQueryRouter)
	r.Register(EnvSetRouter)
}

type EnvQuery struct {
	httpx.MethodGet
	K string `in:"query" name:"k,omitempty"`
}

func (r *EnvQuery) Path() string { return "/env" }

func (r *EnvQuery) Output(ctx context.Context) (interface{}, error) {
	if r.K != "" {
		return map[string]string{
			r.K: os.Getenv(r.K),
		}, nil
	}
	vars := os.Environ()
	m := make(map[string]string)
	prefix := ""
	if os.Getenv(consts.EnvProjectName) != "" {
		prefix = os.Getenv(consts.EnvProjectName)
		prefix = strings.ToUpper(strings.Replace(prefix, "-", "_", -1))
	}

	for _, v := range vars {
		pair := strings.SplitN(v, "=", 2)
		if len(pair) != 2 {
			continue
		}
		key, val := pair[0], pair[1]
		if prefix != "" && strings.HasPrefix(key, prefix) {
			m[key] = val
		}
	}
	return m, nil
}

type EnvSet struct {
	httpx.MethodPost
	K string `in:"query"           name:"key"`
	V string `in:"query,omitempty" name:"key"`
}

func (r *EnvSet) Path() string { return "/env" }

func (r *EnvSet) Output(ctx context.Context) (interface{}, error) {
	return nil, os.Setenv(r.K, r.V)
}
