package projectoperator

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

var Root = kit.NewRouter(httptransport.Group("/project_operator"))

func init() {
	Root.Register(kit.NewRouter(&CreateProjectOperator{}))
	Root.Register(kit.NewRouter(&RemoveProjectOperator{}))
	Root.Register(kit.NewRouter(&GetProjectOperator{}))

	access_key.RouterRegister(Root, enums.ApiGroupProjectOperator, enums.ApiGroupProjectOperatorDesc)
}
