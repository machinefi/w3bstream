package buildertestutil

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

func BeExpr(q string, args ...interface{}) types.GomegaMatcher {
	return &SqlExprMatcher{
		QueryMatcher: gomega.Equal(strings.TrimSpace(q)),
		ArgsMatcher:  gomega.Equal(args),
	}
}

type SqlExprMatcher struct {
	QueryMatcher types.GomegaMatcher
	ArgsMatcher  types.GomegaMatcher
}

func (m *SqlExprMatcher) Match(actual interface{}) (success bool, err error) {
	e, ok := actual.(builder.SqlExpr)
	if !ok {
		return false, fmt.Errorf("actual should be SqlExpr")
	}
	if builder.IsNilExpr(e) {
		return m.QueryMatcher.Match("")
	}
	ex := e.Ex(context.Background())
	queryMatched, err := m.QueryMatcher.Match(ex.Query())
	if err != nil {
		return false, err
	}
	argsMatched, err := m.ArgsMatcher.Match(ex.Args())
	if err != nil {
		return false, err
	}

	return queryMatched && argsMatched, nil
}

func (m *SqlExprMatcher) FailureMessage(actual interface{}) (msg string) {
	e := actual.(builder.SqlExpr).Ex(context.Background())
	return m.QueryMatcher.FailureMessage(e.Query()) + "\n" +
		m.ArgsMatcher.FailureMessage(e.Args())
}

func (m *SqlExprMatcher) NegatedFailureMessage(actual interface{}) (msg string) {
	e := actual.(builder.SqlExpr).Ex(context.Background())
	return m.QueryMatcher.NegatedFailureMessage(e.Query()) + "\n" +
		m.ArgsMatcher.NegatedFailureMessage(e.Args())
}

func MustCols(cols *builder.Columns, err error) *builder.Columns {
	return cols
}

type Point struct {
	X float64
	Y float64
}

func (Point) DataType(engine string) string { return "POINT" }

func (Point) ValueEx() string { return `ST_GeomFromText(?)` }

func (p Point) Value() (driver.Value, error) { return fmt.Sprintf("POINT(%v %v)", p.X, p.Y), nil }
