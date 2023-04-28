package applet

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/applet"))

func init() {
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreateApplet{}))
	Root.Register(kit.NewRouter(&GetApplet{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ListApplet{}))
	Root.Register(kit.NewRouter(&RemoveApplet{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &BatchRemoveApplet{}))
	Root.Register(kit.NewRouter(&UpdateApplet{}))
}
