package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"

	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

func main() {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "strfmt.go", nil, parser.ParseComments)
	file := g.NewFile("strfmt", "strfmt_generated.go")

	regexps := make([]string, 0)
	for key, obj := range f.Scope.Objects {
		if obj.Kind == ast.Con {
			regexps = append(regexps, key)
		}
	}
	sort.Strings(regexps)

	snippets := make([]g.Snippet, 0)
	for _, key := range regexps {
		var (
			name           = strings.Replace(key, "regexpString", "", 1)
			validatorName  = strings.Replace(stringsx.LowerSnakeCase(name), "_", "-", -1)
			validatorAlias = stringsx.LowerCamelCase(name)
			args           = []g.Snippet{g.Ident(key), g.Valuer(validatorName)}
			prefix         = stringsx.UpperCamelCase(name)
			snippet        g.Snippet
		)
		if validatorName != validatorAlias {
			args = append(args, g.Valuer(validatorAlias))
		}
		snippet = g.Func().Named("init").Do(
			g.Ref(
				g.Ident(file.Use(pkg, "DefaultFactory")),
				g.Call(
					"Register",
					g.Ident(prefix+"Validator"),
				),
			),
		)
		snippets = append(snippets, snippet)
		snippet = g.DeclVar(
			g.Assign(g.Var(nil, prefix+"Validator")).
				By(g.Call(file.Use(pkg, "NewRegexpStrfmtValidator"), args...)),
		)
		snippets = append(snippets, snippet)

	}
	file.WriteSnippet(snippets...)
	_, _ = file.Write()
}

var pkg = "github.com/machinefi/w3bstream/pkg/depends/kit/validator"
