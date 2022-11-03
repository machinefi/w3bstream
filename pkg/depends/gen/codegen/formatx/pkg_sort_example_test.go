package formatx_test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	. "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
)

func CreateDemoFile(filename string) *File {
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

func ExampleFormat() {
	cwd, _ := os.Getwd()
	filename := path.Join(cwd, "hello/hello.go")

	f := CreateDemoFile(filename)

	defer os.RemoveAll(filepath.Dir(f.Name))

	_, err := f.Write()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(f.Formatted()))

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
