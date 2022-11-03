package builder

import "context"

type SelectStatement interface {
	SqlExpr
	selectStatement()
}

type StmtSelect struct {
	SelectStatement

	expr      SqlExpr
	tbl       *Table
	modifiers []string
	adds      []Addition
}

func Select(e SqlExpr, modifiers ...string) *StmtSelect {
	return &StmtSelect{expr: e, modifiers: modifiers}
}

func (s *StmtSelect) IsNil() bool { return s == nil }

func (s StmtSelect) From(tbl *Table, adds ...Addition) *StmtSelect {
	s.tbl, s.adds = tbl, adds
	return &s
}

func (s *StmtSelect) Ex(ctx context.Context) *Ex {
	multi := false

	for i := range s.adds {
		add := s.adds[i]
		if IsNilExpr(add) {
			continue
		}
		if add.AdditionKind() == AdditionJoin {
			multi = true
		}
	}
	if multi {
		ctx = ContextWithToggleMultiTable(ctx, true)
	}

	e := Expr("SELECT")
	e.Grow(len(s.adds) + 2)

	if len(s.modifiers) > 0 {
		for i := range s.modifiers {
			e.WriteQueryByte(' ')
			e.WriteQuery(s.modifiers[i])
		}
	}
	expr := s.expr
	if IsNilExpr(expr) {
		expr = Expr("*")
	}
	e.WriteQueryByte(' ')
	e.WriteExpr(expr)

	if !IsNilExpr(s.tbl) {
		e.WriteQuery(" FROM ")
		e.WriteExpr(s.tbl)
	}
	WriteAdditions(e, s.adds...)
	return e.Ex(ctx)
}

func ForUpdate() *addition { return AsAddition(Expr("FOR UPDATE")) }
