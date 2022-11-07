package builder

import "context"

type onconflict struct {
	cols        *Columns
	doNothing   bool
	assignments []*Assignment
	AdditionType
}

func OnConflict(cols *Columns) *onconflict {
	return &onconflict{AdditionType: AdditionOnConflict, cols: cols}
}

func (a onconflict) DoNothing() *onconflict { a.doNothing = true; return &a }

func (a onconflict) DoUpdateSet(ass ...*Assignment) *onconflict {
	a.assignments = ass
	return &a
}

func (a *onconflict) IsNil() bool {
	return a == nil || IsNilExpr(a.cols) || (!a.doNothing && len(a.assignments) == 0)
}

func (a *onconflict) Ex(ctx context.Context) *Ex {
	e := Expr("ON CONFLICT ")
	e.WriteGroup(func(e *Ex) {
		e.WriteExpr(a.cols)
	})
	e.WriteQuery(" DO ")

	if a.doNothing {
		e.WriteQuery("NOTHING")
	} else {
		e.WriteQuery("UPDATE SET ")
		for i := range a.assignments {
			if i > 0 {
				e.WriteQuery(", ")
			}
			e.WriteExpr(a.assignments[i])
		}
	}
	return e.Ex(ctx)
}
