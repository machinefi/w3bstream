package codegen_test

import (
	"fmt"
	"go/token"
	"reflect"
	"runtime"

	. "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
)

func ExampleVar() {
	var v = Var(String, "a")

	fmt.Println(string(v.Bytes()))
	fmt.Println(string(v.WithComments("A is a string").Bytes()))

	v = v.WithTag(`json:"a"`)
	fmt.Println(string(v.Bytes()))

	v = v.WithoutTag()
	fmt.Println(string(v.Bytes()))

	v = Var(String, "AliasStringType").AsAlias().
		WithComments("AliasStringType is a Alias of string")
	fmt.Println(string(v.Bytes()))

	// Output:
	// a string
	// // A is a string
	// a string
	// a string `json:"a"`
	// a string
	// // AliasStringType is a Alias of string
	// AliasStringType = string
}

func ExampleType() {

	var t = Type("Name")

	fmt.Println(string(t.Bytes()))
	fmt.Println(string(Chan(t).Bytes()))
	fmt.Println(string(Chan(Star(t)).Bytes()))
	fmt.Println(string(Array(t, 10).Bytes()))
	fmt.Println(string(Array(Star(t), 10).Bytes()))
	fmt.Println(string(Slice(t).Bytes()))
	fmt.Println(string(Map(Bool, Star(t)).Bytes()))

	// Output:
	// Name
	// chan Name
	// chan *Name
	// [10]Name
	// [10]*Name
	// []Name
	// map[bool]*Name
}

func ExampleTypeAssert() {

	ta := TypeAssert(String, Ident("value"))
	fmt.Println(string(ta.Bytes()))

	ta = TypeAssert(Type("Name"), Ident("value"))
	fmt.Println(string(ta.Bytes()))

	ta = TypeAssert(Star(Type("Name")), Ident("value"))
	fmt.Println(string(ta.Bytes()))

	// Output:
	// value.(string)
	// value.(Name)
	// value.(*Name)
}

func ExampleIdent() {
	ids := Idents("a", "b", "C", "Type.Field")
	for _, id := range ids {
		fmt.Println(string(id.Bytes()))
	}

	// Output:
	// a
	// b
	// C
	// Type.Field
}

func ExampleDefine() {
	fmt.Println(Stringify(
		Define(
			Ident("a"),
			Ident("b"),
		).By(
			Ident("c"),
			Call("Fn", Valuer(1)),
		)),
	)

	fmt.Println(Stringify(
		Define(
			Ident("a"),
			Ident("b"),
		).By(
			Ident("c"),
			Call("Fn", Valuer("abc")),
		)),
	)

	fmt.Println(Stringify(
		DeclVar(
			Var(Int, "a"),
			Assign(Var(String, "a", "b")).By(Call("Fn", Valuer(1))),
		),
	))

	fmt.Println(Stringify(
		DeclConst(
			Assign(Var(Int, "a")).By(Iota),
			Assign(Ident("b")),
			Assign(Ident("c")),
		),
	))

	fmt.Println(Stringify(
		DeclType(
			Var(Type("time.Time"), "aliasTime").AsAlias(),
		),
	))

	// Output:
	// a, b := c, Fn(1)
	// a, b := c, Fn("abc")
	// var (
	// a int
	// a, b string = Fn(1)
	// )
	// const (
	// a int = iota
	// b
	// c
	// )
	// type aliasTime = time.Time
}

func ExampleFor() {
	i := Ident("i")
	fmt.Println(
		Stringify(
			For(
				Define(i).By(Valuer(0)),
				AssignWith(token.LSS, i).By(Valuer(10)),
				Inc(i),
			).Do(),
		),
	)

	// Output:
	// for i := 0; i < 10; i++ {
	// }
}

func ExampleForRange() {
	k, v, ranger := Ident("k"), Ident("v"), Ident("ranger")
	fmt.Println(
		Stringify(
			ForRange(ranger, k, v).Do(),
		),
	)

	// Output:
	// for k, v := range ranger {
	// }
}

var clauses []*SnippetCaseClause

func ExampleCaseClause() {
	clauses = append(clauses,
		CaseClause(Valuer(1), Valuer(2)),
		CaseClause(Arrow(Ref(Ident("time"), Call("After")))),
		CaseClause().Do(Return(Ident("a"), Call("Fn", Valuer(1)))),
	)

	for _, c := range clauses {
		fmt.Println(Stringify(c))
	}

	// Output:
	// case 1, 2:
	// case <-time.After():
	// default:
	// return a, Fn(1)
}

