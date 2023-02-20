package http

import (
	"context"
	"os"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

func RegisterEnvHandler(r *kit.Router) {
	r.Register(EnvRouter)
}

var EnvRouter = kit.NewRouter(&Env{})

type Env struct {
	httpx.MethodGet
	Key string `in:"query" name:"key"`
}

func (r *Env) Path() string { return "/debug/env" }

func (r *Env) Output(ctx context.Context) (interface{}, error) {
	return os.Getenv(r.Key), nil
}
