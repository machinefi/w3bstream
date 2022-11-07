package builder

type addition struct {
	AdditionType // immutable
	SqlExpr
}

func AsAddition(e SqlExpr) *addition { return &addition{AdditionOther, e} }

func (a *addition) IsNil() bool { return a == nil || IsNilExpr(a.SqlExpr) }