func ExampleSwitch() {
	fmt.Println(Stringify(Switch(nil).When(clauses...)))

	os := Ident("os")
	tar := Ident("tar")
	_ = runtime.GOOS
	fmt.Print(Stringify(
		Switch(os).
			InitWith(Define(os).By(Ref(Ident("runtime"), Ident("GOOS")))).
			When(
				CaseClause(Valuer("darwin")).Do(
					Assign(tar).By(Ident("OS_DARWIN")),
				),
				CaseClause(Valuer("linux")).Do(
					Assign(tar).By(Ident("OS_LINUX")),
				),
				CaseClause(Valuer("windows")).Do(
					Assign(tar).By(Ident("OS_WINDOWS")),
				),
				CaseClause().Do(
					Assign(tar).By(Ident("OS_UNKNOWN")),
				),
				CaseClause(Valuer("doNothing")),
			),
	))

	// Output:
	// switch {
	// case 1, 2:
	// case <-time.After():
	// default:
	// return a, Fn(1)
	// }
	// switch os := runtime.GOOS; os {
	// case "darwin":
	// tar = OS_DARWIN
	// case "linux":
	// tar = OS_LINUX
	// case "windows":
	// tar = OS_WINDOWS
	// default:
	// tar = OS_UNKNOWN
	// case "doNothing":
	// }
}

func ExampleSelect() {
	ch1, ch2 := "ch1", "ch2"

	fmt.Println(Stringify(DeclVar(
		Assign(Var(nil, ch1)).By(CallMakeChan(Int, 10)),
		Assign(Var(nil, ch2)).By(CallMakeChan(String, 10)),
		Var(String, "ret"),
	)))
	fmt.Println(Stringify(
		Select(
			CaseClause(Define(Ident("i")).By(Arrow(Ident(ch1)))).Do(
				Assign(Ident("ret")).By(
					Ref(Ident("fmt"), Call("Sprintf", Ident("i")))),
			),
			CaseClause(Define(Ident("s")).By(Arrow(Ident(ch2)))).Do(
				Assign(Ident("ret")).By(Ident("s")),
			),
			CaseClause().Do(),
		),
	))
	fmt.Println(Stringify(Return(Ident("ret"))))

	// Output:
	// var (
	// ch1 = make(chan int, 10)
	// ch2 = make(chan string, 10)
	// ret string
	// )
	// select {
	// case i := <-ch1:
	// ret = fmt.Sprintf(i)
	// case s := <-ch2:
	// ret = s
	// default:
	// }
	// return ret

}

func ExampleStar() {
	fmt.Println(Stringify(Star(Int)))

	// Output:
	// *int
}

func ExampleAddr() {
	fmt.Println(Stringify(Addr(Ident("i"))))

	// Output:
	// &i
}

func ExampleParen() {
	fmt.Println(Stringify(Paren(Valuer(1))))

	// Output:
	// (1)
}

func ExampleArrow() {
	fmt.Println(Stringify(Arrow(Ident("ch"))))
	fmt.Println(Stringify(Define(Ident("val")).By(Arrow(Ident("ch")))))

	// Output:
	// <-ch
	// val := <-ch
}

func ExampleCasting() {
	fmt.Println(Stringify(DeclVar(
		Assign(Ident("a")).By(Casting(Float64, Valuer(0.1))),
		Assign(Ident("b")).By(Casting(Float32, Ident("a"))),
	)))

	// Output:
	// var (
	// a = float64(0.1)
	// b = float32(a)
	// )
}

func ExampleCall() {
	fmt.Println(Stringify(Call("NewBuffer", Nil)))

	// Output:
	// NewBuffer(nil)
}

func ExampleCallWith() {
	fmt.Println(Stringify(Ref(Ident("bytes"), CallWith(Ident("NewBuffer"), Nil))))

	// Output:
	// bytes.NewBuffer(nil)
}

func ExampleComments() {
	fmt.Println(Stringify(DeclVar(
		Var(Int, "a").WithOneLineComment("this is a int var"),
	)))

	// Output:
	// var a int // this is a int var
}

func ExampleKeyValue() {
	fmt.Println(string(KeyValue(Valuer("key"), Valuer("value")).Bytes()))

	// Output:
	// "key": "value"
}

func ExampleInc() {
	fmt.Println(string(Inc(Ident("i")).Bytes()))

	// Output:
	// i++
}

func ExampleDec() {
	fmt.Println(string(Dec(Ident("i")).Bytes()))

	// Output:
	// i--
}

