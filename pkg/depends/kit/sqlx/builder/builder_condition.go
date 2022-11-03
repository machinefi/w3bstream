package builder

import "context"

type SqlConditionMarker interface {
	asCondition()
}

type SqlCondition interface {
	SqlExpr
	SqlConditionMarker

	And(SqlCondition) SqlCondition
	Or(SqlCondition) SqlCondition
	Xor(SqlCondition) SqlCondition
}

var EmptyCondition SqlCondition = (*Condition)(nil)

type Condition struct {
	expr SqlExpr
	SqlConditionMarker
}

func (c *Condition) IsNil() bool { return c == nil || c.expr == nil }

func (c *Condition) Ex(ctx context.Context) *Ex {
	if IsNilExpr(c.expr) {
		return nil
	}
	return c.expr.Ex(ctx)
}

func (c *Condition) And(cond SqlCondition) SqlCondition {
	if IsNilExpr(cond) {
		return c
	}
	return And(c, cond)
}

func (c *Condition) Or(cond SqlCondition) SqlCondition {
	if IsNilExpr(cond) {
		return c
	}
	return Or(c, cond)
}

func (c *Condition) Xor(cond SqlCondition) SqlCondition {
	if IsNilExpr(cond) {
		return c
	}
	return Xor(c, cond)
}

func And(conds ...SqlCondition) SqlCondition {
	return ComposeConditions("AND", conds...)
}

func Or(conds ...SqlCondition) SqlCondition {
	return ComposeConditions("OR", conds...)
}

func Xor(conds ...SqlCondition) SqlCondition {
	return ComposeConditions("XOR", conds...)
}

type CondCompose struct {
	SqlConditionMarker

	op    string
	conds []SqlCondition
}

func ComposeConditions(op string, conds ...SqlCondition) *CondCompose {
	c := &CondCompose{op: op}
	for i := range conds {
		cond := conds[i]
		if IsNilExpr(cond) {
			continue
		}
		c.conds = append(c.conds, cond)
	}
	return c
}

func (c *CondCompose) IsNil() bool {
	if c == nil || c.op == "" {
		return true
	}
	for i := range c.conds {
		if !IsNilExpr(c.conds[i]) {
			return false
		}
	}
	return true
}

func (c *CondCompose) Ex(ctx context.Context) *Ex {
	e := Expr("")
	for i := range c.conds {
		if i > 0 {
			e.WriteQueryByte(' ')
			e.WriteQuery(c.op)
			e.WriteQueryByte(' ')
		}
		e.WriteGroup(func(e *Ex) {
			e.WriteExpr(c.conds[i])
		})
	}
	return e.Ex(ctx)
}

func (c *CondCompose) And(cond SqlCondition) SqlCondition { return And(c, cond) }

func (c *CondCompose) Or(cond SqlCondition) SqlCondition { return Or(c, cond) }

func (c *CondCompose) Xor(cond SqlCondition) SqlCondition { return Xor(c, cond) }

func AsCond(e SqlExpr) *Condition { return &Condition{expr: e} }
