package operator

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/operator"))

func init() {
	Root.Register(kit.NewRouter(&CreateOperator{}))
	Root.Register(kit.NewRouter(&RemoveOperator{}))
	Root.Register(kit.NewRouter(&ListOperator{}))
}
