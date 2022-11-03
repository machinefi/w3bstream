package ptrx_test

import (
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
)

func Example() {
	fmt.Println(*(ptrx.Int(10)))
	fmt.Println(*(ptrx.Float64(10)))
	fmt.Println(*(ptrx.String("abc")))

	// Output:
	// 10
	// 10
	// abc
}
