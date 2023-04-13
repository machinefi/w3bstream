package openapi

import (
	"bytes"
	"go/types"
	"sort"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(typeVar *types.Var, operators ...*OperatorWithTypeName) *Router {
	return &Router{
		typeVar:   typeVar,
		operators: operators,
	}
}

type Router struct {
	typeVar   *types.Var
	parent    *Router
	operators []*OperatorWithTypeName
	children  map[*Router]bool
}

func (r *Router) Name() string {
	if r.typeVar == nil {
		return "Anonymous"
	}
	return r.typeVar.Pkg().Name() + "." + r.typeVar.Name()
}

func (r *Router) String() string {
	buf := bytes.NewBufferString(r.Name())

	buf.WriteString("<")
	for i := range r.operators {
		o := r.operators[i]
		if i != 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(o.String())
	}
	buf.WriteString(">")

	buf.WriteString("[")

	i := 0
	for sub := range r.children {
		if i != 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(sub.Name())
		i++
	}
	buf.WriteString("]")

	return buf.String()
}

func (r *Router) AppendOperators(ops ...*OperatorWithTypeName) {
	r.operators = append(r.operators, ops...)
}

func (r *Router) With(operators ...*OperatorWithTypeName) {
	r.Register(NewRouter(nil, operators...))
}

func (r *Router) Register(child *Router) {
	if r.children == nil {
		r.children = map[*Router]bool{}
	}
	child.parent = r
	r.children[child] = true
}

func (r *Router) Route() *Route {
	parent := r.parent
	operators := r.operators

	for parent != nil {
		operators = append(parent.operators, operators...)
		parent = parent.parent
	}

	route := Route{
		last:      r.children == nil,
		Operators: operators,
	}

	return &route
}

func (r *Router) Routes() (routes []*Route) {
	for child := range r.children {
		route := child.Route()

		if route.last {
			routes = append(routes, route)
		}

		if child.children != nil {
			routes = append(routes, child.Routes()...)
		}
	}

	sort.Slice(routes, func(i, j int) bool {
		return routes[i].String() < routes[j].String()
	})

	return routes
}

type Route struct {
	Operators []*OperatorWithTypeName
	last      bool
}

func (r *Route) String() string {
	buf := bytes.NewBufferString(r.Method())
	buf.WriteString(" ")
	buf.WriteString(r.Path())

	for i := range r.Operators {
		buf.WriteString(" ")
		buf.WriteString(r.Operators[i].String())
	}

	return buf.String()
}

func (r *Route) Method() string {
	method := ""
	for _, m := range r.Operators {
		if m.Method != "" {
			method = m.Method
		}
	}
	return method
}

func (r *Route) Path() string {
	basePath := "/"
	fullPath := ""

	for _, operator := range r.Operators {
		if operator.BasePath != "" {
			basePath = operator.BasePath
		}
		if operator.Path != "" {
			fullPath += operator.Path
		}
	}

	return httprouter.CleanPath(basePath + fullPath)
}
