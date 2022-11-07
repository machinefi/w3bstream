package formatx

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

type Proc func(fset *token.FileSet, f *ast.File, filename, group string) error

func Format(file, mod string, src []byte, procs ...Proc) ([]byte, error) {
	fset, f, err := Parse(file, src)
	if err != nil {
		return nil, err
	}

	for _, proc := range procs {
		if proc == nil {
			continue
		}
		if err = proc(fset, f, file, mod); err != nil {
			return nil, err
		}
	}

	buf := bytes.NewBuffer(nil)
	if err = format.Node(buf, fset, f); err != nil {
		return nil, fmt.Errorf("[FORMAT] %s: %v", file, err)
	}
	return buf.Bytes(), nil
}

func Parse(file string, src []byte) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, src, parser.ParseComments)
	if err != nil {
		return nil, nil, fmt.Errorf("[PARSE] %s: %v", file, err)
	}
	return fset, f, nil
}

func MustFormat(file, mod string, src []byte, procs ...Proc) []byte {
	code, err := Format(file, mod, src, procs...)
	if err != nil {
		panic(err)
	}
	return code
}
