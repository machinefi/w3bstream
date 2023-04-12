package schema

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
)

func Cols(names ...string) *Columns {
	cs := &Columns{}

	for _, name := range names {
		cs.Add(Col(name))
	}
	return cs
}

type Columns struct {
	lst     []*Column
	autoInc *Column
	*mapx.Map[string, *Column]
}

var _ builder.SqlExpr = (*Columns)(nil)

func (cs *Columns) IsNil() bool { return cs == nil || cs.Len() == 0 }

func (cs *Columns) Ex(ctx context.Context) *builder.Ex {
	e := builder.Expr("")
	e.Grow(cs.Len())

	cs.Range(func(c *Column, idx int) {
		if idx > 0 {
			e.WriteQueryByte(',')
		}
		e.WriteExpr(c)
	})
	return e.Ex(ctx)
}

func (cs *Columns) AutoIncrement() *Column { return cs.autoInc }

func (cs *Columns) Col(name string) *Column {
	if cs.Map != nil {
		if c, ok := cs.Load(strings.ToLower(name)); ok {
			return c
		}
	}
	return nil
}

func (cs *Columns) Cols(names ...string) (*Columns, error) {
	if len(names) == 0 {
		return cs.Clone(), nil
	}
	cols := &Columns{}
	for _, name := range names {
		c := cs.Col(name)
		if c == nil {
			return nil, errors.Errorf("unknown column: %s", name)
		}
		cols.Add(c)
	}
	return cols, nil
}

func (cs *Columns) MustCols(names ...string) *Columns {
	cols, err := cs.Cols(names...)
	must.NoError(err)
	return cols
}

func (cs *Columns) Clone() *Columns {
	cloned := &Columns{}
	cloned.Add(cs.lst...)
	return cloned
}

func (cs *Columns) Len() int {
	if cs == nil || cs.lst == nil {
		return 0
	}
	return len(cs.lst)
}

func (cs *Columns) Add(cols ...*Column) {
	if cs.Map == nil {
		cs.Map = mapx.New[string, *Column]()
	}
	for _, c := range cols {
		if c != nil {
			cs.lst = append(cs.lst, c)
			if c.Constrains != nil && c.Constrains.AutoIncrement {
				cs.autoInc = c
			}
			cs.Store(strings.ToLower(c.Name), c)
		}
	}
}

func (cs *Columns) Range(f func(c *Column, idx int)) {
	for idx, c := range cs.lst {
		f(c, idx)
	}
}

func (cs *Columns) Reset() {
	cs.lst = cs.lst[0:0]
	cs.autoInc = nil
	if cs.Map != nil {
		cs.Clear()
	}
}
