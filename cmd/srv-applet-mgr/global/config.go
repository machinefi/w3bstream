package global

import (
	"context"
	"os"

	_ "github.com/machinefi/w3bstream/cmd/srv-applet-mgr/types"
	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	confapp "github.com/machinefi/w3bstream/pkg/depends/conf/app"
	"github.com/machinefi/w3bstream/pkg/depends/conf/filesystem"
	"github.com/machinefi/w3bstream/pkg/depends/conf/filesystem/amazonS3"
	"github.com/machinefi/w3bstream/pkg/depends/conf/filesystem/local"
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
	"github.com/machinefi/w3bstream/pkg/enums"
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

	fs filesystem.FileSystemOp
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
		FileSystem  *types.FileSystem
		AmazonS3    *amazonS3.AmazonS3
		LocalFS     *local.LocalFileSystem
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
		FileSystem:  &types.FileSystem{},
		AmazonS3:    &amazonS3.AmazonS3{},
		LocalFS:     &local.LocalFileSystem{},
	}

	name := os.Getenv(consts.EnvProjectName)
	if name == "" {
		name = "srv-applet-mgr"
	}
	_ = os.Setenv(consts.EnvProjectName, name)

	tasks = mem_mq.New(0)
	worker = mq.NewTaskWorker(tasks, mq.WithWorkerCount(3), mq.WithChannel(name))

	App = confapp.New(
		confapp.WithName(name),
		confapp.WithRoot(".."),
		confapp.WithLogger(conflog.Std()),
	)
	App.Conf(config, worker)

	if config.FileSystem.Type == enums.FILE_SYSTEM_MODE__S3 &&
		!config.AmazonS3.IsZero() {
		fs = config.AmazonS3
	} else {
		fs = config.LocalFS
	}

	confhttp.RegisterCheckerBy(config, worker)
	config.StdLogger.(conflog.LevelSetter).SetLevel(conflog.InfoLevel)

	WithContext = contextx.WithContextCompose(
		types.WithMgrDBExecutorContext(config.Postgres),
		types.WithMonitorDBExecutorContext(config.MonitorDB),
		types.WithWasmDBEndpointContext(config.WasmDB),
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
		types.WithFileSystemOpContext(fs),
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
