package httptransport

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/julienschmidt/httprouter"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

type MethodDescriber interface {
	Method() string
}

type PathDescriber interface {
	Path() string
}

type BasePathDescriber interface {
	BasePath() string
}

var pkgHTTPx = reflect.TypeOf(httpx.MethodGet{}).PkgPath()

func NewOperatorFactoryWithRouteMeta(op kit.Operator, last bool) *OperatorFactoryWithRouteMeta {
	f := kit.NewOperatorFactory(op, last)

	m := &OperatorFactoryWithRouteMeta{OperatorFactory: f}

	m.ID = m.Type.Name()

	if with, ok := op.(MethodDescriber); ok {
		m.Method = with.Method()
	}

	if m.Type.Kind() == reflect.Struct {
		for i := 0; i < m.Type.NumField(); i++ {
			ft := m.Type.Field(i)
			if !ft.Anonymous || ft.Type.PkgPath() != pkgHTTPx ||
				!strings.HasPrefix(ft.Name, "Method") {
				continue
			}
			// here can parse output operator
			if path, ok := ft.Tag.Lookup("path"); ok {
				vs := strings.Split(path, ",")
				m.Path = vs[0]

				if len(vs) > 0 {
					for i := range vs {
						switch vs[i] {
						case "deprecated":
							m.Deprecated = true
						}
					}
				}
			}

			if basePath, ok := ft.Tag.Lookup("basePath"); ok {
				m.BasePath = basePath
			}

			if summary, ok := ft.Tag.Lookup("summary"); ok {
				m.Summary = summary
			}

			break
		}
	}

	if with, ok := op.(BasePathDescriber); ok {
		m.BasePath = with.BasePath()
	}

	if with, ok := m.Operator.(PathDescriber); ok {
		m.Path = with.Path()
	}

	return m
}

type RouteMeta struct {
	ID         string // ID operator name
	Method     string // Method http method implement MethodDescriber
	Path       string // Path in tag `path`
	BasePath   string // BasePath base path
	Summary    string // Summary operator's desc
	Deprecated bool
}

type OperatorFactoryWithRouteMeta struct {
	*kit.OperatorFactory
	RouteMeta
}

func NewHttpRouteMeta(route *kit.Route) *HttpRouteMeta {
	metas := make([]*OperatorFactoryWithRouteMeta, len(route.Operators))

	for i := range route.Operators {
		metas[i] = NewOperatorFactoryWithRouteMeta(
			route.Operators[i],
			i == len(route.Operators)-1,
		)
	}

	return &HttpRouteMeta{
		Route: route,
		Metas: metas,
	}
}

type HttpRouteMeta struct {
	Route *kit.Route
	Metas []*OperatorFactoryWithRouteMeta
}

func (hr *HttpRouteMeta) OperatorNames() string {
	names := make([]string, 0)

	for _, m := range hr.Metas {
		if m.NoOutput {
			continue
		}
		if m.IsLast {
			names = append(names, color.MagentaString(m.String()))
		} else {
			names = append(names, color.CyanString(m.String()))
		}
	}

	return strings.Join(names, " ")
}

func (hr *HttpRouteMeta) Key() string {
	return regexpHttpRouterPath.ReplaceAllString(hr.Path(), "/{$1}") + " " + hr.OperatorNames()
}

func (hr *HttpRouteMeta) String() string {
	method := hr.Method()

	return methodColor(method)("%s %s", method[0:3], hr.Key())
}

func (hr *HttpRouteMeta) Log() {
	method := hr.Method()

	last := hr.Metas[len(hr.Metas)-1]

	firstLine := methodColor(method)(
		"%s %s", method[0:3],
		regexpHttpRouterPath.ReplaceAllString(hr.Path(), "/{$1}"))

	if last.Deprecated {
		firstLine = firstLine + " Deprecated"
	}

	if last.Summary != "" {
		firstLine = firstLine + " " + last.Summary
	}

	outputln(firstLine)
	outputln("\t%s", hr.OperatorNames())
}

var regexpHttpRouterPath = regexp.MustCompile("/:([^/]+)")

func (hr *HttpRouteMeta) Method() string {
	method := ""
	for _, m := range hr.Metas {
		if m.Method != "" {
			method = m.Method
		}
	}
	return method
}

func (hr *HttpRouteMeta) Path() string {
	basePath := "/"
	p := ""

	for _, m := range hr.Metas {
		if m.BasePath != "" {
			basePath = m.BasePath
		}

		if m.Path != "" {
			p += m.Path
		}
	}

	return httprouter.CleanPath(basePath + p)
}

type MetaOperator struct {
	kit.EmptyOperator
	path     string
	basePath string
}

func BasePath(basePath string) *MetaOperator { return &MetaOperator{basePath: basePath} }

func Group(path string) *MetaOperator { return &MetaOperator{path: path} }

func (g *MetaOperator) Path() string { return g.path }

func (g *MetaOperator) BasePath() string { return g.basePath }
