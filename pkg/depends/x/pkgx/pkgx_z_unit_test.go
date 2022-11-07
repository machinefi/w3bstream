package pkgx_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	. "github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func TestImportPathAndExpose(t *testing.T) {
	cases := []struct {
		imported string
		expose   string
		s        string
	}{
		{"", "B", "B"},
		{"testing", "B", "testing.B"},
		{"a.b.c.d/c", "B", "a.b.c.d/c.B"},
		{"e", "B", "a.b.c.d/vendor/e.B"},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			imported, expose := ImportPathAndExpose(c.s)
			NewWithT(t).Expect(imported).To(Equal(c.imported))
			NewWithT(t).Expect(expose).To(Equal(c.expose))
		})
	}
}

var root = "./__tests__"

func TestCommentScanner(t *testing.T) {
	fset := token.NewFileSet()

	fpth, err := filepath.Abs(path.Join(root, "comments.go"))
	if err != nil {
		t.Error(err)
		return
	}

	fast, err := parser.ParseFile(fset, fpth, nil, parser.ParseComments)
	if err != nil {
		t.Error(err)
		return
	}

	ast.Inspect(fast, func(node ast.Node) bool {
		comments := strings.Split(NewCommentScanner(fset, fast).CommentsOf(node), "\n")
		NewWithT(t).Expect(3 >= len(comments)).To(BeTrue())
		return true
	})
}

func TestPkgComments(t *testing.T) {
	cwd, err := os.Getwd()
	NewWithT(t).Expect(err).To(BeNil())
	pkg, err := LoadFrom(path.Join(cwd, root))
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(pkg.Imports()).NotTo(BeEmpty())

	for _, v := range []struct {
		object types.Object // identifier object
		expect string       // expect identifier's comment
	}{
		{pkg.TypeName("Date"), "type Date"},
		{pkg.Var("test"), "var"},
		{pkg.Const("A"), "a\n\nA"},
		{pkg.Func("Print"), "func Print"},
	} {
		NewWithT(t).
			Expect(pkg.CommentsOf(pkg.IdentOf(v.object))).
			To(Equal(v.expect))
	}
}

func TestPkgFuncReturns(t *testing.T) {
	cwd, err := os.Getwd()
	NewWithT(t).Expect(err).To(BeNil())
	pkg, err := LoadFrom(path.Join(cwd, root))
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(pkg.Imports()).NotTo(BeEmpty())

	var pkgid string

	{
		_, current, _, _ := runtime.Caller(0)
		dir := filepath.Join(filepath.Dir(current), "./__tests__")
		pkgid = must.String(PkgIdByPath(dir))
	}

	var cases = []struct {
		FuncName string
		Results  [][]string
	}{
		{
			"FuncSingleReturn",
			[][]string{{"untyped int(2)"}},
		},
		{
			"FuncSelectExprReturn",
			[][]string{{"string"}},
		},
		{
			"FuncWillCall",
			[][]string{
				{"interface{}"},
				{strings.Join([]string{pkgid, "String"}, ".")},
			},
		},
		{
			"FuncReturnWithCallDirectly",
			[][]string{
				{"interface{}"},
				{strings.Join([]string{pkgid, "String"}, ".")},
			},
		},
		{
			"FuncWithNamedReturn",
			[][]string{
				{"interface{}"},
				{strings.Join([]string{pkgid, "String"}, ".")},
			},
		},
		{
			"FuncSingleNamedReturnByAssign",
			[][]string{
				{`untyped string("1")`},
				{strings.Join([]string{pkgid, `String("2")`}, ".")},
			},
		},
		{
			"FuncWithSwitch",
			[][]string{
				{
					`untyped string("a1")`,
					`untyped string("a2")`,
					`untyped string("a3")`,
				},
				{
					strings.Join([]string{pkgid, `String("b1")`}, "."),
					strings.Join([]string{pkgid, `String("b2")`}, "."),
					strings.Join([]string{pkgid, `String("b3")`}, "."),
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.FuncName, func(t *testing.T) {
			values, n := pkg.FuncResultsOf(pkg.Func(c.FuncName))
			NewWithT(t).Expect(values).To(HaveLen(n))
			NewWithT(t).Expect(c.Results).To(Equal(PrintValues(pkg.Fset, values)))
		})
	}
}

func PrintValues(fs *token.FileSet, res map[int][]TypeAndValueExpr) [][]string {
	if res == nil {
		return [][]string{}
	}
	ret := make([][]string, len(res))
	for i := range ret {
		tve := res[i]
		ret[i] = make([]string, len(tve))
		for j, v := range tve {
			fmt.Println(v.Type, v.Value)
			if v.Value == nil {
				ret[i][j] = v.Type.String()
			} else {
				ret[i][j] = fmt.Sprintf("%s(%s)", v.Type, v.Value)
			}
		}
	}
	return ret
}

func PrintAstInfo(t *testing.T) {
	fset := token.NewFileSet()
	fpth, err := filepath.Abs(path.Join(root, "ast.go"))
	if err != nil {
		t.Error(err)
		return
	}
	fast, err := parser.ParseFile(fset, fpth, nil, parser.AllErrors)
	if err != nil {
		t.Error(err)
		return
	}
	_ = ast.Print(fset, fast)
}
