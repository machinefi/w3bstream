package schema

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
)

func NewSchema(name string) *Schema {
	return &Schema{
		Name: name,
		Map:  mapx.New[string, *Table](),
	}
}

type Schema struct {
	Name   string   `json:"name,omitempty"`
	Tables []*Table `json:"tables,omitempty"`

	*mapx.Map[string, *Table] `json:"-"`
}

func (s *Schema) WithName(name string) { s.Name = "wasm_project__" + name }

func (s *Schema) DBExecutor(d sqlx.DBExecutor) sqlx.DBExecutor {
	return d.WithSchema(s.Name)
}

func (s *Schema) AddTable(t *Table) {
	if s.Map == nil {
		s.Map = mapx.New[string, *Table]()
	}
	if !s.StoreNX(t.Name, t) {
		panic(errors.Errorf("duplicated table: %s", t.Name))
	}
}

func (s *Schema) T(name string) *Table {
	if s.Map == nil {
		return nil
	}
	t, _ := s.Load(name)
	return t
}

func (s *Schema) Init() error {
	for _, t := range s.Tables {
		if err := t.Init(); err != nil {
			return err
		}
		s.AddTable(t)
		t.WithSchema(s.Name)
	}
	return nil
}

func FromConfig(data []byte) (*Schema, error) {
	s := &Schema{}
	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Schema) CreateSchema() builder.SqlExpr {
	e := builder.Expr("CREATE SCHEMA IF NOT EXISTS ")
	e.WriteQuery(s.Name)
	e.WriteEnd()
	return e
}

func (Schema) DataType(drv string) string { return "text" }

func (s Schema) Value() (driver.Value, error) {
	return datatypes.JSONValue(s)
}

func (s *Schema) Scan(src interface{}) error {
	return datatypes.JSONScan(src, s)
}
