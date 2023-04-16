package event

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/event"))

func init() {
	Root.Register(kit.NewRouter(&BatchEventHandleProxy{}))
	Root.Register(kit.NewRouter(&jwt.Auth{}, &middleware.ContextPublisherAuth{}, &EventHandle{}))
}
