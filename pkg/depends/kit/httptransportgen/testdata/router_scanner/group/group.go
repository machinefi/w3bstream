package group

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var (
	Router      = kit.NewRouter(httptransport.Group("/group"))
	HeathRouter = kit.NewRouter(&Health{})
)

func init() {
	Router.Register(HeathRouter)
}

type Health struct {
	httpx.MethodHead
}

func (Health) MiddleOperators() kit.MiddleOperators {
	return kit.MiddleOperators{
		httptransport.Group("/health"),
	}
}

func (*Health) Output(ctx context.Context) (result interface{}, err error) {
	return
}
