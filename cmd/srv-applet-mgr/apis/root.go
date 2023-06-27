package apis

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/account"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/account_access"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/applet"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/configuration"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/cronjob"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/deploy"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/event"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/login"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/monitor"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/operator"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/project"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/project_config"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/projectoperator"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/publisher"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/resource"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/strategy"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/traffic_limit"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/version"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/wasmlog"
	confhttp "github.com/machinefi/w3bstream/pkg/depends/conf/http"
	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/openapi"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var (
	RootMgr   = kit.NewRouter(httptransport.Group("/"))
	RootEvent = kit.NewRouter(httptransport.Group("/"))
)

func init() {
	// root router register for applet-mgr
	{
		serve := kit.NewRouter(httptransport.Group("/srv-applet-mgr"))
		RootMgr.Register(serve)
		RootMgr.Register(kit.NewRouter(&version.VersionRouter{}))
		RootMgr.Register(kit.NewRouter(&confhttp.Liveness{}))

		var (
			v0   = kit.NewRouter(httptransport.Group("/v0"))
			auth = kit.NewRouter(&jwt.Auth{}, &middleware.ContextAccountAuth{})
		)
		serve.Register(v0)
		serve.Register(kit.NewRouter(&openapi.OpenAPI{}))

		v0.Register(login.Root)
		v0.Register(account.RegisterRoot)
		v0.Register(configuration.Root)
		v0.Register(auth)

		auth.Register(account.Root)
		auth.Register(account_access.Root)
		auth.Register(project.Root)
		auth.Register(project_config.Root)
		auth.Register(applet.Root)
		auth.Register(deploy.Root)
		auth.Register(publisher.Root)
		auth.Register(strategy.Root)
		auth.Register(monitor.Root)
		auth.Register(cronjob.Root)
		auth.Register(resource.Root)
		auth.Register(wasmlog.Root)
		auth.Register(operator.Root)
		auth.Register(traffic_limit.Root)
		auth.Register(projectoperator.Root)
		auth.Register(event.Root2)
	}

	// root router register for event http transport
	{
		// TODO should use another root name to differentiate with w3bstream core? `/srv-event`
		serve := kit.NewRouter(httptransport.Group("/srv-applet-mgr"))
		RootEvent.Register(serve)
		RootEvent.Register(kit.NewRouter(&version.VersionRouter{}))
		RootEvent.Register(kit.NewRouter(&confhttp.Liveness{}))

		var (
			v0   = kit.NewRouter(httptransport.Group("/v0"))
			auth = kit.NewRouter(&jwt.Auth{}, &middleware.ContextPublisherAuth{})
		)
		serve.Register(kit.NewRouter(&openapi.OpenAPI{}))
		serve.Register(v0)

		v0.Register(auth)

		auth.Register(event.Root)
	}
}
