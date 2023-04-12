package schema

import (
	"context"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
)

type WithTableDefinition struct{ t *Table }

func (with *WithTableDefinition) WithTable(t *Table) { with.t = t }

func (with *WithTableDefinition) T() *Table { return with.t }

type TableDefinition interface {
	WithTable(*Table)
	T() *Table
}

// T is used to contribute a table structure
func T(name string, defs ...TableDefinition) *Table {
	t := &Table{Name: name}
	for _, def := range defs {
		if c, ok := def.(*Column); ok {
			t.AddCol(c)
		}
		if k, ok := def.(*Key); ok {
			t.AddKey(k)
		}
	}
	return t
}

type Table struct {
	Name             string    `json:"name"`
	Desc             string    `json:"desc,omitempty"`
	Cols             []*Column `json:"cols"`
	Keys             []*Key    `json:"keys"`
	WithSoftDeletion bool      `json:"withSoftDeletion,omitempty"`
	WithPrimaryKey   bool      `json:"withPrimaryKey,omitempty"`

	schema string
	cols   Columns
	keys   Keys
}

var _ builder.SqlExpr = (*Table)(nil)

func (t *Table) IsNil() bool { return t == nil || t.Name == "" }

func (t *Table) Ex(ctx context.Context) *builder.Ex {
	if t.schema != "" {
		return builder.Expr(t.schema + "." + t.Name).Ex(ctx)
	}
	return builder.Expr(t.Name).Ex(ctx)
}

func (t *Table) WithSchema(schema string) { t.schema = schema }

func (t *Table) Init() error {
	t.cols.Reset()
	t.keys.Reset()

	if t.WithPrimaryKey {
		c := &Column{
			Name: "f_id",
			Constrains: &ColumnType{
				Datatype:      DATATYPE__UINT64,
				AutoIncrement: true,
			},
			WithTableDefinition: WithTableDefinition{t: t},
		}
		t.cols.Add(c)
	}
	for _, c := range t.Cols {
		dfv := c.Constrains.Default
		if dfv != nil && len(*dfv) > 0 && (*dfv)[0] != '\'' {
			c.Constrains.Default = ptrx.Ptr("'" + *c.Constrains.Default + "'")
		}
	}

	t.cols.Add(t.Cols...)
	t.cols.Add(
		&Column{
			Name: "f_created_at",
			Constrains: &ColumnType{
				Datatype: DATATYPE__TIMESTAMP,
				Default:  ptrx.Ptr("'0'"),
			},
			WithTableDefinition: WithTableDefinition{t: t},
		},
		&Column{
			Name: "f_updated_at",
			Constrains: &ColumnType{
				Datatype: DATATYPE__TIMESTAMP,
				Default:  ptrx.Ptr("'0'"),
			},
			WithTableDefinition: WithTableDefinition{t: t},
		},
	)
	if t.WithSoftDeletion {
		t.cols.Add(
			&Column{
				Name: "f_deleted_at",
				Constrains: &ColumnType{
					Datatype: DATATYPE__TIMESTAMP,
					Default:  ptrx.Ptr("'0'"),
				},
				WithTableDefinition: WithTableDefinition{t: t},
			},
		)
	}

	t.keys.Add(t.Keys...)

	if t.WithSoftDeletion {
		t.keys.Range(func(k *Key, idx int) {
			if k.IsUnique {
				k.ColumnNames = append(k.ColumnNames, "f_deleted_at")
			}
		})
	}
	t.cols.Range(func(c *Column, idx int) {
		c.WithTable(t)
	})
	t.keys.Range(func(k *Key, idx int) {
		k.WithTable(t)
	})
	return nil
}

func (t *Table) AddKey(k *Key) {
	if k != nil {
		t.keys.Add(k)
		t.Keys = append(t.Keys, k)
	}
}

func (t *Table) AddCol(c *Column) {
	if c != nil {
		t.cols.Add(c)
		t.Cols = append(t.Cols, c)
	}
}

func (t *Table) CreateIfNotExists() []builder.SqlExpr {
	e := builder.Expr("CREATE TABLE IF NOT EXISTS ")
	e.WriteExpr(t)
	e.WriteQueryByte(' ')
	e.WriteGroup(func(ex *builder.Ex) {
		if t.cols.IsNil() {
			return
		}

		t.cols.Range(func(c *Column, idx int) {
			if idx > 0 {
				ex.WriteQueryByte(',')
			}
			e.WriteQueryByte('\n')
			e.WriteQueryByte('\t')
			e.WriteExpr(c)
			e.WriteQueryByte(' ')
			e.WriteExpr(c.Constrains.Ex(context.Background()))
		})

		t.keys.Range(func(k *Key, idx int) {
			if k.IsPrimary() {
				e.WriteQueryByte(',')
				e.WriteQueryByte('\n')
				e.WriteQueryByte('\t')
				e.WriteQuery("PRIMARY KEY ")
				e.WriteGroup(func(e *builder.Ex) {
					e.WriteExpr(Cols(k.IndexDef.ColumnNames...))
				})
			}
		})

		e.WriteQueryByte('\n')
	})
	e.WriteEnd()

	es := []builder.SqlExpr{e}

	t.keys.Range(func(k *Key, idx int) {
		if !k.IsPrimary() {
			es = append(es, t.AddIndex(k))
		}
	})
	return es
}

func (t *Table) DropIfExists() []builder.SqlExpr {
	exprs := make([]builder.SqlExpr, 0)
	for _, k := range t.Keys {
		exprs = append(exprs, t.DropIndex(k))
	}

	e := builder.Expr("DROP TABLE IF EXISTS ")
	e.WriteExpr(t)
	e.WriteEnd()

	exprs = append(exprs, e)
	return exprs
}

func (t *Table) AddIndex(k *Key) builder.SqlExpr {
	if k.IsPrimary() {
		e := builder.Expr("ALTER TABLE ")
		e.WriteExpr(k.t)
		e.WriteQuery(" ADD PRIMARY KEY ")
		e.WriteGroup(func(e *builder.Ex) {
			e.WriteExpr(Cols(k.IndexDef.ColumnNames...))
		})
		e.WriteEnd()
		return e
	}

	e := builder.Expr("CREATE ")
	if k.IsUnique {
		e.WriteQuery("UNIQUE ")
	}
	e.WriteQuery("INDEX IF NOT EXISTS ")

	e.WriteQuery(k.t.Name)
	e.WriteQuery("_")
	e.WriteQuery(k.Name)
	e.WriteQuery(" ON ")
	e.WriteExpr(k.t)

	if mtd := strings.ToUpper(k.Method); mtd != "" {
		e.WriteQuery(" USING " + mtd)
	}
	e.WriteQueryByte(' ')
	e.WriteGroup(func(e *builder.Ex) {
		e.WriteExpr(Cols(k.IndexDef.ColumnNames...))
	})
	e.WriteEnd()
	return e
}

func (t *Table) DropIndex(k *Key) builder.SqlExpr {
	if k.IsPrimary() {
		e := builder.Expr("ALTER TABLE ")
		e.WriteExpr(t)
		e.WriteQuery(" DROP CONSTRAINT ")
		e.WriteExpr(t)
		e.WriteQuery("_pkey")
		e.WriteEnd()
		return e
	}
	e := builder.Expr("DROP ")

	e.WriteQuery("INDEX IF EXISTS ")
	e.WriteExpr(t)
	e.WriteQueryByte('_')
	e.WriteQuery(k.Name)
	e.WriteEnd()

	return e
}
