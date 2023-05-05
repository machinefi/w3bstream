package requires

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/clients/applet_mgr"
	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/http"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	confjwt "github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/conf/postgres"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/client"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/migration"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/retry"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
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

// Database executor for access database for testing
func Databases() {
	ep := &postgres.Endpoint{
		Master: base.Endpoint{
			Scheme:   "postgresql",
			Hostname: "localhost",
			Port:     5432,
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
		if _dbMgrEp, _dbMgr, err = migrate(models.DB); err != nil {
			panic(err)
		}
	}
	if _dbMonitor == nil {
		if _dbMonitorEp, _dbMonitor, err = migrate(models.MonitorDB); err != nil {
			panic(err)
		}
	}
	if _dbWasm == nil {
		if _dbWasmEp, _dbWasm, err = migrate(models.WasmDB); err != nil {
			panic(err)
		}
	}
}

func Mqtt() {
	if _broker != nil {
		return
	}
	_broker = &mqtt.Broker{
		Server: base.Endpoint{
			Scheme:   "mqtt",
			Hostname: "localhost",
			Port:     1883,
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

func Serve() (stop func()) {
	go kit.Run(apis.RootMgr, _server.WithContextInjector(_injection))
	time.Sleep(3 * time.Second)

	return func() {
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(os.Interrupt)
		time.Sleep(3 * time.Second)
	}
}

func Server() {
	if _server == nil {
		_server = &http.Server{
			Port:  8888,
			Debug: ptrx.Ptr(true),
		}
		_server.SetDefault()
	}
}

func Context() context.Context {
	return _ctx
}

var (
	_server      *http.Server
	_client      *client.Client
	_broker      *mqtt.Broker
	_dbMgr       sqlx.DBExecutor
	_dbMonitor   sqlx.DBExecutor
	_dbWasm      sqlx.DBExecutor
	_dbMgrEp     *postgres.Endpoint
	_dbMonitorEp *postgres.Endpoint
	_dbWasmEp    *postgres.Endpoint
	_injection   contextx.WithContext
	_ctx         context.Context
)

func init() {
	Databases()
	Mqtt()
	Server()
	Client()

	std := conflog.Std()
	jwt := &confjwt.Jwt{
		Issuer:  "w3bstream_test",
		ExpIn:   *base.AsDuration(time.Hour),
		SignKey: stringsx.GenRandomVisibleString(16),
	}

	_injection = contextx.WithContextCompose(
		types.WithMgrDBExecutorContext(_dbMgr),
		types.WithMonitorDBExecutorContext(_dbMonitor),
		types.WithWasmDBEndpointContext(_dbWasmEp),
		types.WithLoggerContext(std),
		types.WithMqttBrokerContext(_broker),
		conflog.WithLoggerContext(std),
		confid.WithSFIDGeneratorContext(confid.MustNewSFIDGenerator()),
		confjwt.WithConfContext(jwt),
	)

	_ctx = _injection(context.Background())
}
