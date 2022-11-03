package main

import (
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx/__tests__/sub"
)

func Example() {

}

type String string

func (String) Method() string {
	return strings.Join(strings.Split("1", ","), ",")
}

func FuncSingleReturn() interface{} {
	// should skip
	_ = func() string {
		return ""
	}()

	var a interface{}
	a = "" + "1"
	a = 2

	return a
}

func FuncSelectExprReturn() string {
	s := "2"
	return (struct{ s string }{s: s}).s
}

func FuncWillCall() (a interface{}, s String) {
	return FuncSingleReturn(), String(FuncSelectExprReturn())
}

func FuncReturnWithCallDirectly() (a interface{}, b String) {
	return FuncWillCall()
}

func FuncWithNamedReturn() (a interface{}, b String) {
	a, b = FuncWillCall()
	return
}

func FuncSingleNamedReturnByAssign() (a interface{}, s String) {
	a = "" + "1"
	s = "2"
	return
}

func FuncWithSwitch() (a interface{}, b String) {
	switch a {
	case "1":
		a = "a1"
		b = "b1"
		return
	case "2":
		a = "a2"
		b = "b2"
		return
	default:
		a = "a3"
		b = "b3"
	}
	return
}

func str(a string, b string) string {
	return a + b
}

func FuncWithIf() (a interface{}) {
	if true {
		a = "a0"
		return
	} else if true {
		a = "a1"
		return
	} else {
		a = str("a", "b")
		return
	}
}

func FuncCallReturnAssign() (a interface{}, b String) {
	return FuncSingleReturn(), String(FuncSelectExprReturn())
}

func FuncCallWithFuncLit() (a interface{}, b String) {
	call := func() interface{} {
		return 1
	}

	return call(), "1"
}

func FuncWithCurryCall() interface{} {
	return sub.CurryCall()()()
}
