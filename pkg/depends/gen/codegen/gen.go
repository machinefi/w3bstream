package codegen

import (
	"bytes"
	"fmt"
	"io"
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
	Imps        map[string]string   // key: package path ; val: package name
	Pkgs        map[string][]string // key: package name ; val: package paths
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

func (f *File) Import(pkg string) string {
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
		name := pkgs[0].Name

		if name == "" {
			name = filepath.Base(pkgs[0].PkgPath)
			// panic("cannot load package name: " + pkg)
		}

		if len(f.Pkgs[name]) == 0 {
			f.Imps[pkg] = name
		} else {
			f.Imps[pkg] = fmt.Sprintf("%s%d", name, len(f.Pkgs[name])+1)
		}
		f.Pkgs[name] = append(f.Pkgs[name], pkg)
		f.OrderedImps = append(f.OrderedImps, [2]string{f.Imps[pkg], pkg})
	}
	return f.Imps[pkg]
}

func (f *File) Use(pkg, name string) string {
	if name == "" {
		panic("should give using name")
	}
	if pkg == "" {
		return name
	}
	return f.Import(pkg) + "." + name
}

func (f *File) Expr(format string, args ...interface{}) SnippetExpr {
	return ExprWithAlias(f.Import)(format, args...)
}

func (f *File) Type(t reflect.Type) SnippetType {
	return TypeWithAlias(f.Import)(t)
}

func (f *File) Value(v interface{}) Snippet { return ValueWithAlias(f.Import)(v) }

func (f *File) WriteSnippet(ss ...Snippet) {
	for _, s := range ss {
		if s != nil {
			f.Buffer.Write(s.Bytes())
			f.Buffer.WriteString("\n\n")
		}
	}
}

func (f *File) Write(opts ...WriteOptionSetter) (size int, err error) {
	for _, setter := range opts {
		setter(&f.opts)
	}

	if dir := filepath.Dir(f.Name); dir != "" {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return -1, err
		}
	}

	var wr io.Writer

	defer func() {
		if closer, ok := wr.(io.Closer); ok {
			closer.Close()
		}
	}()

	if f.opts.Output != nil {
		wr = f.opts.Output
	} else {
		var file *os.File
		file, err = os.Create(f.Name)
		if err != nil {
			return -1, err
		}
		wr = file
		defer func() {
			if err = file.Sync(); err != nil {
				size = -1
			}

		}()
	}

	size, err = wr.Write(f.Bytes())
	return
}

type WriteOption struct {
	WithCommit      bool
	WithTimestamp   bool
	WithToolVersion bool
	MustFormat      bool
	Output          io.Writer
}

type WriteOptionSetter func(v *WriteOption)

func WriteOptionWithCommit(v bool) WriteOptionSetter {
	return func(o *WriteOption) { o.WithCommit = v }
}

func WriteOptionWithTimestamp(v bool) WriteOptionSetter {
	return func(o *WriteOption) { o.WithTimestamp = v }
}

func WriteOptionWithToolVersion(v bool) WriteOptionSetter {
	return func(o *WriteOption) { o.WithToolVersion = v }
}

func WriteOptionMustFormat(v bool) WriteOptionSetter {
	return func(o *WriteOption) { o.MustFormat = v }
}
func WriteOptionWithOutput(v io.Writer) WriteOptionSetter {
	return func(o *WriteOption) { o.Output = v }
}
