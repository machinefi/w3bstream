package projectoperator

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/project_operator"))

func init() {
	Root.Register(kit.NewRouter(&Create{}))
	Root.Register(kit.NewRouter(&Remove{}))
	Root.Register(kit.NewRouter(&Get{}))
}
