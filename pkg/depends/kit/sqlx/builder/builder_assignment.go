package builder

import (
	"context"
	"math"
)

type Assignment struct {
	AssignmentMarker

	columns SqlExpr
	values  []interface{}
	colc    int
}

func ColumnsAndValues(columns SqlExpr, values ...interface{}) *Assignment {
	colc := 1
	if canLen, ok := columns.(interface{ Len() int }); ok {
		colc = canLen.Len()
	}
	return &Assignment{colc: colc, values: values, columns: columns}
}

func (a *Assignment) IsNil() bool {
	return a == nil || IsNilExpr(a.columns) || len(a.values) == 0
}

func (a *Assignment) Ex(ctx context.Context) *Ex {
	e := Expr("")
	e.Grow(len(a.values))

	if TogglesFromContext(ctx).Is(ToggleUseValues) {
		e.WriteGroup(func(e *Ex) {
			e.WriteExpr(ExprBy(func(ctx context.Context) *Ex {
				return a.columns.Ex(ContextWithToggleMultiTable(ctx, false))
			}))
		})
		if len(a.values) == 1 {
			if expr, ok := a.values[0].(SelectStatement); ok {
				e.WriteQueryByte(' ')
				e.WriteExpr(expr)
				return e.Ex(ctx)
			}
		}
		e.WriteQuery(" VALUES ")

		groupCount := int(math.Round(float64(len(a.values)) / float64(a.colc)))
		for i := 0; i < groupCount; i++ {
			if i > 0 {
				e.WriteQueryByte(',')
			}
			e.WriteGroup(func(e *Ex) {
				for j := 0; j < a.colc; j++ {
					e.WriteHolder(j)
				}
			})
		}
		e.AppendArgs(a.values...)
		return e.Ex(ctx)
	}

	e.WriteExpr(ExprBy(func(ctx context.Context) *Ex {
		return a.columns.Ex(ContextWithToggleMultiTable(ctx, false))
	}))
	e.WriteQuery(" = ?")
	e.AppendArgs(a.values[0])
	return e.Ex(ctx)
}

type Assignments []*Assignment

func WriteAssignments(e *Ex, assignments ...*Assignment) {
	count := 0
	for i := range assignments {
		a := assignments[i]
		if IsNilExpr(a) {
			continue
		}
		if count > 0 {
			e.WriteQuery(", ")
		}
		e.WriteExpr(a)
		count++
	}
}

type AssignmentMarker interface {
	asCondition()
}

type SqlAssignment interface {
	SqlExpr
	AssignmentMarker
}
