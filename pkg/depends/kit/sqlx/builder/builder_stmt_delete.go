package builder

import "context"

func Delete() *StmtDelete { return &StmtDelete{} }

type StmtDelete struct {
	tbl  *Table
	adds Additions
}

func (s StmtDelete) From(tbl *Table, adds ...Addition) *StmtDelete {
	s.tbl, s.adds = tbl, adds
	return &s
}

func (s *StmtDelete) IsNil() bool { return s == nil || IsNilExpr(s.tbl) }

func (s *StmtDelete) Ex(ctx context.Context) *Ex {
	e := Expr("DELETE FROM ")
	e.WriteExpr(s.tbl)
	WriteAdditions(e, s.adds...)
	return e.Ex(ctx)
}
