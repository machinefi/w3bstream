package builder

import "context"

type StmtUpdate struct {
	tbl         *Table
	modifiers   []string
	assignments []*Assignment
	adds        []Addition
}

func Update(tbl *Table, modifiers ...string) *StmtUpdate {
	return &StmtUpdate{tbl: tbl, modifiers: modifiers}
}

func (s *StmtUpdate) IsNil() bool {
	return s == nil || IsNilExpr(s.tbl) || len(s.assignments) == 0
}

func (s StmtUpdate) Set(assignments ...*Assignment) *StmtUpdate {
	s.assignments = assignments
	return &s
}

func (s StmtUpdate) Where(c SqlCondition, adds ...Addition) *StmtUpdate {
	s.adds = []Addition{Where(c)}
	if len(adds) > 0 {
		s.adds = append(s.adds, adds...)
	}
	return &s
}

func (s *StmtUpdate) Ex(ctx context.Context) *Ex {
	e := Expr("UPDATE")
	if len(s.modifiers) > 0 {
		for i := range s.modifiers {
			e.WriteQueryByte(' ')
			e.WriteQuery(s.modifiers[i])
		}
	}
	e.WriteQueryByte(' ')
	e.WriteExpr(s.tbl)
	e.WriteQuery(" SET ")

	WriteAssignments(e, s.assignments...)
	WriteAdditions(e, s.adds...)
	return e.Ex(ctx)
}

// TODO Update without condition warning
