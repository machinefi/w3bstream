package sqlx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

type DBExecutor interface {
	ExprExecutor
	TableResolver

	Dialect() builder.Dialect
	D() *Database
	WithSchema(string) DBExecutor

	Context() context.Context
	WithContext(ctx context.Context) DBExecutor
}

type WithDBName interface {
	WithDBName(string) driver.Connector
}

type SqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type ExprExecutor interface {
	SqlExecutor
	Exec(builder.SqlExpr) (sql.Result, error)
	Query(builder.SqlExpr) (*sql.Rows, error)
	QueryAndScan(builder.SqlExpr, interface{}) error
}

type TableResolver interface {
	T(model builder.Model) *builder.Table
}

type TxExecutor interface {
	IsTx() bool
	BeginTx(*sql.TxOptions) (DBExecutor, error)
	Begin() (DBExecutor, error)
	Commit() error
	Rollback() error
}

type Migrator interface {
	Migrate(ctx context.Context, db DBExecutor) error
}

func NewDatabase(name string) *Database {
	return &Database{Name: name, Tables: builder.Tables{}}
}

type Database struct {
	Name   string
	Schema string
	Tables builder.Tables
}

func (d Database) WithSchema(schema string) *Database {
	d.Schema = schema

	tables := builder.Tables{}
	d.Tables.Range(func(t *builder.Table, _ int) {
		tables.Add(t.WithSchema(schema))
	})
	d.Tables = tables
	return &d
}

func (d *Database) OpenDB(connector driver.Connector) *DB {
	if c, ok := connector.(WithDBName); ok {
		connector = c.WithDBName(d.Name)
	}
	dialect, ok := connector.(builder.Dialect)
	if !ok {
		panic("connect MUST be a builder.Dialect")
	}
	return &DB{
		Database:    d,
		dialect:     dialect,
		SqlExecutor: sql.OpenDB(connector),
	}
}

func (d *Database) T(model builder.Model) *builder.Table {
	if t, ok := model.(builder.TableDefinition); ok {
		return t.T()
	}
	if t, ok := model.(*builder.Table); ok {
		return t
	}
	return d.Table(model.TableName())
}

func (d *Database) Table(name string) *builder.Table { return d.Tables.Table(name) }

func (d *Database) Register(m builder.Model) *builder.Table {
	t := builder.TableFromModel(m)
	t.Schema = d.Schema
	d.AddTable(t)
	return t
}

func (d *Database) AddTable(t *builder.Table) { d.Tables.Add(t) }

type DB struct {
	ctx context.Context

	dialect builder.Dialect
	*Database
	SqlExecutor
}

func (d *DB) WithContext(ctx context.Context) DBExecutor {
	e := new(DB)
	*e = *d
	e.ctx = ctx
	return e
}

func (d *DB) Context() context.Context {
	if d.ctx == nil {
		return context.Background()
	}
	return d.ctx
}

func (d DB) WithSchema(schema string) DBExecutor {
	d.Database = d.Database.WithSchema(schema)
	return &d
}

func (d *DB) Dialect() builder.Dialect { return d.dialect }

func (d *DB) D() *Database { return d.Database }

func (d *DB) Migrate(ctx context.Context, db DBExecutor) error {
	if migrator, ok := d.dialect.(Migrator); ok {
		return migrator.Migrate(ctx, db)
	}
	return nil
}

func (d *DB) Exec(e builder.SqlExpr) (sql.Result, error) {
	ex := builder.ResolveExprContext(d.Context(), e)
	if builder.IsNilExpr(ex) {
		return nil, nil
	}
	if err := ex.Err(); err != nil {
		return nil, err
	}
	res, err := d.ExecContext(d.Context(), ex.Query(), ex.Args()...)
	if err != nil {
		if d.dialect.IsErrorConflict(err) {
			return nil, NewSqlError(sqlErrTypeConflict, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (d *DB) Query(e builder.SqlExpr) (*sql.Rows, error) {
	ex := builder.ResolveExprContext(d.Context(), e)
	if builder.IsNilExpr(ex) {
		return nil, nil
	}
	if err := ex.Err(); err != nil {
		return nil, err
	}
	return d.QueryContext(d.Context(), ex.Query(), ex.Args()...)
}

func (d *DB) QueryAndScan(e builder.SqlExpr, v interface{}) error {
	rows, err := d.Query(e)
	if err != nil {
		return err
	}
	return Scan(d.Context(), rows, v)
}

func (d *DB) IsTx() bool { _, ok := d.SqlExecutor.(*sql.Tx); return ok }

func (d *DB) Begin() (DBExecutor, error) { return d.BeginTx(nil) }

func (d *DB) BeginTx(opt *sql.TxOptions) (DBExecutor, error) {
	if d.IsTx() {
		return nil, ErrNotDB
	}
	db, err := d.SqlExecutor.(*sql.DB).BeginTx(d.Context(), opt)
	if err != nil {
		return nil, err
	}
	e := *d
	e.SqlExecutor = db
	e.ctx = d.Context()
	return &e, nil
}

func (d *DB) Commit() error {
	if !d.IsTx() {
		return ErrNotTx
	}
	if d.Context().Err() == context.Canceled {
		return context.Canceled
	}
	return d.SqlExecutor.(*sql.Tx).Commit()
}

func (d *DB) Rollback() error {
	if !d.IsTx() {
		return ErrNotTx
	}
	if d.Context().Err() == context.Canceled {
		return context.Canceled
	}
	return d.SqlExecutor.(*sql.Tx).Rollback()
}

func (d *DB) SetMaxOpenConns(n int) { d.SqlExecutor.(*sql.DB).SetMaxOpenConns(n) }

func (d *DB) SetMaxIdleConns(n int) { d.SqlExecutor.(*sql.DB).SetMaxIdleConns(n) }

func (d *DB) SetConnMaxLifetime(du time.Duration) { d.SqlExecutor.(*sql.DB).SetConnMaxLifetime(du) }
