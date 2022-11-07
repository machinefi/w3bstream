package global

import (
	"context"
	"os"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	confapp "github.com/machinefi/w3bstream/pkg/depends/conf/app"
	confhttp "github.com/machinefi/w3bstream/pkg/depends/conf/http"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	confjwt "github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	confpostgres "github.com/machinefi/w3bstream/pkg/depends/conf/postgres"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq/mem_mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/migration"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

// global vars

var (
	postgres      = &confpostgres.Endpoint{Database: models.DB}
	monitorDB     = &confpostgres.Endpoint{Database: models.MonitorDB}
	mqtt          = &confmqtt.Broker{}
	server        = &confhttp.Server{}
	jwt           = &confjwt.Jwt{}
	logger        = &conflog.Log{Name: "srv-demo"}
	std           = conflog.Std()
	uploadConf    = &types.UploadConfig{}
	ethClientConf = &types.ETHClientConfig{}

	App *confapp.Ctx
)

var (
	tasks  mq.TaskManager
	worker *mq.TaskWorker
)

func init() {
	name := os.Getenv(consts.EnvProjectName)
	if name == "" {
		name = "srv-applet-mgr"
	}
	os.Setenv(consts.EnvProjectName, name)
	logger.Name = name
	App = confapp.New(
		confapp.WithName(name),
		confapp.WithRoot(".."),
		confapp.WithVersion("0.0.1"),
		confapp.WithLogger(conflog.Std()),
	)
	App.Conf(postgres, monitorDB, server, jwt, logger, mqtt, uploadConf, ethClientConf)

	confhttp.RegisterCheckerBy(postgres, mqtt, server)
	std.(conflog.LevelSetter).SetLevel(conflog.InfoLevel)

	tasks = mem_mq.New(0)
	worker = mq.NewTaskWorker(tasks, mq.WithWorkerCount(3), mq.WithChannel(name))

	WithContext = contextx.WithContextCompose(
		types.WithDBExecutorContext(postgres),
		types.WithMonitorDBExecutorContext(monitorDB),
		types.WithPgEndpointContext(postgres),
		types.WithLoggerContext(conflog.Std()),
		types.WithMqttBrokerContext(mqtt),
		types.WithUploadConfigContext(uploadConf),
		confid.WithSFIDGeneratorContext(confid.MustNewSFIDGenerator()),
		confjwt.WithConfContext(jwt),
		types.WithETHClientConfigContext(ethClientConf),
		types.WithTaskWorkerContext(worker),
		types.WithTaskBoardContext(mq.NewTaskBoard(tasks)),
	)
}

var WithContext contextx.WithContext

func Server() kit.Transport { return server.WithContextInjector(WithContext) }

func TaskServer() kit.Transport { return worker.WithContextInjector(WithContext) }

func Migrate() {
	ctx, log := conflog.StdContext(context.Background())

	log.Start(ctx, "Migrate")
	defer log.End()
	if err := migration.Migrate(postgres.WithContext(ctx), nil); err != nil {
		log.Panic(err)
	}
	if err := migration.Migrate(monitorDB.WithContext(ctx), nil); err != nil {
		log.Panic(err)
	}
}
