package applet

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var (
	Root     = kit.NewRouter(httptransport.Group("/applet"))
	Provider = kit.NewRouter(httptransport.Group("/applet/x"))
)

func init() {
	Root.Register(kit.NewRouter(&RemoveApplet{}))
	Root.Register(kit.NewRouter(&UpdateApplet{}))
	Root.Register(kit.NewRouter(&GetApplet{}))
	Root.Register(kit.NewRouter(&GetAppletDetail{}))
	Provider.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreateApplet{}))
	Provider.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ListApplet{}))
	Provider.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ListAppletDetail{}))

}
