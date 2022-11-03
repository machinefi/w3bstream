package builder

import "context"

type orderby struct {
	orders []*Order
	AdditionType
}

func (o *orderby) IsNil() bool { return o == nil || len(o.orders) == 0 }

func (o *orderby) Ex(ctx context.Context) *Ex {
	e := Expr("ORDER BY ")
	for i := range o.orders {
		if i > 0 {
			e.WriteQueryByte(',')
		}
		e.WriteExpr(o.orders[i])
	}
	return e.Ex(ctx)
}

func OrderBy(orders ...*Order) *orderby {
	o := &orderby{AdditionType: AdditionOrderBy}
	for i := range orders {
		if !IsNilExpr(orders[i]) {
			o.orders = append(o.orders, orders[i])
		}
	}
	return o
}

type Order struct {
	target SqlExpr
	order  string
}

func (o *Order) IsNil() bool { return o == nil || IsNilExpr(o.target) }

func (o *Order) Ex(ctx context.Context) *Ex {
	e := Expr("")
	e.Grow(1)
	e.WriteGroup(func(e *Ex) {
		e.WriteExpr(o.target)
	})
	if o.order != "" {
		e.WriteQueryByte(' ')
		e.WriteQuery(o.order)
	}
	return e.Ex(ctx)
}

func AscOrder(target SqlExpr) *Order {
	return &Order{target: target, order: "ASC"}
}
func DescOrder(target SqlExpr) *Order {
	return &Order{target: target, order: "DESC"}
}
