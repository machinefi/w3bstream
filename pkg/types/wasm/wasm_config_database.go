package wasm

import (
	"context"
	"os"

	"github.com/machinefi/w3bstream/pkg/depends/conf/postgres"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/migration"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
)

func NewDatabase(name string) *Database {
	return &Database{Name: name}
}

type Database struct {
	Name    string              `json:"-"`                            // Name: database name, this should be assigned by host
	Dialect enums.WasmDBDialect `json:"dialect,omitempty,default=''"` // Dialect database dialect
	Schemas []*Schema           `json:"schemas,omitempty"`            // Schemas
	schemas map[string]*Schema  // schemas reference of Schemas; key: schema name

	ep *postgres.Endpoint // database endpoint
}

type Schema struct {
	Name   string   `json:"schema,omitempty,default='public'"` // Name: schema name, use postgres driver, default schema is `public`
	Tables []*Table `json:"tables,omitempty"`                  // Tables: tables define

	d sqlx.DBExecutor // database executor with schema
}

type Table struct {
	Name string    `json:"name"`           // Name table name
	Desc string    `json:"desc,omitempty"` // Desc table description
	Cols []*Column `json:"cols"`           // Cols table column define
	Keys []*Key    `json:"keys"`           // Keys table index or primary define
}

func (t *Table) Build() *builder.Table {
	tbl := builder.T(t.Name)
	tbl.Desc = []string{t.Desc}
	for _, c := range t.Cols {
		tbl.AddCol(c.Build())
	}
	for _, k := range t.Keys {
		tbl.AddKey(k.Build())
	}
	return tbl
}

type Column struct {
	Name       string      `json:"name"`       // Name column name
	Constrains *Constrains `json:"constrains"` // Constrains column constrains
}

func (c *Column) Build() *builder.Column {
	col := builder.Col(c.Name)
	dt := c.Constrains
	col.ColumnType = &builder.ColumnType{
		DataType:      dt.Datatype.String(),
		Length:        dt.Length,
		Decimal:       dt.Decimal,
		Default:       dt.Default,
		Null:          dt.Null,
		AutoIncrement: dt.AutoIncrement,
		Comment:       dt.Desc,
		Desc:          []string{dt.Desc},
	}
	return col
}

type Constrains struct {
	Datatype      enums.WasmDBDatatype `json:"datatype"`
	Length        uint64               `json:"length,omitempty"`
	Decimal       uint64               `json:"decimal,omitempty"`
	Default       *string              `json:"default,omitempty"`
	Null          bool                 `json:"null,omitempty"`
	AutoIncrement bool                 `json:"autoincrement,omitempty"`
	Desc          string               `json:"desc,omitempty"`
}

type Key struct {
	Name        string   `json:"name,omitempty"`
	Method      string   `json:"method,omitempty"`
	IsUnique    bool     `json:"isUnique"`
	ColumnNames []string `json:"columnNames"`
	Expr        string   `json:"expr,omitempty"`
}

func (k *Key) Build() *builder.Key {
	return &builder.Key{
		Name:     k.Name,
		IsUnique: k.IsUnique,
		Method:   k.Method,
		Def: builder.IndexDef{
			ColNames: k.ColumnNames,
			Expr:     k.Expr,
		},
	}
}
func (d *Database) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__PROJECT_DATABASE
}

func (d *Database) WithContext(ctx context.Context) context.Context {
	return WithSQLStore(ctx, d)
}

func (d *Database) WithSchema(name string) sqlx.DBExecutor {
	if name == "" {
		name = "public"
	}
	if s, ok := d.schemas[name]; ok {
		return s.d
	}
	return d.ep.WithSchema("public")
}

func (d *Database) WithDefaultSchema() sqlx.DBExecutor {
	return d.WithSchema("public")
}

func (d *Database) Init(ctx context.Context) (err error) {
	var prj = types.MustProjectFromContext(ctx)

	// init database endpoint
	d.Name = "w3b_" + prj.ProjectID.String()
	d.ep = types.MustWasmDBEndpointFromContext(ctx)
	d.ep.Database = sqlx.NewDatabase(d.Name)
	if d.schemas == nil {
		d.schemas = make(map[string]*Schema)
	}

	// combine schema tables
	for _, s := range d.Schemas {
		if s.Name == "" {
			s.Name = "public" // pg default
		}
		if _, ok := d.schemas[s.Name]; !ok {
			d.schemas[s.Name] = &Schema{Name: s.Name}
		}

		d.schemas[s.Name].Tables = append(d.schemas[s.Name].Tables, s.Tables...)
	}

	if err = d.ep.Init(); err != nil {
		return err
	}

	// init each schema
	for _, s := range d.schemas {
		ep := d.ep
		for _, t := range s.Tables {
			ep.AddTable(t.Build())
		}
		s.d = ep.WithSchema(s.Name)
		if err = migration.Migrate(s.d, os.Stderr); err != nil {
			return err
		}
		if err = migration.Migrate(s.d, nil); err != nil {
			return err
		}
	}
	return nil
}
