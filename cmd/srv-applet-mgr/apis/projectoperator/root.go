package projectoperator

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/project_operator"))

func init() {
	Root.Register(kit.NewRouter(&CreateProjectOperator{}))
	Root.Register(kit.NewRouter(&RemoveProjectOperator{}))
	Root.Register(kit.NewRouter(&GetProjectOperator{}))
}
