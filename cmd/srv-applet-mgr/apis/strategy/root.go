package strategy

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var (
	Root     = kit.NewRouter(httptransport.Group("/strategy"))
	Provider = kit.NewRouter(httptransport.Group("/strategy/x"))
)

func init() {
	Provider.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreateStrategy{}))
	Root.Register(kit.NewRouter(&UpdateStrategy{}))
	Root.Register(kit.NewRouter(&GetStrategy{}))
	Provider.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ListStrategy{}))
	Provider.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ListStrategyDetail{}))
	Provider.Register(kit.NewRouter(&middleware.ProjectProvider{}, &RemoveStrategy{}))
}
