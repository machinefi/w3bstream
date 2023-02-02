package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/driver/postgres"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/retry"
)

type Endpoint struct {
	Master   types.Endpoint
	Slave    types.Endpoint
	Database *sqlx.Database `env:"-"`
	Retry    *retry.Retry

	Extensions      []string
	PoolSize        int
	ConnMaxLifetime types.Duration

	*sqlx.DB `env:"-"`
	slave    *sqlx.DB `env:"-"`
}

var _ types.DefaultSetter = (*Endpoint)(nil)

func (e Endpoint) LivenessCheck() map[string]string {
	s := map[string]string{}

	_, err := e.DB.ExecContext(context.Background(), "SELECT 1")
	if err != nil {
		s[e.Master.Host()] = err.Error()
	} else {
		s[e.Master.Host()] = "ok"
	}

	if e.slave != nil {
		_, err := e.slave.ExecContext(context.Background(), "SELECT 1")
		if err != nil {
			s[e.Slave.Host()] = err.Error()
		} else {
			s[e.Slave.Host()] = "ok"
		}
	}

	return s
}

func (e *Endpoint) SetDefault() {
	if e.PoolSize == 0 {
		e.PoolSize = 10
	}
	if e.ConnMaxLifetime == 0 {
		e.ConnMaxLifetime = types.Duration(time.Hour)
	}
	if e.Master.IsZero() {
		e.Master.Hostname, e.Master.Port = "127.0.0.1", 5432
	}
	e.Master.Scheme = "postgres"
	if e.Database.Name == "" && len(e.Master.Base) > 0 {
		e.Database.Name = e.Master.Base
	}
	if e.Retry == nil {
		e.Retry = retry.Default
	}
}

func (e Endpoint) UseSlave() sqlx.DBExecutor {
	if e.slave == nil {
		return e.DB
	}
	return e.slave
}

// Establish a new connection in master mode with pg
func (e Endpoint) NewConnection() (sqlx.DBExecutor, error) {
	return e.conn(e.masterURL(), false)
}

func (e *Endpoint) conn(url string, readonly bool) (*sqlx.DB, error) {
	connector := &postgres.Connector{
		Host:  url,
		Extra: e.Master.Param.Encode(),
	}
	if !readonly {
		connector.Extensions = e.Extensions
	}
	db := e.Database.OpenDB(connector)
	db.SetMaxOpenConns(e.PoolSize)
	db.SetMaxIdleConns(e.PoolSize / 2)
	db.SetConnMaxLifetime(e.ConnMaxLifetime.Duration())

	_, err := db.ExecContext(context.Background(), "SELECT 1")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (e *Endpoint) Init() {
	// cover default database name
	if len(e.Master.Base) > 0 {
		e.Database.Name = e.Master.Base
	}
	// must try master
	must.NoError(e.Retry.Do(func() error {
		db, err := e.conn(e.masterURL(), false)
		if err != nil {
			return err
		}
		e.DB = db
		return nil
	}))
	// try slave if config
	if !e.Slave.IsZero() {
		must.NoError(e.Retry.Do(func() error {
			db, err := e.conn(e.slaveURL(), false)
			if err != nil {
				return err
			}
			e.slave = db
			return nil
		}))
	}
}

func (e Endpoint) masterURL() string {
	passwd := e.Master.Password
	if passwd != "" {
		passwd = ":" + passwd
	}
	return fmt.Sprintf("postgres://%s%s@%s", e.Master.Username, passwd, e.Master.Host())
}

func (e Endpoint) slaveURL() string {
	passwd := e.Master.Password
	if passwd != "" {
		passwd = ":" + passwd
	}
	return fmt.Sprintf("postgres://%s%s@%s", e.Master.Username, passwd, e.Slave.Host())
}

func (e Endpoint) Name() string { return "pgcli" }

func SwitchSlave(db sqlx.DBExecutor) sqlx.DBExecutor {
	if slave, ok := db.(CanSlave); ok {
		return slave.UseSlave()
	}
	return db
}

type CanSlave interface {
	UseSlave() sqlx.DBExecutor
}
