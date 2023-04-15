package codegen_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/gomega"

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
		// NOTE: this test should always FAILED if the generated file
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
	// 	pkg2 "github.com/another/pkg"
	// 	pkg3 "github.com/one_more/pkg"
	// 	"github.com/some/pkg"
	// )
	//
	// func main() {
	// 	fmt.Println("Hello, 世界")
	// 	pkg.Println("Hello World!")
	// 	pkg2.Println("Hello World!")
	// 	pkg3.Println("Hello World!")
	// 	_ = bytes.NewBuffer(nil)
	// }
}

func DISABLE_TestFile_Import(t *testing.T) {
	f := NewFile("fake", "fake")

	// remote
	path := "github.com/golang-jwt/jwt/v4"
	f.Import(path)
	NewWithT(t).Expect(f.Imps[path]).To(Equal("jwt"))

	// renamed
	path = "github.com/golang-jwt/jwt/v5"
	f.Import(path)
	NewWithT(t).Expect(f.Imps[path]).To(Equal("jwt2"))

	buf := bytes.NewBuffer(nil)
	_, err := f.Write(
		WriteOptionWithOutput(buf),
		WriteOptionMustFormat(false),
		WriteOptionWithCommit(false),
	)
	fmt.Println(string(buf.Bytes()))
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(buf.Bytes()).To(Equal(
		[]byte(`
package fake
import (
jwt "github.com/golang-jwt/jwt/v4"
jwt2 "github.com/golang-jwt/jwt/v5"
)
`)))
}
