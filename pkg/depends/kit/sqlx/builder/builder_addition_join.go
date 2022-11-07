package builder

import (
	"context"
	"strings"
)

type join struct {
	AdditionType

	prefix   string
	target   SqlExpr
	joinCond SqlCondition
	joinCols []*Column
}

func (j join) On(c SqlCondition) *join { j.joinCond = c; return &j }

func (j join) Using(cols ...*Column) *join { j.joinCols = cols; return &j }

func (j *join) IsNil() bool {
	return j == nil || IsNilExpr(j.target) ||
		(j.prefix != "CROSS" && IsNilExpr(j.joinCond) && len(j.joinCols) == 0)
}

func (j *join) Ex(ctx context.Context) *Ex {
	e := Expr("JOIN ")
	if j.prefix != "" {
		e = Expr(j.prefix + " JOIN ")
	}
	e.WriteExpr(j.target)
	if !IsNilExpr(j.joinCond) {
		e.WriteExpr(ExprBy(func(ctx context.Context) *Ex {
			e := Expr(" ON ")
			e.WriteExpr(j.joinCond)
			return e.Ex(ctx)
		}))
	}
	if len(j.joinCols) > 0 {
		e.WriteExpr(ExprBy(func(ctx context.Context) *Ex {
			e := Expr(" USING ")
			e.WriteGroup(func(e *Ex) {
				for i := range j.joinCols {
					if i > 0 {
						e.WriteQuery(", ")
					}
					e.WriteExpr(j.joinCols[i])
				}
			})
			return e.Ex(ContextWithToggleMultiTable(ctx, false))
		}))
	}
	return e.Ex(ctx)
}

func Join(tar SqlExpr, prefixes ...string) *join {
	return &join{
		AdditionType: AdditionJoin,
		prefix:       strings.Join(prefixes, " "),
		target:       tar,
	}
}

func InnerJoin(tar SqlExpr) *join { return Join(tar, "INNER") }
func LeftJoin(tar SqlExpr) *join  { return Join(tar, "LEFT") }
func RightJoin(tar SqlExpr) *join { return Join(tar, "RIGHT") }
func FullJoin(tar SqlExpr) *join  { return Join(tar, "FULL") }
func CrossJoin(tar SqlExpr) *join { return Join(tar, "CROSS") }
