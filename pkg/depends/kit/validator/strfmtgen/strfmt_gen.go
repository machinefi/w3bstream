package strfmtgen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/kit/validator"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

func NewGenerator(pkg *pkgx.Pkg, path string) *Generator {
	return &Generator{
		pkg:  pkg,
		path: path,
	}
}

type Generator struct {
	pkg  *pkgx.Pkg
	path string
}

func (sg *Generator) Output(_ string) {
	path, _ := filepath.Abs(sg.path)
	fset := token.NewFileSet()
	fast, _ := parser.ParseFile(fset, path, nil, parser.ParseComments)

	parts := strings.SplitN(filepath.Base(path), ".", 2)
	keyPrefix := "regexpString"

	file := g.NewFile(sg.pkg.Name, filepath.Join(filepath.Dir(path), parts[0]+"_generated.go"))

	regexps := make([]string, 0)
	for key, obj := range fast.Scope.Objects {
		if obj.Kind != ast.Con && obj.Name == "_" {
			continue
		}
		if strings.HasPrefix(key, keyPrefix) {
			regexps = append(regexps, key)
		}
	}
	sort.Strings(regexps)

	assigns := make([]g.SnippetSpec, 0)
	registers := make([]g.Snippet, 0)
	for _, key := range regexps {
		var (
			ident  = strings.TrimPrefix(key, keyPrefix)
			prefix = stringsx.UpperCamelCase(ident)
			name   = strings.Replace(stringsx.LowerSnakeCase(ident), "_", "-", -1)
		)

		args := []g.Snippet{g.Ident(key), g.Valuer(name)}
		if alias := stringsx.LowerCamelCase(name); alias != name {
			args = append(args, g.Valuer(alias))
		}

		assigns = append(assigns, g.Assign(g.Var(nil, prefix+"Validator")).
			By(g.Call(file.Use(pkgVldt, "NewRegexpStrfmtValidator"), args...)),
		)
		registers = append(registers, g.Ref(
			g.Ident(file.Use(pkgVldt, "DefaultFactory")),
			g.Call(
				"Register",
				g.Ident(prefix+"Validator"),
			),
		))
	}

	file.WriteSnippet(
		g.DeclVar(assigns...),
		g.Func().Named("init").Do(registers...),
	)
	_, _ = file.Write()
}

var pkgVldt = pkgx.Import(reflect.TypeOf(validator.Rule{}).PkgPath())
