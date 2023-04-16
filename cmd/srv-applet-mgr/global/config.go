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

	db          = &confpostgres.Endpoint{Database: models.DB}
	monitordb   = &confpostgres.Endpoint{Database: models.MonitorDB}
	wasmdb      = &confpostgres.Endpoint{Database: models.WasmDB}
	server      = &confhttp.Server{}
	serverEvent = &confhttp.Server{} // serverEvent support event http transport
)

func init() {
	config := &struct {
		Postgres    *confpostgres.Endpoint
		MonitorDB   *confpostgres.Endpoint
		WasmDB      *confpostgres.Endpoint
		MqttBroker  *confmqtt.Broker
		Redis       *confredis.Redis
		Server      *confhttp.Server
		Jwt         *confjwt.Jwt
		Logger      *conflog.Log
		StdLogger   conflog.Logger
		UploadConf  *types.UploadConfig
		EthClient   *types.ETHClientConfig
		WhiteList   *types.WhiteList
		ServerEvent *confhttp.Server
	}{
		Postgres:    db,
		MonitorDB:   monitordb,
		WasmDB:      wasmdb,
		MqttBroker:  &confmqtt.Broker{},
		Redis:       &confredis.Redis{},
		Server:      server,
		Jwt:         &confjwt.Jwt{},
		Logger:      &conflog.Log{},
		StdLogger:   conflog.Std(),
		UploadConf:  &types.UploadConfig{},
		EthClient:   &types.ETHClientConfig{},
		WhiteList:   &types.WhiteList{},
		ServerEvent: serverEvent,
	}

	name := os.Getenv(consts.EnvProjectName)
	if name == "" {
		name = "srv-applet-mgr"
	}
	os.Setenv(consts.EnvProjectName, name)
	config.Logger.Name = name

	tasks = mem_mq.New(0)
	worker = mq.NewTaskWorker(tasks, mq.WithWorkerCount(3), mq.WithChannel(name))

	App = confapp.New(
		confapp.WithName(name),
		confapp.WithRoot(".."),
		confapp.WithVersion("0.0.1"),
		confapp.WithLogger(conflog.Std()),
	)
	App.Conf(config, worker)

	confhttp.RegisterCheckerBy(config, worker)
	config.StdLogger.(conflog.LevelSetter).SetLevel(conflog.InfoLevel)

	WithContext = contextx.WithContextCompose(
		types.WithMgrDBExecutorContext(config.Postgres),
		types.WithMonitorDBExecutorContext(config.MonitorDB),
		types.WithWasmDBExecutorContext(config.WasmDB),
		types.WithPgEndpointContext(config.Postgres),
		types.WithRedisEndpointContext(config.Redis),
		types.WithLoggerContext(config.StdLogger),
		conflog.WithLoggerContext(config.StdLogger),
		types.WithMqttBrokerContext(config.MqttBroker),
		types.WithUploadConfigContext(config.UploadConf),
		confid.WithSFIDGeneratorContext(confid.MustNewSFIDGenerator()),
		confjwt.WithConfContext(config.Jwt),
		types.WithTaskWorkerContext(worker),
		types.WithTaskBoardContext(mq.NewTaskBoard(tasks)),
		types.WithETHClientConfigContext(config.EthClient),
		types.WithWhiteListContext(config.WhiteList),
	)
}

func Server() kit.Transport { return server.WithContextInjector(WithContext) }

func TaskServer() kit.Transport { return worker.WithContextInjector(WithContext) }

func EventServer() kit.Transport { return serverEvent.WithContextInjector(WithContext) }

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
