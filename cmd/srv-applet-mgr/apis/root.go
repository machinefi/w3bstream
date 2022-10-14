package apis

import (
	"github.com/iotexproject/Bumblebee/conf/jwt"
	"github.com/iotexproject/Bumblebee/kit/httptransport"
	"github.com/iotexproject/Bumblebee/kit/kit"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/account"
	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/applet"
	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/deploy"
	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/event"
	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/login"
	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/monitor"
	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/project"
	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/publisher"
)

var (
	name         = "srv-applet-mgr"
	Root         = kit.NewRouter(httptransport.Group("/"))
	RouterServer = kit.NewRouter(httptransport.Group("/" + name))
	RouterV0     = kit.NewRouter(httptransport.Group("/v0"))
	RouterAuth   = kit.NewRouter(&jwt.Auth{}, &middleware.ContextAccountAuth{})
)

func init() {
	Root.Register(RouterServer)
	RouterServer.Register(RouterV0)

	RouterV0.Register(login.Root)
	RouterV0.Register(event.Root)
	RouterV0.Register(monitor.Root)
	RouterV0.Register(RouterAuth)
	{
		RouterAuth.Register(account.Root)
		RouterAuth.Register(project.Root)
		RouterAuth.Register(applet.Root)
		RouterAuth.Register(deploy.Root)
		RouterAuth.Register(publisher.Root)
	}
}
