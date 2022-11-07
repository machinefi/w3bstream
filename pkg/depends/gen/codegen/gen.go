package codegen

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"time"

	"golang.org/x/tools/go/packages"

	"github.com/machinefi/w3bstream/pkg/depends/gen/codegen/formatx"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

type File struct {
	Pkg         string
	Name        string
	Imps        map[string]string
	Pkgs        map[string][]string
	OrderedImps [][2]string
	opts        WriteOption
	bytes.Buffer
}

func NewFile(pkg, name string) *File {
	return &File{
		Pkg:  pkg,
		Name: name,
		opts: WriteOption{WithCommit: true, MustFormat: true},
	}
}

func (f *File) bytes() []byte {
	buf := bytes.NewBuffer(nil)

	if f.opts.WithCommit {
		cmt := Comments(
			`This is a generated source file. DO NOT EDIT`,
			`Source: `+path.Join(f.Pkg, path.Base(f.Name)),
		)
		if f.opts.WithToolVersion {
			cmt.Append(`Version: ` + Version)
		}
		if f.opts.WithTimestamp {
			cmt.Append(`Date: ` + time.Now().Format(time.Stamp))
		}

		buf.Write(cmt.Bytes())
		buf.WriteRune('\n')
	}

	buf.Write([]byte("\npackage " + stringsx.LowerSnakeCase(f.Pkg) + "\n"))

	if len(f.Imps) > 0 {
		if len(f.Imps) == 1 {
			buf.Write([]byte("import "))
		} else if len(f.Imps) > 1 {
			buf.Write([]byte("import (\n"))
		}

		for _, imp := range f.OrderedImps {
			if IsReserved(imp[0]) {
				panic("[CONFLICT] package name conflict reserved")
			}
			if imp[0] != path.Base(imp[1]) {
				buf.WriteString(imp[0])
				buf.WriteByte(' ')
			}
			buf.WriteByte('"')
			buf.WriteString(imp[1])
			buf.WriteByte('"')
			buf.WriteByte('\n')
		}

		if len(f.Imps) > 1 {
			buf.Write([]byte(")\n"))
		}
	}

	buf.Write(f.Buffer.Bytes())

	if f.opts.MustFormat {
		return formatx.MustFormat(f.Name, "", buf.Bytes(), formatx.SortImports)
	}
	return buf.Bytes()
}

func (f *File) Bytes() []byte {
	return f.bytes()
}

// Raw test only
func (f File) Raw() []byte { return f.bytes() }

// Formatted test only
func (f File) Formatted() []byte { return f.bytes() }

func (f *File) _import(pkg string) string {
	if f.Imps == nil {
		f.Imps = make(map[string]string)
		f.Pkgs = make(map[string][]string)
	}

	if _, ok := f.Imps[pkg]; !ok {
		pkgs, err := packages.Load(nil, pkg)
		if err != nil {
			panic(err)
		}
		if len(pkgs) == 0 {
			panic(pkg + " not found")
		}
		pkg = pkgs[0].PkgPath
		min := path.Base(pkg)

		if len(f.Pkgs[min]) == 0 {
			f.Imps[pkg] = min
		} else {
			f.Imps[pkg] = stringsx.LowerSnakeCase(
				fmt.Sprintf("gen %s %d", min, len(f.Pkgs[min])),
			)
		}
		f.Pkgs[min] = append(f.Pkgs[min], pkg)
		f.OrderedImps = append(f.OrderedImps, [2]string{f.Imps[pkg], pkg})
	}
	return f.Imps[pkg]
}

func (f *File) Use(pkg, name string) string { return f._import(pkg) + "." + name }

func (f *File) Expr(format string, args ...interface{}) SnippetExpr {
	return ExprWithAlias(f._import)(format, args...)
}

func (f *File) Type(t reflect.Type) SnippetType {
	return TypeWithAlias(f._import)(t)
}

func (f *File) Value(v interface{}) Snippet { return ValueWithAlias(f._import)(v) }

func (f *File) WriteSnippet(ss ...Snippet) {
	for _, s := range ss {
		if s != nil {
			f.Buffer.Write(s.Bytes())
			f.Buffer.WriteString("\n\n")
		}
	}
}

func (f *File) Write(opts ...WriterOptionSetter) (int, error) {
	for _, setter := range opts {
		setter(&f.opts)
	}

	if dir := filepath.Dir(f.Name); dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return -1, err
		}
	}

	fl, err := os.Create(f.Name)
	if err != nil {
		return -1, err
	}
	defer fl.Close()

	size, err := fl.Write(f.Bytes())
	if err != nil {
		return -1, err
	}

	if err := fl.Sync(); err != nil {
		return -1, err
	}
	return size, nil
}

type WriteOption struct {
	WithCommit      bool
	WithTimestamp   bool
	WithToolVersion bool
	MustFormat      bool
}

func WriteOptionWithCommit(v *WriteOption)      { v.WithCommit = true }
func WriteOptionWithTimestamp(v *WriteOption)   { v.WithCommit, v.WithTimestamp = true, true }
func WriteOptionWithToolVersion(v *WriteOption) { v.WithCommit, v.WithToolVersion = true, true }
func WriteOptionMustFormat(v *WriteOption)      { v.MustFormat = true }

type WriterOptionSetter func(v *WriteOption)
