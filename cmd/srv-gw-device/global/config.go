package global

import (
	"context"
	"os"

	"github.com/machinefi/w3bstream/cmd/srv-gw-device/types"
	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	confapp "github.com/machinefi/w3bstream/pkg/depends/conf/app"
	confhttp "github.com/machinefi/w3bstream/pkg/depends/conf/http"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	confpostgres "github.com/machinefi/w3bstream/pkg/depends/conf/postgres"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/migration"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/models"
)

type config struct {
	ServiceName string
	Postgres    *confpostgres.Endpoint
	Server      *confhttp.Server
	Transporter *types.Transporter
	Logger      *conflog.Log
	Core        *types.Core
}

var (
	postgres = &confpostgres.Endpoint{Database: models.GwDB}
	server   = &confhttp.Server{}

	App            *confapp.Ctx
	WithAppContext contextx.WithContext
)

func init() {
	cfg := &config{
		Postgres:    postgres,
		Server:      server,
		Transporter: &types.Transporter{},
		Logger:      &conflog.Log{},
		Core:        &types.Core{},
	}

	name := os.Getenv(consts.EnvProjectName)
	if name == "" {
		name = "srv-gw-device"
	}
	os.Setenv(consts.EnvProjectName, name)

	App = confapp.New(
		confapp.WithName(name),
		confapp.WithRoot(".."),
		confapp.WithVersion("0.0.1"),
		confapp.WithLogger(conflog.Std()),
	)
	App.Conf(cfg)

	WithAppContext = contextx.WithContextCompose(
		types.WithDBExecutorContext(cfg.Postgres),
		types.WithLoggerContext(conflog.Std()), // TODO impl conflog.Log as Logger and inject cfg.Logger
		types.WithCoreContext(cfg.Core),
	)
}

func Server() kit.Transport {
	return server.WithContextInjector(WithAppContext)
}

func Migrate() {
	ctx, l := conflog.StdContext(context.Background())

	l.Start(ctx, "Migrate")
	defer l.End()

	if err := migration.Migrate(postgres.WithContext(ctx), nil); err != nil {
		l.Panic(err)
	}
}
