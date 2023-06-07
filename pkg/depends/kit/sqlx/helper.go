package sqlx

import (
	"context"
	"database/sql"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/scanner"
)

func InsertToDB(db DBExecutor, m builder.Model, zeros []string, additions ...builder.Addition) builder.SqlExpr {
	t := db.T(m)
	fvs := FieldValuesFromModel(t, m, zeros...)
	cs, vs := t.ColumnsAndValuesByFieldValues(fvs)
	return builder.Insert().Into(t, additions...).Values(cs, vs...)
}

func AsAssignments(db DBExecutor, m builder.Model, zeros ...string) builder.Assignments {
	fvs := FieldValuesFromModel(db.T(m), m, zeros...)
	return db.T(m).AssignmentsByFieldValues(fvs)

}

func FieldValuesFromModel(table *builder.Table, model builder.Model, zeros ...string) builder.FieldValues {
	fvs := builder.FieldValueFromStructByNoneZero(model, zeros...)
	if autoIncrementCol := table.AutoIncrement(); autoIncrementCol != nil {
		delete(fvs, autoIncrementCol.FieldName)
	}
	return fvs
}

type ScanIterator = scanner.ScanIterator

func Scan(ctx context.Context, rows *sql.Rows, v interface{}) error {
	if err := scanner.Scan(ctx, rows, v); err != nil {
		if err == scanner.RecordNotFound {
			return NewSqlError(SqlErrTypeNotFound, "record is not found")
		}
		return err
	}
	return nil
}
