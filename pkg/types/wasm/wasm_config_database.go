package wasm

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/schema"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Schema struct {
	schema.Schema
	db sqlx.DBExecutor
}

// TODO impl Schema.Init

func (s *Schema) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__PROJECT_SCHEMA
}

func (s *Schema) WithContext(ctx context.Context) context.Context {
	db := types.MustWasmDBExecutorFromContext(ctx)
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
