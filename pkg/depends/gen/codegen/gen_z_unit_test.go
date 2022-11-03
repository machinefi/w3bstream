package codegen_test

import (
	"fmt"
	"os"

	. "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
)

func CreateDemoFile() *File {
	filename := "examples/hello/hello.go"
	f := NewFile("main", filename)
	f.WriteSnippet(Func().Named("main").Do(
		Call(f.Use("fmt", "Println"), f.Value("Hello, 世界")),
		Call(f.Use("github.com/some/pkg", "Println"), f.Value("Hello World!")),
		Call(f.Use("github.com/another/pkg", "Println"), f.Value("Hello World!")),
		Call(f.Use("github.com/one_more/pkg", "Println"), f.Value("Hello World!")),

		Assign(AnonymousIdent).By(Call(f.Use("bytes", "NewBuffer"), f.Value(nil))),
	))

	return f
}

func ExampleNewFileFullFormat() {
	f := CreateDemoFile()

	defer os.RemoveAll("examples")

	if _, err := f.Write(); err != nil {
		panic(err)
	}
	if raw, err := os.ReadFile(f.Name); err != nil {
		panic(err)
	} else {
		// NOTE: this test should always FAILED because the generated file
		// contains `time` and `version` information
		fmt.Println(string(raw))
	}
	// Output:
	// // This is a generated source file. DO NOT EDIT
	// // Source: main/hello.go
	//
	// package main
	//
	// import (
	// 	"bytes"
	// 	"fmt"
	//
	// 	gen_pkg_1 "github.com/another/pkg"
	// 	gen_pkg_2 "github.com/one_more/pkg"
	// 	"github.com/some/pkg"
	// )
	//
	// func main() {
	// 	fmt.Println("Hello, 世界")
	// 	pkg.Println("Hello World!")
	// 	gen_pkg_1.Println("Hello World!")
	// 	gen_pkg_2.Println("Hello World!")
	// 	_ = bytes.NewBuffer(nil)
	// }
}
