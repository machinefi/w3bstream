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
	confredis "github.com/machinefi/w3bstream/pkg/depends/conf/redis"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq/mem_mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/migration"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

var (
	App         *confapp.Ctx
	WithContext contextx.WithContext

	tasks  mq.TaskManager
	worker *mq.TaskWorker

	db        = &confpostgres.Endpoint{Database: models.DB}
	monitordb = &confpostgres.Endpoint{Database: models.MonitorDB}
	server    = &confhttp.Server{}
)

func init() {
	config := &struct {
		Postgres   *confpostgres.Endpoint
		MonitorDB  *confpostgres.Endpoint
		MqttBroker *confmqtt.Broker
		Redis      *confredis.Redis
		Server     *confhttp.Server
		Jwt        *confjwt.Jwt
		Logger     *conflog.Log
		StdLogger  conflog.Logger
		UploadConf *types.UploadConfig
		EthClient  *types.ETHClientConfig
	}{
		Postgres:   db,
		MonitorDB:  monitordb,
		MqttBroker: &confmqtt.Broker{},
		Redis:      &confredis.Redis{},
		Server:     server,
		Jwt:        &confjwt.Jwt{},
		Logger:     &conflog.Log{},
		StdLogger:  conflog.Std(),
		UploadConf: &types.UploadConfig{},
		EthClient:  &types.ETHClientConfig{},
	}

	name := os.Getenv(consts.EnvProjectName)
	if name == "" {
		name = "srv-applet-mgr"
	}
	os.Setenv(consts.EnvProjectName, name)
	config.Logger.Name = name
	App = confapp.New(
		confapp.WithName(name),
		confapp.WithRoot(".."),
		confapp.WithVersion("0.0.1"),
		confapp.WithLogger(conflog.Std()),
	)
	App.Conf(config)

	confhttp.RegisterCheckerBy(
		config.Postgres,
		config.MqttBroker,
		config.Redis,
		config.Server,
	)
	config.StdLogger.(conflog.LevelSetter).SetLevel(conflog.InfoLevel)

	tasks = mem_mq.New(0)
	worker = mq.NewTaskWorker(tasks, mq.WithWorkerCount(3), mq.WithChannel(name))

	WithContext = contextx.WithContextCompose(
		types.WithDBExecutorContext(config.Postgres),
		types.WithMonitorDBExecutorContext(config.MonitorDB),
		types.WithPgEndpointContext(config.Postgres),
		types.WithRedisEndpointContext(config.Redis),
		types.WithLoggerContext(conflog.Std()),
		types.WithMqttBrokerContext(config.MqttBroker),
		types.WithUploadConfigContext(config.UploadConf),
		confid.WithSFIDGeneratorContext(confid.MustNewSFIDGenerator()),
		confjwt.WithConfContext(config.Jwt),
		types.WithETHClientConfigContext(config.EthClient),
		types.WithTaskWorkerContext(worker),
		types.WithTaskBoardContext(mq.NewTaskBoard(tasks)),
	)
}

func Server() kit.Transport { return server.WithContextInjector(WithContext) }

func TaskServer() kit.Transport { return worker.WithContextInjector(WithContext) }

func Migrate() {
	ctx, log := conflog.StdContext(context.Background())

	log.Start(ctx, "Migrate")
	defer log.End()
	if err := migration.Migrate(db.WithContext(ctx), nil); err != nil {
		log.Panic(err)
	}
	if err := migration.Migrate(monitordb.WithContext(ctx), nil); err != nil {
		log.Panic(err)
	}
}
