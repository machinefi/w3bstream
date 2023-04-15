package routes

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/openapi"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var RootRouter = kit.NewRouter(httptransport.BasePath("/demo"))

func init() {
	RootRouter.Register(openapi.Router)
}
