package strategy

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/strategy"))

func init() {
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreateStrategy{}))
	Root.Register(kit.NewRouter(&UpdateStrategy{}))
	Root.Register(kit.NewRouter(&GetStrategy{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ListStrategy{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &RemoveStrategy{}))
}
