package kit

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type Router struct {
	parent    *Router
	operators []Operator
	children  map[*Router]bool
}

func NewRouter(operators ...Operator) *Router {
	ops := make([]Operator, 0)
	for i := range operators {
		op := operators[i]

		if with, ok := op.(WithMiddleOperators); ok {
			ops = append(ops, with.MiddleOperators()...)
		}

		ops = append(ops, op)
	}

	return &Router{
		operators: ops,
	}
}

func (r *Router) Register(x *Router) {
	if r.children == nil {
		r.children = map[*Router]bool{}
	}
	if x.parent != nil {
		panic(fmt.Errorf("router %v already registered to router %v", x, x.parent))
	}
	x.parent = r
	r.children[x] = true
}

func (r *Router) route() *Route {
	parent := r.parent
	operators := r.operators

	for parent != nil {
		operators = append(parent.operators, operators...)
		parent = parent.parent
	}

	return &Route{
		Operators: operators,
		last:      len(r.children) == 0,
	}
}

func (r *Router) Routes() (routes Routes) {
	maybeAppendRoute := func(router *Router) {
		route := router.route()

		if route.last && len(route.Operators) > 0 {
			routes = append(routes, route)
		}

		if len(router.children) > 0 {
			routes = append(routes, router.Routes()...)
		}
	}

	if len(r.children) == 0 {
		maybeAppendRoute(r)
		return
	}

	for child := range r.children {
		maybeAppendRoute(child)
	}

	return
}

type Routes []*Route

func (rs Routes) String() string {
	keys := make([]string, len(rs))
	for i, r := range rs {
		keys[i] = r.String()
	}
	sort.Strings(keys)
	return strings.Join(keys, "\n")
}

type Route struct {
	Operators []Operator
	last      bool
}

func (rs *Route) OperatorFactories() (factories []*OperatorFactory) {
	length := len(rs.Operators)
	for i, op := range rs.Operators {
		factories = append(factories, NewOperatorFactory(op, i == length-1))
	}
	return
}

func (rs *Route) String() string {
	buf := &bytes.Buffer{}
	factories := rs.OperatorFactories()
	for i, operatorFactory := range factories {
		if i > 0 {
			buf.WriteString(" |> ")
		}
		buf.WriteString(operatorFactory.String())
	}
	return buf.String()
}
