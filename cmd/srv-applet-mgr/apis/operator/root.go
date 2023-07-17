package operator

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

var Root = kit.NewRouter(httptransport.Group("/operator"))

func init() {
	Root.Register(kit.NewRouter(&CreateOperator{}))
	Root.Register(kit.NewRouter(&RemoveOperator{}))
	Root.Register(kit.NewRouter(&ListOperator{}))

	access_key.RouterRegister(Root, enums.ApiGroupOperator, enums.ApiGroupOperatorDesc)
}
