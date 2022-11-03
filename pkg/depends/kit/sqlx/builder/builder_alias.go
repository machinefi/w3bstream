package builder

import "context"

type alias struct {
	name string
	SqlExpr
}

func Alias(e SqlExpr, name string) *alias {
	return &alias{name: name, SqlExpr: e}
}

func (as *alias) IsNil() bool {
	return as == nil || as.name == "" || IsNilExpr(as.SqlExpr)
}

func (as *alias) Ex(ctx context.Context) *Ex {
	return Expr(
		"? AS ?",
		as.SqlExpr,
		Expr(as.name),
	).Ex(ContextWithToggleNeedAutoAlias(ctx, false))
}
