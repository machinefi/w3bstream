package traffic_limit

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/traffic"))

func init() {
	Root.Register(kit.NewRouter(&GetTrafficLimit{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ListTrafficLimit{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreateTrafficLimit{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &UpdateTrafficLimit{}))
}
