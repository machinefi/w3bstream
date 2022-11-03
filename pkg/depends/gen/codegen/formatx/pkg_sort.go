package formatx

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func SortImports(fset *token.FileSet, f *ast.File, file, group string) error {
	ast.SortImports(fset, f)
	dir := filepath.Dir(file)

	for _, decl := range f.Decls {
		d, ok := decl.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT || len(d.Specs) == 0 {
			break
		}
		g := &GroupSet{}
		for i := range d.Specs {
			s := d.Specs[i].(*ast.ImportSpec)
			path, _ := strconv.Unquote(s.Path.Value)

			if Stds[path] {
				g.AppendStd(path, s)
				continue
			}
			if group != "" && strings.HasPrefix(path, group) {
				g.AppendGroup(path, s)
				continue
			}
			if strings.Contains(dir, path) {
				g.AppendLocal(path, s)
				continue
			}
			g.AppendVendor(path, s)
		}
		_fset, _f, err := Parse(file, bytes.Replace(
			FmtNode(fset, f),
			FmtNode(fset, d),
			g.Bytes(),
			1,
		))
		if err != nil {
			return err
		}
		*fset, *f = *_fset, *_f
	}

	return nil
}

func FmtNode(fset *token.FileSet, node ast.Node) []byte {
	buf := bytes.NewBuffer(nil)
	if err := format.Node(buf, fset, node); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

type dep struct {
	Pkg  string
	Spec *ast.ImportSpec
}

type GroupSet [4][]*dep

func (g GroupSet) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("import ")
	if g.Len() > 1 {
		buf.WriteRune('(')
	}
	for _, deps := range g {
		for _, d := range deps {
			buf.WriteRune('\n')
			if d.Spec.Doc != nil {
				for _, c := range d.Spec.Doc.List {
					buf.WriteString(c.Text)
					buf.WriteRune('\n')
				}
			}
			if d.Spec.Name != nil {
				buf.WriteString(d.Spec.Name.String())
				buf.WriteRune(' ')
			}
			buf.WriteString(d.Spec.Path.Value)
			if d.Spec.Comment != nil {
				for _, c := range d.Spec.Comment.List {
					buf.WriteString(c.Text)
				}
			}
		}
		buf.WriteRune('\n')
	}
	if g.Len() > 1 {
		buf.WriteRune(')')
	}

	return buf.Bytes()
}

func (g *GroupSet) Len() int {
	return len(g[0]) + len(g[1]) + len(g[2]) + len(g[3])
}

func (g *GroupSet) append(idx int, pkg string, spec *ast.ImportSpec) {
	(*g)[idx] = append((*g)[idx], &dep{Pkg: pkg, Spec: spec})
}

func (g *GroupSet) AppendStd(pkg string, spec *ast.ImportSpec)    { g.append(0, pkg, spec) }
func (g *GroupSet) AppendVendor(pkg string, spec *ast.ImportSpec) { g.append(1, pkg, spec) }
func (g *GroupSet) AppendLocal(pkg string, spec *ast.ImportSpec)  { g.append(2, pkg, spec) }
func (g *GroupSet) AppendGroup(pkg string, spec *ast.ImportSpec)  { g.append(3, pkg, spec) }

type StdLibSet map[string]bool

func (s StdLibSet) WalkInit(root, prefix string) {
	ds, _ := ioutil.ReadDir(root)
	for _, d := range ds {
		if !d.IsDir() {
			continue
		}
		name := d.Name()
		if name == "vendor" || name == "internal" || name == "testdata" {
			continue
		}
		pkg := name
		if prefix != "" {
			pkg = filepath.Join(prefix, pkg)
		}
		s.WalkInit(filepath.Join(root, name), pkg)
		s[pkg] = true
	}
}

var Stds StdLibSet

func init() {
	Stds = make(StdLibSet)
	Stds.WalkInit(filepath.Join(runtime.GOROOT(), "src"), "")
	// root := filepath.Join(runtime.GOROOT(), "src")
	// err := filepath.Walk(
	// 	root,
	// 	func(path string, info fs.FileInfo, err error) error {
	// 		if path == root {
	// 			return nil
	// 		}
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if !info.IsDir() {
	// 			return nil
	// 		}
	// 		Stds[path] = true
	// 		return nil
	// 	},
	// )
	// if err != nil {
	// 	panic(err)
	// }
}
