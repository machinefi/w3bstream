package migration

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

type SqlMetaEnum struct {
	TName string `db:"F_table_name,size=64"`
	CName string `db:"F_column_name,size=64"`
	Value int    `db:"F_value"`
	Type  string `db:"F_type,size=255"`
	Key   string `db:"F_key,size=255"`
	Label string `db:"F_label,size=255"`
}

func (*SqlMetaEnum) TableName() string {
	return "t_sql_meta_enum"
}

func (*SqlMetaEnum) UniqueIndexes() builder.Indexes {
	return builder.Indexes{"I_enum": {"TName", "CName", "Value"}}
}

func SyncEnum(db sqlx.DBExecutor) error {
	tbl := builder.T((&SqlMetaEnum{}).TableName()).WithSchema(db.D().Schema)

	dialect := db.Dialect()

	builder.ScanDefToTable(tbl, &SqlMetaEnum{})

	tx := sqlx.NewTasks(db).With(
		func(db sqlx.DBExecutor) error {
			_, err := db.Exec(dialect.DropTable(tbl))
			return err
		},
	)

	exprs := dialect.CreateTableIsNotExists(tbl)
	for i := range exprs {
		expr := exprs[i]
		tx = tx.With(func(db sqlx.DBExecutor) error {
			_, err := db.Exec(expr)
			return err
		})
	}

	// insert values
	stmt := builder.Insert().Into(tbl)
	vals := make([]interface{}, 0)
	cols := &builder.Columns{}

	db.D().Tables.Range(func(table *builder.Table, idx int) {
		table.Columns.Range(func(col *builder.Column, idx int) {
			rv, ok := typesx.TryNew(col.ColumnType.Type)
			if !ok {
				return
			}
			enumValue, ok := rv.Interface().(enum.Enum)
			if !ok {
				return
			}
			for _, e := range enumValue.ConstValues() {
				meta := &SqlMetaEnum{
					TName: table.Name,
					CName: col.Name,
					Type:  e.TypeName(),
					Value: e.Int(),
					Key:   e.String(),
					Label: e.Label(),
				}
				fvs := builder.FieldValueFromStructByNoneZero(meta, "Value")
				cs, vs := tbl.ColumnsAndValuesByFieldValues(fvs)
				vals = append(vals, vs...)
				cols = cs
			}
		})
	})

	if len(vals) > 0 {
		stmt = stmt.Values(cols, vals...)
		tx = tx.With(func(db sqlx.DBExecutor) error {
			_, err := db.Exec(stmt)
			return err
		})
	}

	return tx.Do()
}
