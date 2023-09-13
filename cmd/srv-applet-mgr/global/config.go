package global

import (
	"context"
	"os"
	"time"

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
	conflogger "github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	confmq "github.com/machinefi/w3bstream/pkg/depends/conf/mq"
	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	confpostgres "github.com/machinefi/w3bstream/pkg/depends/conf/postgres"
	confrate "github.com/machinefi/w3bstream/pkg/depends/conf/rate_limit"
	confredis "github.com/machinefi/w3bstream/pkg/depends/conf/redis"
	conftracer "github.com/machinefi/w3bstream/pkg/depends/conf/tracer"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/client"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/migration"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/operator/pool"
	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmapi"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm/kvdb"
)

var (
	App         *confapp.Ctx
	WithContext contextx.WithContext
	Context     context.Context

	proxy *client.Client // proxy client for forward mqtt event

	db        = &confpostgres.Endpoint{Database: models.DB}
	monitordb = &confpostgres.Endpoint{Database: models.MonitorDB}

	ServerMgr   = &confhttp.Server{}
	ServerEvent = &confhttp.Server{} // serverEvent support event http transport
	TaskMgr     = &confmq.Config{}

	fs  filesystem.FileSystemOp
	std = conflog.Std().(conflog.LevelSetter).SetLevel(conflog.InfoLevel)
)

func init() {
	// TODO config struct should be defined outside this method and impl it's Init() interface{}
	// TODO split this init too long
	config := &struct {
		Postgres      *confpostgres.Endpoint
		MonitorDB     *confpostgres.Endpoint
		MqttBroker    *confmqtt.Broker
		Redis         *confredis.Redis
		NewLogger     *conflogger.Config
		Tracer        *conftracer.Config
		Server        *confhttp.Server
		Jwt           *confjwt.Jwt
		Logger        *conflog.Log
		UploadConf    *types.UploadConfig
		EthClient     *types.ETHClientConfig
		ChainConfig   *types.ChainConfig
		WhiteList     *types.EthAddressWhiteList
		ServerEvent   *confhttp.Server
		FileSystem    *types.FileSystem
		AmazonS3      *amazonS3.AmazonS3
		LocalFS       *local.LocalFileSystem
		WasmDBConfig  *types.WasmDBConfig
		RateLimit     *confrate.RateLimit
		MetricsCenter *types.MetricsCenterConfig
		RobotNotifier *types.RobotNotifierConfig
		Risc0Config   *types.Risc0Config
		Mq            *confmq.Config
	}{
		Postgres:      db,
		MonitorDB:     monitordb,
		MqttBroker:    &confmqtt.Broker{},
		Redis:         &confredis.Redis{},
		NewLogger:     &conflogger.Config{},
		Tracer:        &conftracer.Config{},
		Server:        ServerMgr,
		Jwt:           &confjwt.Jwt{},
		Logger:        &conflog.Log{},
		UploadConf:    &types.UploadConfig{},
		EthClient:     &types.ETHClientConfig{},
		ChainConfig:   &types.ChainConfig{},
		WhiteList:     &types.EthAddressWhiteList{},
		ServerEvent:   ServerEvent,
		FileSystem:    &types.FileSystem{},
		AmazonS3:      &amazonS3.AmazonS3{},
		LocalFS:       &local.LocalFileSystem{},
		WasmDBConfig:  &types.WasmDBConfig{},
		RateLimit:     &confrate.RateLimit{},
		MetricsCenter: &types.MetricsCenterConfig{},
		RobotNotifier: &types.RobotNotifierConfig{},
		Risc0Config:   &types.Risc0Config{},
		Mq:            TaskMgr,
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

	App = confapp.New(
		confapp.WithName(name),
		confapp.WithRoot(".."),
		confapp.WithLogger(conflogger.Std()),
	)
	App.Conf(config /*, worker*/)

	if config.FileSystem.Type == enums.FILE_SYSTEM_MODE__S3 &&
		!config.AmazonS3.IsZero() {
		fs = config.AmazonS3
	} else {
		fs = config.LocalFS
	}

	if config.RobotNotifier.IsZero() {
		config.RobotNotifier = nil
	}

	confhttp.RegisterCheckerBy(config)

	proxy = &client.Client{Port: uint16(ServerEvent.Port), Timeout: 10 * time.Second}
	proxy.SetDefault()

	redisKvDB := kvdb.NewRedisDB(config.Redis)
	operatorPool := pool.NewPool(config.Postgres)

	sfIDGenerator := confid.MustNewSFIDGenerator()

	wasmApiServer, err := wasmapi.NewServer(std, config.Redis, config.Postgres, redisKvDB, config.ChainConfig,
		config.Mq, operatorPool, sfIDGenerator, config.Risc0Config)
	if err != nil {
		std.Fatal(err)
	}

	WithContext = contextx.WithContextCompose(
		types.WithMgrDBExecutorContext(config.Postgres),
		types.WithMonitorDBExecutorContext(config.MonitorDB),
		types.WithRedisEndpointContext(config.Redis),
		types.WithLoggerContext(std),
		conflog.WithLoggerContext(std),
		types.WithUploadConfigContext(config.UploadConf),
		types.WithMqttBrokerContext(config.MqttBroker),
		confid.WithSFIDGeneratorContext(sfIDGenerator),
		confjwt.WithConfContext(config.Jwt),
		confmq.WithMqContext(config.Mq),
		types.WithETHClientConfigContext(config.EthClient),
		types.WithChainConfigContext(config.ChainConfig),
		types.WithEthAddressWhiteListContext(config.WhiteList),
		types.WithFileSystemOpContext(fs),
		types.WithProxyClientContext(proxy),
		types.WithWasmDBConfigContext(config.WasmDBConfig),
		confrate.WithRateLimitKeyContext(config.RateLimit),
		kvdb.WithRedisDBKeyContext(redisKvDB),
		types.WithMetricsCenterConfigContext(config.MetricsCenter),
		types.WithRobotNotifierConfigContext(config.RobotNotifier),
		types.WithWasmApiServerContext(wasmApiServer),
		types.WithOperatorPoolContext(operatorPool),
	)
	Context = WithContext(context.Background())
}

func Server() kit.Transport {
	return ServerMgr.WithContextInjector(WithContext).WithName("srv-applet-mgr")
}

func TaskServer() kit.Transport {
	return TaskMgr.WithContextInjector(WithContext)
}

func EventServer() kit.Transport {
	return ServerEvent.WithContextInjector(WithContext).WithName("srv-event")
}

func Migrate() {
	ctx, l := conflogger.NewSpanContext(context.Background(), "global.Migrate")
	defer l.End()

	if err := migration.Migrate(db.WithContext(ctx), nil); err != nil {
		l.Error(err)
		panic(err)
	}

	if err := migration.Migrate(monitordb.WithContext(ctx), nil); err != nil {
		l.Error(err)
		panic(err)
	}
}
