package routes

import (
	"context"
	"net/url"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var RedirectRouter = kit.NewRouter(httptransport.Group("/redirect"))

func init() {
	RootRouter.Register(RedirectRouter)
	RedirectRouter.Register(kit.NewRouter(Redirect{}))
	RedirectRouter.Register(kit.NewRouter(RedirectWhenError{}))
}

type Redirect struct {
	httpx.MethodGet
}

func (Redirect) Output(ctx context.Context) (interface{}, error) {
	return httpx.RedirectWithStatusFound(&url.URL{
		Path: "/other",
	}), nil
}

type RedirectWhenError struct {
	httpx.MethodPost
}

func (RedirectWhenError) Output(ctx context.Context) (interface{}, error) {
	return nil, httpx.RedirectWithStatusMovedPermanently(&url.URL{
		Path: "/other",
	})
}
