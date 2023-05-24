package global

import (
	"context"
	"os"
	"time"

	_ "github.com/machinefi/w3bstream/cmd/srv-applet-mgr/types"
	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
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
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/client"
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
	Context     context.Context

	tasks  mq.TaskManager
	worker *mq.TaskWorker

	proxy *client.Client // proxy client for forward mqtt event

	db        = &confpostgres.Endpoint{Database: models.DB}
	monitordb = &confpostgres.Endpoint{Database: models.MonitorDB}
	wasmdb    = &base.Endpoint{}

	ServerMgr   = &confhttp.Server{}
	ServerEvent = &confhttp.Server{} // serverEvent support event http transport

	fs  filesystem.FileSystemOp
	std = conflog.Std().(conflog.LevelSetter).SetLevel(conflog.DebugLevel)
)

func init() {
	// TODO config struct should be defined outside this method and impl it's Init() interface{}
	// TODO split this init too long
	config := &struct {
		Postgres     *confpostgres.Endpoint
		MonitorDB    *confpostgres.Endpoint
		WasmDB       *base.Endpoint
		MqttBroker   *confmqtt.Broker
		Redis        *confredis.Redis
		Server       *confhttp.Server
		Jwt          *confjwt.Jwt
		Logger       *conflog.Log
		UploadConf   *types.UploadConfig
		EthClient    *types.ETHClientConfig
		WhiteList    *types.WhiteList
		ServerEvent  *confhttp.Server
		FileSystem   *types.FileSystem
		AmazonS3     *amazonS3.AmazonS3
		LocalFS      *local.LocalFileSystem
		WasmDBConfig *types.WasmDBConfig
	}{
		Postgres:     db,
		MonitorDB:    monitordb,
		WasmDB:       wasmdb,
		MqttBroker:   &confmqtt.Broker{},
		Redis:        &confredis.Redis{},
		Server:       ServerMgr,
		Jwt:          &confjwt.Jwt{},
		Logger:       &conflog.Log{},
		UploadConf:   &types.UploadConfig{},
		EthClient:    &types.ETHClientConfig{},
		WhiteList:    &types.WhiteList{},
		ServerEvent:  ServerEvent,
		FileSystem:   &types.FileSystem{},
		AmazonS3:     &amazonS3.AmazonS3{},
		LocalFS:      &local.LocalFileSystem{},
		WasmDBConfig: &types.WasmDBConfig{},
	}

	name := os.Getenv(consts.EnvProjectName)
	if name == "" {
		name = "srv-applet-mgr"
	}
	_ = os.Setenv(consts.EnvProjectName, name)

	group := os.Getenv(consts.EnvResourceGroup)
	if group == "" {
		group = "srv-applet-mgr"
	}
	_ = os.Setenv(consts.EnvResourceGroup, group)

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

	proxy = &client.Client{Port: uint16(ServerEvent.Port), Timeout: 10 * time.Second}
	proxy.SetDefault()

	WithContext = contextx.WithContextCompose(
		types.WithMgrDBExecutorContext(config.Postgres),
		types.WithMonitorDBExecutorContext(config.MonitorDB),
		types.WithWasmDBEndpointContext(config.WasmDB),
		types.WithRedisEndpointContext(config.Redis),
		types.WithLoggerContext(std),
		conflog.WithLoggerContext(std),
		types.WithUploadConfigContext(config.UploadConf),
		types.WithMqttBrokerContext(config.MqttBroker),
		confid.WithSFIDGeneratorContext(confid.MustNewSFIDGenerator()),
		confjwt.WithConfContext(config.Jwt),
		types.WithTaskWorkerContext(worker),
		types.WithTaskBoardContext(mq.NewTaskBoard(tasks)),
		types.WithETHClientConfigContext(config.EthClient),
		types.WithWhiteListContext(config.WhiteList),
		types.WithFileSystemOpContext(fs),
		types.WithProxyClientContext(proxy),
		types.WithWasmDBConfigContext(config.WasmDBConfig),
	)
	Context = WithContext(context.Background())
}

func Server() kit.Transport { return ServerMgr.WithContextInjector(WithContext) }

func TaskServer() kit.Transport { return worker.WithContextInjector(WithContext) }

func EventServer() kit.Transport { return ServerEvent.WithContextInjector(WithContext) }

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
