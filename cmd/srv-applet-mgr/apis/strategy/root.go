package strategy

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

var Root = kit.NewRouter(httptransport.Group("/strategy"))

func init() {
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreateStrategy{}))
	Root.Register(kit.NewRouter(&UpdateStrategy{}))
	Root.Register(kit.NewRouter(&GetStrategy{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ListStrategy{}))
	Root.Register(kit.NewRouter(&RemoveStrategy{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &BatchRemoveStrategy{}))

	access_key.RouterRegister(Root, enums.ApiGroupStrategy, enums.ApiGroupStrategyDesc)
}
