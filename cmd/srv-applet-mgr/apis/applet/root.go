package applet

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/applet"))

func init() {
	Root.Register(kit.NewRouter(&CreateApplet{}))
	Root.Register(kit.NewRouter(&GetApplet{}))
	Root.Register(kit.NewRouter(&ListApplet{}))
	Root.Register(kit.NewRouter(&RemoveApplet{}))
	Root.Register(kit.NewRouter(&UpdateApplet{}))
}
