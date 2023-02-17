package wasm

import (
	"context"
	"database/sql"
	"fmt"

	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/schema"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Schema struct {
	schema.Schema
	db sqlx.DBExecutor
}

func (s *Schema) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__PROJECT_SCHEMA
}

func (s *Schema) WithContext(ctx context.Context) context.Context {
	db, err := types.MustWasmPgEndpointFromContext(ctx).NewConnection()
	if err != nil {
		panic(err)
	}

	if s.Name == "" {
		prj := types.MustProjectFromContext(ctx)
		s.WithName(prj.Name)
	}

	// limit the scope of sql to the schema
	if _, err := db.ExecContext(ctx, fmt.Sprintf("SET search_path TO %s", s.Name)); err != nil {
		panic(err)
	}
	s.db = db
	return WithSQLStore(ctx, s)
}

func (s *Schema) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *Schema) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *Schema) Init(ctx context.Context) error {
	_, l := conflog.FromContext(ctx).Start(ctx)
	defer l.End()

	err := s.Schema.Init()
	if err != nil {
		return err
	}
	db := types.MustWasmDBExecutorFromContext(ctx)

	expr := s.CreateSchema()
	l.Info(builder.ResolveExpr(expr).Query())
	_, err = db.Exec(expr)
	if err != nil {
		return err
	}

	db = db.WithSchema(s.Name)
	for _, t := range s.Tables {
		es := t.CreateIfNotExists()
		for _, e := range es {
			if e.IsNil() {
				continue
			}
			l.Info(builder.ResolveExpr(e).Query())
			if _, err = db.Exec(e); err != nil {
				return err
			}
		}
	}
	return nil
}
