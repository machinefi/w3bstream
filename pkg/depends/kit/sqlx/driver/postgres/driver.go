package postgres

import (
	"bytes"
	"context"
	"database/sql/driver"
	"strconv"
	"strings"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/timer"
)

type Driver struct {
	drv pq.Driver
}

func (d *Driver) Open(dsn string) (driver.Conn, error) {
	cfg, err := pq.ParseURL(dsn)
	if err != nil {
		return nil, err
	}
	opts := ParseOption(cfg)
	if passwd, ok := opts["password"]; ok {
		opts["password"] = strings.Repeat("*", len(passwd))
	}
	conn, err := d.drv.Open(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Driver.Open")
	}
	return &LoggingConn{opts, conn}, nil
}

type LoggingConn struct {
	opts Opts
	driver.Conn
}

var _ interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.QueryerContext
} = (*LoggingConn)(nil)

func (c *LoggingConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	l := log.FromContext(ctx)
	l.Debug("=========== Beginning Transaction ===========")
	tx, err := c.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
	if err != nil {
		l.Error(errors.Wrap(err, "failed to begin transaction"))
		return nil, err
	}
	return &LoggingTx{tx: tx, l: l}, nil
}

func (c *LoggingConn) Prepare(string) (driver.Stmt, error) { panic("dont use Prepare") }

func (c *LoggingConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	cost := timer.Start()
	_ctx, l := log.Start(ctx, "Query")

	defer func() {
		qs := interpolate(query, args)

		if err != nil {
			if pgErr, ok := sqlx.UnwrapAll(err).(*pq.Error); !ok {
				l.Error(errors.Wrapf(err, "query failed: %s", qs))
			} else {
				l.Warn(errors.Wrapf(pgErr, "query failed: %s", qs))
			}
		} else {
			l.WithValues("cost", cost().String()).Debug("%s", qs)
		}

		l.End()
	}()

	rows, err = c.Conn.(driver.QueryerContext).QueryContext(_ctx, replaceValueHolder(query), args)
	return
}

func (c *LoggingConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (res driver.Result, err error) {
	cost := timer.Start()
	_ctx, l := log.Start(ctx, "Exec")

	defer func() {
		qs := interpolate(query, args)

		if err != nil {
			if pgError, ok := sqlx.UnwrapAll(err).(*pq.Error); !ok {
				l.Error(errors.Wrapf(err, "exec failed: %s", qs))
			} else if pgError.Code == "23505" {
				l.Warn(errors.Wrapf(pgError, "exec failed: %s", qs))
			} else {
				l.Error(errors.Wrapf(pgError, "exec failed: %s", qs))
			}
			return
		}

		l.WithValues("cost", cost().String()).Debug(qs.String())

		l.End()
	}()

	res, err = c.Conn.(driver.ExecerContext).ExecContext(_ctx, replaceValueHolder(query), args)
	return
}

type LoggingTx struct {
	l  log.Logger
	tx driver.Tx
}

func (tx *LoggingTx) Commit() error {
	if err := tx.tx.Commit(); err != nil {
		tx.l.Debug("failed to commit transaction: %s", err)
		return err
	}
	tx.l.Debug("=========== Committed Transaction ===========")
	return nil

}

func (tx *LoggingTx) Rollback() error {
	if err := tx.tx.Rollback(); err != nil {
		tx.l.Debug("failed to rollback transaction: %s", err)
		return err
	}
	tx.l.Debug("=========== Rollback Transaction ===========")
	return nil
}

func replaceValueHolder(query string) string {
	qc := 0
	buf := bytes.NewBuffer(nil)

	for _, c := range []byte(query) {
		switch c {
		case '?':
			buf.WriteByte('$')
			buf.WriteString(strconv.Itoa(qc + 1))
			qc++
		default:
			buf.WriteByte(c)
		}
	}
	return buf.String()
}
