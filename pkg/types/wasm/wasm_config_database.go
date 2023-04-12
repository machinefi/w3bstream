package wasm

import (
	"context"
	"database/sql"
	"os"

	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/migration"
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
	return WithSQLStore(ctx, s)
}

func (s *Schema) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.DBExecutor(s.db).ExecContext(ctx, query, args...)
}

func (s *Schema) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.DBExecutor(s.db).QueryContext(ctx, query, args...)
}

func (s *Schema) Init(ctx context.Context) (err error) {
	var (
		prj  = types.MustProjectFromContext(ctx)
		ep   = types.MustWasmDBEndpointFromContext(ctx)
		_, l = conflog.FromContext(ctx).Start(ctx, "schema.Init")
		name = "w3b_" + prj.ProjectID.String()
	)

	ep.Database = sqlx.NewDatabase(name)
	err = ep.Init()
	if err != nil {
		return err
	}

	for _, t := range s.Tables {
		tbl := builder.T(t.Name)
		tbl.Desc = []string{t.Desc}
		for _, c := range t.Cols {
			col := builder.Col(c.Name)
			dt := c.Constrains
			col.ColumnType = &builder.ColumnType{
				DataType:      dt.DatabaseDatatype(),
				Length:        dt.Length,
				Decimal:       dt.Decimal,
				Default:       dt.Default,
				Null:          dt.Null,
				AutoIncrement: dt.AutoIncrement,
				Comment:       dt.Desc,
				Desc:          []string{dt.Desc},
			}
			tbl.AddCol(col)
		}
		for _, k := range t.Keys {
			key := &builder.Key{
				Name:     k.Name,
				IsUnique: k.IsUnique,
				Method:   k.Method,
			}
			tbl.AddKey(key)
		}
		ep.AddTable(tbl)
	}
	if s.Name == "" {
		s.Name = prj.Name
	}
	s.db = ep.WithSchema(s.Name)

	if err = migration.Migrate(s.db, os.Stderr); err != nil {
		l.Error(err)
		return err
	}
	if err = migration.Migrate(s.db, nil); err != nil {
		l.Error(err)
		return err
	}
	return nil
}

// func (s *Schema) Uninit(ctx context.Context) error {
// 	_, l := conflog.FromContext(ctx).Start(ctx, "schema.Uninit")
// 	db := s.DBExecutor(s.db)
//
// 	l = l.WithValues("schema", s.Name)
// 	exprs := make([]builder.SqlExpr, 0)
// 	for _, t := range s.Tables {
// 		exprs = append(exprs, t.DropIfExists()...)
// 	}
// 	exprs = append(exprs, s.DropSchema())
//
// 	for _, expr := range exprs {
// 		if _, err := db.Exec(expr); err != nil {
// 			l.Error(err)
// 		}
// 	}
// 	return nil
// }
