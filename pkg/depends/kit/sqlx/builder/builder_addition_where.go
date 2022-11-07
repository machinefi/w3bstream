package builder

import "context"

type where struct {
	AdditionType
	SqlCondition
}

func Where(c SqlCondition) *where { return &where{AdditionWhere, c} }

func (w *where) IsNil() bool { return w == nil || IsNilExpr(w.SqlCondition) }

func (w *where) Ex(ctx context.Context) *Ex {
	e := Expr("WHERE ")
	e.WriteExpr(w.SqlCondition)
	return e.Ex(ctx)
}
