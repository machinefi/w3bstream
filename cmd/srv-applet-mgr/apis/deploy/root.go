package deploy

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/deploy"))

func init() {
	Root.Register(kit.NewRouter(&CreateAndStartInstance{}))
	Root.Register(kit.NewRouter(&GetInstanceByInstanceID{}))
	Root.Register(kit.NewRouter(&GetInstanceByAppletID{}))
	Root.Register(kit.NewRouter(&ControlInstance{}))
	Root.Register(kit.NewRouter(&RemoveInstance{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &BatchRemoveInstance{}))
}
