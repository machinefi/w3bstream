package apis

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/account"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/applet"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/deploy"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/event"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/login"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/monitor"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/project"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/project_config"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/publisher"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/resource"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/strategy"
	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
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
	RouterV0.Register(RouterAuth)
	{
		RouterAuth.Register(account.Root)
		RouterAuth.Register(project.Root)
		RouterAuth.Register(project_config.Root)
		RouterAuth.Register(applet.Root)
		RouterAuth.Register(deploy.Root)
		RouterAuth.Register(publisher.Root)
		RouterAuth.Register(strategy.Root)
		RouterAuth.Register(monitor.Root)
		RouterAuth.Register(resource.Root)
	}
}
