package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"net/url"
	"path"
	"path/filepath"
	"runtime"
	"strconv"

	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func main() {
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(
		fset,
		path.Join(must.String(pkgx.PkgPathByPath("net/http")), "status.go"),
		nil,
		parser.ParseComments,
	)
	// http.StatusOK

	names := make([]string, 0)
	codes := make([]int, 0)

	ast.Inspect(file, func(node ast.Node) bool {
		n, ok := node.(*ast.ValueSpec)
		if !ok {
			return true
		}
		if len(n.Values) != 1 {
			return false
		}
		lit, ok := n.Values[0].(*ast.BasicLit)
		if !ok {
			return false
		}
		if lit.Kind != token.INT {
			return false
		}
		if len(lit.Value) != 3 || lit.Value[0] != '3' {
			return false
		}
		name := n.Names[0].Name
		if name[0] == '_' {
			return false
		}
		names = append(names, name)
		code, _ := strconv.ParseInt(lit.Value, 10, 64)
		codes = append(codes, int(code))
		return false
	})

	write(names, codes)
}

func write(names []string, codes []int) {
	{
		f := g.NewFile("httpx", g.GenerateFileSuffix("./httpx_redirect.go"))

		for _, name := range names {
			f.WriteSnippet(
				f.Expr(`
func RedirectWith`+name+`(u *?) *`+name+` {
	return &`+name+`{Response: &Response{Location: u}}
}

type `+name+` struct { *Response }

func (`+name+`) StatusCode() int { return ? }

func (r `+name+`) Location() *? { return r.Response.Location }
`,
					g.Ident(f.Use("net/url", "URL")),
					g.Ident(f.Use("net/http", name)),
					g.Ident(f.Use("net/url", "URL")),
				),
			)
		}
		if _, err := f.Write(); err != nil {
			panic(err)
		}
	}

	{

		f := g.NewFile("httpx_test", g.GenerateFileSuffix("./httpx_redirect_test.go"))

		for i, name := range names {
			f.WriteSnippet(
				f.Expr(`func Example`+name+`() {
m := `+f.Use(pkg, `RedirectWith`+name)+`(?)
`+f.Use("fmt", "Println")+`(m.StatusCode())
`+f.Use("fmt", "Println")+`(m.Location())
// Output:
// ?
// /test
}`, &url.URL{Path: "/test"}, codes[i]),
			)

		}
		if _, err := f.Write(); err != nil {
			panic(err)
		}
	}
}

var pkg string

func init() {
	_, current, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(current), "../../../httpx")
	pkg = must.String(pkgx.PkgIdByPath(dir))
}
