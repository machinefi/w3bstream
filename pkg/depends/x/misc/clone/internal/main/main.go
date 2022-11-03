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
	pkg := "clone"
	_, path, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(path), "../..")

	if name := filepath.Base(root); name != pkg {
		log.Panicf("wrong execute location: \n\tpath: %s\n\tbase: %s", root, name)
	}

	{
		filename := filepath.Join(root, pkg+".go")
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

	types := []BuiltInType{
		Byte,
		String,
		Int,
		Int8,
		Int16,
		Int32,
		Int64,
		Uint,
		Uint8,
		Uint16,
		Uint32,
		Uint64,
		Rune,
		Float32,
		Float64,
	}

	{
		filename := filepath.Join(root, pkg+"_generated.go")
		file := NewFile(pkg, filename)
		for _, t := range types {
			fn := stringsx.UpperCamelCase(string(t)) + "s"
			st := Slice(t)
			file.WriteSnippet(
				Func(Var(st, "orig")).
					Named(fn).
					Return(Var(st)).
					Do(
						Define(Ident("cloned")).
							By(
								Call("make", st, Call("len", Ident("orig"))),
							),
						Call("copy", Ident("cloned"), Ident("orig")),
						Return(Ident("cloned")),
					),
			)
		}

		if _, err := file.Write(); err != nil {
			log.Panic(err)
		}
	}
}
