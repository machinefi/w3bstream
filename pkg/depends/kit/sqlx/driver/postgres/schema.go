package postgres

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/slice"
)

var regexpUsing = regexp.MustCompile(`USING ([^ ]+)`)

func databaseFromSchema(db sqlx.DBExecutor) (*sqlx.Database, error) {
	d := db.D()

	var (
		tableNames  = slice.ToAnySlice(d.Tables.TableNames()...)
		tableSchema = SchemaDB.T(&ColumnSchema{}).WithSchema("information_schema")
		columns     = make([]ColumnSchema, 0)
	)

	d = sqlx.NewDatabase(d.Name).WithSchema(d.Schema)

	schema := "public"
	if d.Schema != "" {
		schema = d.Schema
	}

	stmt := builder.Select(tableSchema.Columns.Clone()).
		From(
			tableSchema,
			builder.Where(
				builder.And(
					tableSchema.ColByFieldName("TABLE_SCHEMA").Eq(schema),
					tableSchema.ColByFieldName("TABLE_NAME").In(tableNames...),
				),
			),
		)

	err := db.QueryAndScan(stmt, &columns)
	if err != nil {
		return nil, err
	}

	for i := range columns {
		cs := columns[i]

		tbl := d.Table(cs.TABLE_NAME)
		if tbl == nil {
			tbl = builder.T(cs.TABLE_NAME)
			d.AddTable(tbl)
		}

		tbl.AddCol(colFromSchema(&cs))
	}

	if tableSchema.Columns.Len() != 0 {
		v := SchemaDB.T(&IndexSchema{})
		indexes := make([]IndexSchema, 0)

		err = db.QueryAndScan(
			builder.Select(v.Columns.Clone()).
				From(
					v,
					builder.Where(
						builder.And(
							v.ColByFieldName("TABLE_SCHEMA").Eq(schema),
							v.ColByFieldName("TABLE_NAME").In(tableNames...),
						),
					),
				),
			&indexes,
		)

		if err != nil {
			return nil, err
		}

		for _, index := range indexes {
			table := d.Table(index.TABLE_NAME)
			key := &builder.Key{
				Name:     strings.ToLower(index.INDEX_NAME[len(table.Name)+1:]),
				Method:   strings.ToUpper(regexpUsing.FindString(index.INDEX_DEF)[6:]),
				IsUnique: strings.Contains(index.INDEX_DEF, "UNIQUE"),
				Def: builder.IndexDef{
					Expr: strings.TrimSpace(regexpUsing.Split(index.INDEX_DEF, 2)[1]),
				},
			}
			table.AddKey(key)
		}
	}

	return d, nil
}

func colFromSchema(columnSchema *ColumnSchema) *builder.Column {
	col := builder.Col(columnSchema.COLUMN_NAME)

	defaultValue := columnSchema.COLUMN_DEFAULT

	if defaultValue != "" {
		col.AutoIncrement = strings.HasSuffix(columnSchema.COLUMN_DEFAULT, "_seq'::regclass)")

		if !col.AutoIncrement {
			if !strings.Contains(defaultValue, "'::") && '0' <= defaultValue[0] && defaultValue[0] <= '9' {
				defaultValue = fmt.Sprintf("'%s'::integer", defaultValue)
			}
			col.Default = &defaultValue
		}
	}

	dataType := columnSchema.DATA_TYPE

	if col.AutoIncrement {
		if strings.HasPrefix(dataType, "big") {
			dataType = "bigserial"
		} else {
			dataType = "serial"
		}
	}

	col.DataType = dataType

	// numeric type
	if columnSchema.NUMERIC_PRECISION > 0 {
		col.Length = columnSchema.NUMERIC_PRECISION
		col.Decimal = columnSchema.NUMERIC_SCALE
	} else {
		col.Length = columnSchema.CHARACTER_MAXIMUM_LENGTH
	}

	if columnSchema.IS_NULLABLE == "YES" {
		col.Null = true
	}

	return col
}

type ColumnSchema struct {
	TABLE_SCHEMA             string `db:"table_schema"`
	TABLE_NAME               string `db:"table_name"`
	COLUMN_NAME              string `db:"column_name"`
	DATA_TYPE                string `db:"data_type"`
	IS_NULLABLE              string `db:"is_nullable"`
	COLUMN_DEFAULT           string `db:"column_default"`
	CHARACTER_MAXIMUM_LENGTH uint64 `db:"character_maximum_length"`
	NUMERIC_PRECISION        uint64 `db:"numeric_precision"`
	NUMERIC_SCALE            uint64 `db:"numeric_scale"`
}

func (ColumnSchema) TableName() string { return "columns" }

type IndexSchema struct {
	TABLE_SCHEMA string `db:"schemaname"`
	TABLE_NAME   string `db:"tablename"`
	INDEX_NAME   string `db:"indexname"`
	INDEX_DEF    string `db:"indexdef"`
}

func (IndexSchema) TableName() string { return "pg_indexes" }

var SchemaDB = sqlx.NewDatabase("INFORMATION_SCHEMA")

func init() {
	SchemaDB.Register(&ColumnSchema{})
	SchemaDB.Register(&IndexSchema{})
}
