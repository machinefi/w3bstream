package builder

import (
	"context"
	"reflect"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

type Column struct {
	Name      string
	FieldName string
	Table     *Table
	exactly   bool

	*ColumnType
}

func Col(name string) *Column {
	return &Column{Name: strings.ToLower(name), ColumnType: &ColumnType{}}
}

func (c *Column) IsNil() bool { return c == nil }

func (c *Column) Ex(ctx context.Context) *Ex {
	toggles := TogglesFromContext(ctx)
	if c.Table != nil && (c.exactly || toggles.Is(ToggleMultiTable)) {
		if toggles.Is(ToggleNeedAutoAlias) {
			return Expr("?.? AS ?", c.Table, Expr(c.Name), Expr(c.Name)).Ex(ctx)
		}
		return Expr("?.?", c.Table, Expr(c.Name)).Ex(ctx)
	}
	return ExactlyExpr(c.Name).Ex(ctx)
}

func (c *Column) Expr(query string, args ...interface{}) *Ex {
	n, e := len(args), Expr("")
	e.Grow(n)

	idx := 0

	for _, b := range []byte(query) {
		switch b {
		case '#':
			e.WriteExpr(c)
		case '?':
			e.WriteQueryByte(b)
			if n > idx {
				e.AppendArgs(args[idx])
				idx++
			}
		default:
			e.WriteQueryByte(b)
		}
	}
	return e
}

func (c Column) Of(t *Table) *Column {
	return &Column{
		Name:       c.Name,
		FieldName:  c.FieldName,
		Table:      t,
		exactly:    true,
		ColumnType: c.ColumnType,
	}
}

func (c Column) Type(v interface{}, tag string) *Column {
	c.ColumnType = AnalyzeColumnType(typesx.FromReflectType(reflect.TypeOf(v)), tag)
	return &c
}

func (c Column) Field(name string) *Column { c.FieldName = name; return &c }

func (c Column) On(t *Table) *Column { c.Table = t; return &c }

func (c *Column) T() *Table { return c.Table }

func (c *Column) ValueBy(v interface{}) *Assignment { return ColumnsAndValues(c, v) }

func (c *Column) Inc(d int) SqlExpr { return Expr("?+?", c, d) }

func (c *Column) Dec(d int) SqlExpr { return Expr("?-?", c, d) }

func (c *Column) Like(v string) SqlCondition { return AsCond(Expr("? LIKE ?", c, "%"+v+"%")) }

func (c *Column) LLike(v string) SqlCondition { return AsCond(Expr("? LIKE ?", c, "%"+v)) }

func (c *Column) RLike(v string) SqlCondition { return AsCond(Expr("? LIKE ?", c, v+"%")) }

func (c *Column) NotLike(v string) SqlCondition { return AsCond(Expr("? NOT LIKE ?", c, "%"+v+"%")) }

func (c *Column) IsNull() SqlCondition { return AsCond(Expr("? IS NULL", c)) }

func (c *Column) IsNotNull() SqlCondition { return AsCond(Expr("? IS NOT NULL", c)) }

func (c *Column) Eq(v interface{}) SqlCondition { return AsCond(Expr("? = ?", c, v)) }

func (c *Column) Neq(v interface{}) SqlCondition { return AsCond(Expr("? <> ?", c, v)) }

func (c *Column) Gt(v interface{}) SqlCondition { return AsCond(Expr("? > ?", c, v)) }

func (c *Column) Gte(v interface{}) SqlCondition { return AsCond(Expr("? >= ?", c, v)) }

func (c *Column) Lt(v interface{}) SqlCondition { return AsCond(Expr("? < ?", c, v)) }

func (c *Column) Lte(v interface{}) SqlCondition { return AsCond(Expr("? <= ?", c, v)) }

func (c *Column) Between(l, r interface{}) SqlCondition {
	return AsCond(Expr("? BETWEEN ? AND ?", c, l, r))
}

func (c *Column) NotBetween(l, r interface{}) SqlCondition {
	return AsCond(Expr("? NOT BETWEEN ? AND ?", c, l, r))
}

func (c *Column) In(args ...interface{}) SqlCondition {
	n := len(args)
	if n == 0 {
		return nil
	}
	if n == 1 {
		_ = n
		// TODO WithConditionFor this column
	}
	e := Expr("? IN ")
	e.Grow(n + 1)
	e.AppendArgs(c)
	e.WriteGroup(func(e *Ex) {
		for i := 0; i < n; i++ {
			e.WriteHolder(i)
		}
	})
	e.AppendArgs(args...)
	return AsCond(e)
}

func (c *Column) NotIn(args ...interface{}) SqlCondition {
	n := len(args)
	if n == 0 {
		return nil
	}
	e := Expr("? NOT IN ")
	e.Grow(n + 1)
	e.AppendArgs(c)
	e.WriteGroup(func(e *Ex) {
		for i := 0; i < n; i++ {
			e.WriteHolder(i)
		}
	})
	e.AppendArgs(args...)
	return AsCond(e)
}
