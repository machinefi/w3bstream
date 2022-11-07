package builder

import "context"

type StmtInsert struct {
	tbl         *Table
	modifiers   []string
	assignments []*Assignment
	adds        Additions
}

func Insert(modifiers ...string) *StmtInsert {
	return &StmtInsert{modifiers: modifiers}
}

func (s StmtInsert) Into(tbl *Table, adds ...Addition) *StmtInsert {
	s.tbl, s.adds = tbl, adds
	return &s
}

func (s StmtInsert) Values(cols *Columns, values ...interface{}) *StmtInsert {
	s.assignments = Assignments{ColumnsAndValues(cols, values...)}
	return &s
}

func (s *StmtInsert) IsNil() bool {
	return s == nil || s.tbl == nil || len(s.assignments) == 0
}

func (s *StmtInsert) Ex(ctx context.Context) *Ex {
	e := Expr("INSERT")
	if len(s.modifiers) > 0 {
		for i := range s.modifiers {
			e.WriteQueryByte(' ')
			e.WriteQuery(s.modifiers[i])
		}
	}
	e.WriteQuery(" INTO ")
	e.WriteExpr(s.tbl)
	e.WriteQueryByte(' ')

	e.WriteExpr(ExprBy(func(ctx context.Context) *Ex {
		e := Expr("")
		e.Grow(len(s.assignments))
		ctx = ContextWithToggleUseValues(ctx, true)
		WriteAssignments(e, s.assignments...)
		return e.Ex(ctx)
	}))
	WriteAdditions(e, s.adds...)
	return e.Ex(ctx)
}

// TODO OnDuplicateKeyUpdate (mysql feature)
// TODO Returning
