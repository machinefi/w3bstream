package main

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httpswaggergen/testdata/router_scanner/auth"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httpswaggergen/testdata/router_scanner/group"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

type Get struct {
	httpx.MethodGet `path:"/:id"`

	ID string `name:"id" in:"path"`
}

func (get Get) Output(ctx context.Context) (result interface{}, err error) {
	return
}

var Router = kit.NewRouter(httptransport.Group("/root"))

func main() {
	Router.Register(group.Router)
	Router.Register(kit.NewRouter(auth.Auth{}, Get{}))

	ht := &httptransport.HttpTransport{
		Port: 8080,
	}
	ht.SetDefault()

	kit.Run(Router, ht)
}
