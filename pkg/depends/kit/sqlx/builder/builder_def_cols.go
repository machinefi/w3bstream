package builder

import (
	"context"
	"strings"

	"github.com/pkg/errors"
)

type Columns struct {
	lst     []*Column
	autoInc *Column
}

func Cols(names ...string) *Columns {
	cols := &Columns{}
	for _, name := range names {
		cols.Add(Col(name))
	}
	return cols
}

func (c *Columns) IsNil() bool { return c == nil || c.Len() == 0 }

func (c *Columns) Ex(ctx context.Context) *Ex {
	e := Expr("")
	e.Grow(c.Len())

	c.Range(func(col *Column, idx int) {
		if idx > 0 {
			e.WriteQueryByte(',')
		}
		e.WriteExpr(col)
	})
	return e.Ex(ctx)
}

func (c *Columns) AutoIncrement() *Column { return c.autoInc }

func (c *Columns) Len() int {
	if c == nil || c.lst == nil {
		return 0
	}
	return len(c.lst)
}

func (c *Columns) Clone() *Columns {
	clone := &Columns{
		lst: make([]*Column, len(c.lst)),
	}
	copy(clone.lst, c.lst)
	return clone
}

func (c *Columns) Col(name string) *Column {
	name = strings.ToLower(name)
	for i := range c.lst {
		if c.lst[i].Name == name {
			return c.lst[i]
		}
	}
	return nil
}

func (c *Columns) ColByFieldName(name string) *Column {
	for i := range c.lst {
		if c.lst[i].FieldName == name {
			return c.lst[i]
		}
	}
	return nil
}

func (c *Columns) Cols(names ...string) (*Columns, error) {
	if len(names) == 0 {
		return c.Clone(), nil
	}
	cols := &Columns{}
	for _, name := range names {
		col := c.Col(name)
		if col == nil {
			return nil, errors.Errorf("unknown struct column %s", name)
		}
		cols.Add(col)
	}
	return cols, nil
}

func (c *Columns) ColsByFieldNames(names ...string) (*Columns, error) {
	if len(names) == 0 {
		return c.Clone(), nil
	}
	cols := &Columns{lst: make([]*Column, 0, len(names))}
	for _, name := range names {
		col := c.ColByFieldName(name)
		if col == nil {
			return nil, errors.Errorf("unknonw struct field %s", name)
		}
		cols.lst = append(cols.lst, col)
	}
	return cols, nil
}

func (c *Columns) MustCols(names ...string) *Columns {
	cols, err := c.Cols(names...)
	if err != nil {
		panic(err)
	}
	return cols
}

func (c *Columns) MustColsByFieldNames(names ...string) *Columns {
	cols, err := c.ColsByFieldNames(names...)
	if err != nil {
		panic(err)
	}
	return cols
}

func (c *Columns) ColNames() []string {
	names := make([]string, 0, c.Len())
	c.Range(func(col *Column, idx int) {
		if col.Name != "" {
			names = append(names, col.Name)
		}
	})
	return names
}

func (c *Columns) FieldNames() []string {
	names := make([]string, 0, c.Len())
	c.Range(func(col *Column, idx int) {
		if col.FieldName != "" {
			names = append(names, col.FieldName)
		}
	})
	return names
}

func (c *Columns) Add(cols ...*Column) {
	for i := range cols {
		if cols[i] == nil {
			continue
		}
		col := cols[i]
		if col.ColumnType != nil && col.ColumnType.AutoIncrement {
			if c.autoInc != nil {
				panic("auto increment field can only have one")
			}
			c.autoInc = col
		}
		c.lst = append(c.lst, col)
	}
}

func (c *Columns) Range(f func(*Column, int)) {
	for i := range c.lst {
		f(c.lst[i], i)
	}
}

func (c *Columns) List() []*Column {
	if c == nil {
		return nil
	}
	return c.lst
}
