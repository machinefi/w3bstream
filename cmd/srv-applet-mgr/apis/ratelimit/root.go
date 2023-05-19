package ratelimit

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/traffic"))

func init() {
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreateTrafficRateLimit{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &UpdateTrafficRateLimit{}))
}