func ExampleAccess() {
	output(Access(Ident("array"), 10))
	output(AccessWith(Ident("array"), Exprer("?-?", Call("len", Ident("array")), Valuer(1))))
	output(AccessWith(Ident("array"), Exprer("len(array)-1")))

	// Output:
	// array[10]
	// array[len(array)-1]
	// array[len(array)-1]
}

type Anonymous interface{}

type Exported struct {
	F_int        int
	f_unexported interface{}
	F_map        map[string]struct{}
	F_slice      []int
	Anonymous
}

func ExampleTyper() {
	t := reflect.TypeOf(&Exported{})
	output(Typer(t))
	t = reflect.TypeOf(Snippet(SnippetBuiltIn("")))
	output(Typer(t))

	// Output:
	// *codegen_test.Exported
	// codegen.SnippetBuiltIn
}

func ExampleValuer() {
	type unexported Exported
	output(Valuer(1))
	newline()
	output(Valuer(unexported{
		F_int:        1,
		f_unexported: 2,
		F_map:        map[string]struct{}{"a": {}},
		F_slice:      []int{1, 2, 3},
		Anonymous:    3,
	}))
	newline()
	output(Valuer(&Exported{
		F_map:     map[string]struct{}{"a": {}},
		F_slice:   []int{1, 2, 3},
		Anonymous: 3,
	}))

	// Output:
	// 1
	//
	// codegen_test.unexported{
	// F_int: 1,
	// F_map: map[string]struct {
	// }{
	// "a": struct {
	// }{
	// },
	// },
	// F_slice: []int{
	// 1,
	// 2,
	// 3,
	// },
	// Anonymous: 3,
	// }
	//
	// &(codegen_test.Exported{
	// F_map: map[string]struct {
	// }{
	// "a": struct {
	// }{
	// },
	// },
	// F_slice: []int{
	// 1,
	// 2,
	// 3,
	// },
	// Anonymous: 3,
	// })
}

func ExampleNewFile() {
	filename := "examples/hello/hello.go"
	f := NewFile("main", filename)

	f.WriteSnippet(
		DeclVar(Var(Slice(String), "lines")),

		Func().Named("main").Do(
			DeclVar(
				Assign(Var(nil, "ch")).By(Call("make", String, Valuer(10))),
			),

			Call("close", Ident("ch")).AsDefer(),

			Call("PipeReadAndPrint", Ident("ch")).AsRoutine(),

			Call("PipePrintAndWrite", Ident("ch"), Ident("lines")),
		),

		Func(Var(ChanRO(String), "ch"), Var(Ellipsis(String), "v")).
			Named("PipePrintAndWrite").
			Do(
				ForRange(Ident("v"), Ident("i"), Ident("_")).Do(
					Call(f.Use("fmt", "Println"), AccessWith(Ident("v"), Ident("i"))),
					Exprer("? <- ?", Ident("ch"), AccessWith(Ident("v"), Ident("i"))),
				),
			),

		Func(Var(ChanRO(String), "ch")).Named("PipeReadAndPrint").Do(
			For(nil, nil, nil).Do(
				Select(
					CaseClause(Define(Ident("s"), Ident("ok")).By(Arrow(Ident("ch")))).Do(
						If(Ident("ok")).Do(
							Call(f.Use("fmt", "Println"), Ident("s")),
						).Else(
							If(Exprer("!ok")).Do(
								Return(),
							),
						),
					),
				),
			),
		),
	)

	fmt.Println(string(f.Formatted()))

	// Output:
	// // This is a generated source file. DO NOT EDIT
	// // Source: main/hello.go
	//
	// package main
	//
	// import "fmt"
	//
	// var lines []string
	//
	// func main() {
	// 	var ch = make(string, 10)
	// 	defer close(ch)
	// 	go PipeReadAndPrint(ch)
	// 	PipePrintAndWrite(ch, lines)
	// }
	//
	// func PipePrintAndWrite(ch <-chan string, v ...string) {
	// 	for i := range v {
	// 		fmt.Println(v[i])
	// 		ch <- v[i]
	// 	}
	// }
	//
	// func PipeReadAndPrint(ch <-chan string) {
	// 	for {
	// 		select {
	// 		case s, ok := <-ch:
	// 			if ok {
	// 				fmt.Println(s)
	// 			} else if !ok {
	// 				return
	// 			}
	// 		}
	// 	}
	// }
}

func output(s Snippet) { fmt.Println(string(s.Bytes())) }
func newline()         { fmt.Println() }
