package requires

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/clients/applet_mgr"
	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/filesystem/local"
	"github.com/machinefi/w3bstream/pkg/depends/conf/http"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	confjwt "github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/conf/postgres"
	"github.com/machinefi/w3bstream/pkg/depends/conf/redis"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/client"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq/mem_mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/migration"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/retry"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

// Client for request APIs
func Client(transports ...client.HttpTransport) *applet_mgr.Client {
	if _client == nil {
		_client = &client.Client{
			Protocol: "http",
			Host:     "localhost",
			Port:     uint16(_server.Port),
			Timeout:  time.Hour,
		}
		_client.SetDefault()
	}

	_client.Transports = append(_client.Transports, transports...)
	return applet_mgr.NewClient(_client)
}

// AuthClient client with jwt token
func AuthClient(transports ...client.HttpTransport) *applet_mgr.Client {
	return Client(NewAuthPatchRT())
}

// ClientEvent for post wasm event APIs through http transport
func ClientEvent(transports ...client.HttpTransport) *applet_mgr.Client {
	if _clientEvent == nil {
		_clientEvent = &client.Client{
			Protocol: "http",
			Host:     "localhost",
			Port:     uint16(_serverEvent.Port),
			Timeout:  time.Hour,
		}
		_clientEvent.SetDefault()
	}
	_clientEvent.Transports = append(_clientEvent.Transports, transports...)
	return applet_mgr.NewClient(_clientEvent)
}

// AuthClientEvent client with jwt token
func AuthClientEvent(transports ...client.HttpTransport) *applet_mgr.Client {
	return ClientEvent(NewAuthPatchRT())
}

// Database executor for access database for testing
func Databases() {
	ep := &postgres.Endpoint{
		Master: base.Endpoint{
			Scheme:   "postgresql",
			Hostname: "localhost",
			Port:     15432,
			Base:     "w3bstream",
			Username: "root",
			Password: "test_passwd",
			Param:    url.Values{"sslmode": []string{"disable"}},
		},
		Retry: &retry.Retry{
			Repeats:  3,
			Interval: *base.AsDuration(10 * time.Second),
		},
	}

	migrate := func(d *sqlx.Database) (*postgres.Endpoint, sqlx.DBExecutor, error) {
		ep := *ep
		ep.Database = d
		if err := ep.Init(); err != nil {
			return nil, nil, err
		}
		if err := migration.Migrate(ep.WithContext(context.Background()), nil); err != nil {
			return nil, nil, err
		}
		return &ep, &ep, nil
	}

	var err error
	if _dbMgr == nil {
		if _, _dbMgr, err = migrate(models.DB); err != nil {
			panic(err)
		}
	}
	if _dbMonitor == nil {
		if _, _dbMonitor, err = migrate(models.MonitorDB); err != nil {
			panic(err)
		}
	}
	_dbWasmEp = &ep.Master
}

func Mqtt() {
	if _broker != nil {
		return
	}
	_broker = &mqtt.Broker{
		Server: base.Endpoint{
			Scheme:   "mqtt",
			Hostname: "localhost",
			Port:     11883,
		},
		Retry: retry.Retry{
			Repeats:  3,
			Interval: *base.AsDuration(10 * time.Second),
		},
	}
	_broker.SetDefault()
	if err := _broker.Init(); err != nil {
		panic(err)
	}
}

var (
	grp = &sync.WaitGroup{}
	run = &sync.Once{}
)

func Serve() (stop func()) {
	grp.Add(1)

	run.Do(func() {
		go func() {
			go kit.Run(apis.RootMgr, _server.WithContextInjector(_injection))
			go kit.Run(apis.RootEvent, _serverEvent.WithContextInjector(_injection))

			time.Sleep(20 * time.Second)

			grp.Wait()
			_server.Shutdown()
		}()
	})
	time.Sleep(3 * time.Second)

	return func() {
		grp.Done()
	}
}

func Server() {
	if _server == nil {
		_server = &http.Server{
			Port:  18888,
			Debug: ptrx.Ptr(true),
		}
		_server.SetDefault()
	}
}

func ServerEvent() {
	if _serverEvent == nil {
		_serverEvent = &http.Server{
			Port:  18889,
			Debug: ptrx.Ptr(true),
		}
		_serverEvent.SetDefault()
	}
}

func Context() context.Context {
	return _ctx
}

var (
	_server      *http.Server
	_serverEvent *http.Server
	_client      *client.Client
	_clientEvent *client.Client
	_broker      *mqtt.Broker
	_dbMgr       sqlx.DBExecutor
	_dbMonitor   sqlx.DBExecutor
	_dbWasmEp    *base.Endpoint
	_injection   contextx.WithContext
	_ctx         context.Context
)

func init() {
	Databases()
	Mqtt()
	Server()
	ServerEvent()
	Client()
	ClientEvent()

	_jwt := &confjwt.Jwt{
		Issuer:  "w3bstream_test",
		SignKey: "xxxx",
	}
	_uploadConfig := &types.UploadConfig{}
	_fsop := &local.LocalFileSystem{}
	_redis := &redis.Redis{Port: 16379}

	for _, c := range []interface{}{_jwt, _uploadConfig, _fsop, _redis} {
		if canSetDefault, ok := c.(base.DefaultSetter); ok {
			canSetDefault.SetDefault()
		}
		switch v := c.(type) {
		case base.Initializer:
			v.Init()
		case base.ValidatedInitializer:
			if err := v.Init(); err != nil {
				panic(fmt.Sprintf("%v init failed", reflect.TypeOf(v)))
			}
		}
	}

	_tasks := mem_mq.New(0)
	_workers := mq.NewTaskWorker(_tasks, mq.WithWorkerCount(3), mq.WithChannel("apis_tests"))

	_injection = contextx.WithContextCompose(
		types.WithMgrDBExecutorContext(_dbMgr),
		types.WithMonitorDBExecutorContext(_dbMonitor),
		types.WithWasmDBEndpointContext(_dbWasmEp),
		types.WithLoggerContext(conflog.Std()),
		types.WithMqttBrokerContext(_broker),
		conflog.WithLoggerContext(conflog.Std()),
		confid.WithSFIDGeneratorContext(confid.MustNewSFIDGenerator()),
		confjwt.WithConfContext(_jwt),
		types.WithUploadConfigContext(_uploadConfig),
		types.WithFileSystemOpContext(_fsop),
		types.WithRedisEndpointContext(_redis),
		types.WithTaskWorkerContext(_workers),
		types.WithTaskBoardContext(mq.NewTaskBoard(_tasks)),
		types.WithETHClientConfigContext(nil), // can be nil
	)

	_ctx = _injection(context.Background())
}
