package schema

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

func Col(name string) *Column {
	return &Column{
		Name:       name,
		Constrains: &ColumnType{},
	}
}

type Column struct {
	Name       string      `json:"name"`       // Name column name
	Constrains *ColumnType `json:"constrains"` // Constrains column constrains
	exactly    bool

	WithTableDefinition `json:"-"`
}

var _ builder.SqlExpr = (*Column)(nil)

func (c *Column) IsNil() bool { return c == nil }

func (c *Column) Ex(ctx context.Context) *builder.Ex {
	toggles := builder.TogglesFromContext(ctx)
	if c.t != nil && (c.exactly || toggles.Is(builder.ToggleMultiTable)) {
		if toggles.Is(builder.ToggleNeedAutoAlias) {
			return builder.Expr(
				"?.? AS ?",
				c.t, builder.Expr(c.Name), builder.Expr(c.Name),
			).Ex(ctx)
		}
		return builder.Expr("?.?", c.t, builder.Expr(c.Name)).Ex(ctx)
	}
	return builder.ExactlyExpr(c.Name).Ex(ctx)
}
