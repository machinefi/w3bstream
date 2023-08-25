package xvm

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/xvm"))

func init() {
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreateRisc0VM{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreateProof{}))
}
