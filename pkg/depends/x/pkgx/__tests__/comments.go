// package
package main

import (
	"fmt"
	"time"
)

// type Date
type Date time.Time

const (
	// a
	A = "A" // A
	// b
	B = "B" // B
	// c
	C = "C" // C
)

// type Test
type Test struct {
	// field String
	String string
	// field Int
	Int int
	// field Bool
	Bool bool
	// field Date
	Date Date
}

// method
func (Test) Recv() {

}

// type Test2
type (
	Test2 struct {
		// field String
		String string
		// field Int
		Int int
		// field Bool
		Bool bool
	}
)

// var
var test = Test{
	String: "",
	Int:    1 + 1,
	Bool:   true,
}

// var
var (
	// test2
	test2 = Test{
		String: "",
		Int:    1,
		Bool:   true,
	}
	// test3
	test3 = Test{
		String: "",
		Int:    1,
		Bool:   true,
	}
)

// func Print
//
//go:generate echo
func Print(a string, b string) string {
	return a + b
}

// func fn
func fn() {
	// Call
	res := Print("", "")
	if res != "" {
		// print
		fmt.Println(res)
	}
}
