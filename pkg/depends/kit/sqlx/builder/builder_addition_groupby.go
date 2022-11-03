package builder

import "context"

type groupby struct {
	AdditionType
	groups []SqlExpr
	having SqlCondition
}

func GroupBy(grps ...SqlExpr) *groupby {
	return &groupby{
		AdditionType: AdditionGroupBy,
		groups:       grps,
	}
}

func (g groupby) Having(c SqlCondition) *groupby {
	g.having = c
	return &g
}

func (g *groupby) IsNil() bool { return g == nil || len(g.groups) == 0 }

func (g *groupby) Ex(ctx context.Context) *Ex {
	e := Expr("GROUP BY ")

	for i, grp := range g.groups {
		if i > 0 {
			e.WriteQueryByte(',')
		}
		e.WriteExpr(grp)
	}
	if !IsNilExpr(g.having) {
		e.WriteQuery(" HAVING ")
		e.WriteExpr(g.having)
	}
	return e.Ex(ctx)
}
