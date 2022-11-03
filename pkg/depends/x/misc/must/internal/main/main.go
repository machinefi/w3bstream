package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	. "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

func main() {
	pkg := "must"
	_, path, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(path), "../..")

	if name := filepath.Base(root); name != pkg {
		log.Panicf("wrong execute location: \n\tpath: %s\n\tbase: %s", root, name)
	}

	{
		filename := filepath.Join(root, "must.go")
		_, err := os.Stat(filename)
		if err != nil && !os.IsExist(err) {
			file := NewFile(pkg, filename)
			file.WriteSnippet(
				Comments("pls add your assert function here, or add type and re-generate"),
			)
			if _, err := file.Write(); err != nil {
				log.Panic(err)
			}
		}
	}

	types := []struct {
		BuiltInType
		Name string
	}{
		{BuiltInType: "byte"},
		{BuiltInType: "[]byte", Name: "bytes"},
		{BuiltInType: "string"},
		{BuiltInType: "[]string", Name: "strings"},
		{BuiltInType: "int"},
		{BuiltInType: "int8"},
		{BuiltInType: "int16"},
		{BuiltInType: "int32"},
		{BuiltInType: "int64"},
		{BuiltInType: "uint8"},
		{BuiltInType: "uint16"},
		{BuiltInType: "uint32"},
		{BuiltInType: "uint64"},
		{BuiltInType: "rune"},
		{BuiltInType: "float32"},
		{BuiltInType: "float64"},
	}

	{
		filename := filepath.Join(root, "must_generated.go")
		file := NewFile(pkg, filename)
		for _, t := range types {
			name := ""
			if t.Name != "" {
				name += stringsx.UpperCamelCase(t.Name)
			} else {
				name += stringsx.UpperCamelCase(string(t.BuiltInType))
			}
			file.WriteSnippet(
				Func(Var(t.BuiltInType, "v"), Var(Error, "err")).
					Named(name).
					Return(Var(t.BuiltInType)).
					Do(
						If(Exprer("err != nil")).
							Do(
								Call(file.Use("log", "Panic"), Ident("err")),
							),
						Return(Ident("v")),
					),
				Func(Var(t.BuiltInType, "v"), Var(Bool, "ok")).
					Named(name+"OK").
					Return(Var(t.BuiltInType)).
					Do(
						If(Exprer("!ok")).
							Do(
								Call(file.Use("log", "Panic"), Valuer(name+" not ok")),
							),
						Return(Ident("v")),
					),
			)
		}

		if _, err := file.Write(); err != nil {
			log.Panic(err)
		}
	}
}
