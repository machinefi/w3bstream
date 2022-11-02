package global

import (
	"context"
	"log"
	"os"
	"path/filepath"

	uconfig "go.uber.org/config"

	"github.com/machinefi/Bumblebee/base/consts"
	confapp "github.com/machinefi/Bumblebee/conf/app"
	confhttp "github.com/machinefi/Bumblebee/conf/http"
	confid "github.com/machinefi/Bumblebee/conf/id"
	confjwt "github.com/machinefi/Bumblebee/conf/jwt"
	conflog "github.com/machinefi/Bumblebee/conf/log"
	confmqtt "github.com/machinefi/Bumblebee/conf/mqtt"
	confpostgres "github.com/machinefi/Bumblebee/conf/postgres"
	"github.com/machinefi/Bumblebee/kit/kit"
	"github.com/machinefi/Bumblebee/kit/mq"
	"github.com/machinefi/Bumblebee/kit/mq/mem_mq"
	"github.com/machinefi/Bumblebee/kit/sqlx/migration"
	"github.com/machinefi/Bumblebee/x/contextx"

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

type Config struct {
	ethClient types.ETHClientConfig `yaml:"ethClient"`
}

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
	// TODO: read config file path from parameter
	configFile := filepath.Join(App.Root(), "config.yml")
	switch _, err := os.Stat(configFile); {
	case os.IsNotExist(err):
	case err == nil:
		// var cfg Config
		yaml, err := uconfig.NewYAML(uconfig.File(configFile))
		if err != nil {
			log.Fatalln("failed to load yaml file", err)
		}
		var config Config
		if err := yaml.Get(uconfig.Root).Populate(&config); err != nil {
			log.Fatalln("failed to populate config", err)
		}
		if config.ethClient.ChainEndpoint != "" {
			if _, ok := os.LookupEnv("SRV_APPLET_MGR__ETHCLIENTCONFIG__ChainEndpoint"); !ok {
				if err := os.Setenv("SRV_APPLET_MGR__ETHCLIENTCONFIG__ChainEndpoint", config.ethClient.ChainEndpoint); err != nil {
					log.Fatalln("failed to set env", err)
				}
			}
		}
		if config.ethClient.PrivateKey != "" {
			if _, ok := os.LookupEnv("SRV_APPLET_MGR__ETHCLIENTCONFIG__PrivateKey"); !ok {
				if err := os.Setenv("SRV_APPLET_MGR__ETHCLIENTCONFIG__PrivateKey", config.ethClient.PrivateKey); err != nil {
					log.Fatalln("failed to set env", err)
				}
			}
		}
	default:
		log.Fatalln("failed to load", configFile, err)
	}

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
