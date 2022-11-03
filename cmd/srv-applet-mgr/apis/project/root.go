package project

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/project"))

func init() {
	Root.Register(kit.NewRouter(&CreateProject{}))
	Root.Register(kit.NewRouter(&GetProjectByProjectID{}))
	Root.Register(kit.NewRouter(&ListProject{}))
	Root.Register(kit.NewRouter(&RemoveProject{}))
}
