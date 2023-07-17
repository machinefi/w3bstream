package event

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

var (
	Root = kit.NewRouter(httptransport.Group("/event"), &middleware.EventReqRateLimit{})
)

func init() {
	Root.Register(kit.NewRouter(&HandleEvent{}))

	access_key.RouterRegister(Root, enums.ApiGroupEvent, enums.ApiGroupEventDesc)
}
