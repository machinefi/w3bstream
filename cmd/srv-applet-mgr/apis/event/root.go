package event

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var (
	Root = kit.NewRouter(httptransport.Group("/event"), &middleware.EventReqRateLimit{})
)

func init() {
	Root.Register(kit.NewRouter(&HandleEvent{}))
}
