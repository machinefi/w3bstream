package pkgx

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

func StringifyNode(fs *token.FileSet, n ast.Node) string {
	if cg, ok := n.(*ast.CommentGroup); ok {
		return StringifyCommentGroup(cg)
	}
	buf := bytes.NewBuffer(nil)
	if err := format.Node(buf, fs, n); err != nil {
		panic(err)
	}
	return buf.String()
}

func StringifyCommentGroup(cgs ...*ast.CommentGroup) string {
	if len(cgs) == 0 {
		return ""
	}
	comments := ""
	for _, cg := range cgs {
		for _, line := range strings.Split(cg.Text(), "\n") {
			if strings.HasPrefix(line, "go:") {
				continue
			}
			comments = comments + "\n" + line
		}
	}
	return strings.TrimSpace(comments)
}

func Import(path string) string {
	parts := strings.Split(path, "/vendor/")
	return parts[len(parts)-1]
}

func ImportPathAndExpose(s string) (string, string) {
	args := strings.Split(s, ".")
	if _len := len(args); _len > 1 {
		return Import(strings.Join(args[0:_len-1], ".")), args[_len-1]
	}
	return "", s
}

const (
	ModeLoadFiles = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles
)

func FindPkgInfoByPath(path string, modes ...packages.LoadMode) (*packages.Package, error) {
	mode := ModeLoadFiles
	if len(modes) > 0 {
		for _, v := range modes {
			mode |= v
		}
	}
	pkgs, err := packages.Load(&packages.Config{Mode: ModeLoadFiles}, path)
	if err != nil {
		panic(err)
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("package `%s` not found", path)
	}
	return pkgs[0], nil
}

func PkgIdByPath(path string, modes ...packages.LoadMode) (string, error) {
	pkg, err := FindPkgInfoByPath(path, modes...)
	if err != nil {
		return "", err
	}
	return pkg.ID, nil
}

func PkgPathByPath(path string, modes ...packages.LoadMode) (string, error) {
	pkg, err := FindPkgInfoByPath(path, modes...)
	if err != nil {
		return "", err
	}
	return filepath.Dir(pkg.GoFiles[0]), nil
}
