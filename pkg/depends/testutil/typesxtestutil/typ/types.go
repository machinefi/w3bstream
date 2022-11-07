package typ

import (
	"context"
	"encoding"
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/testutil/typesxtestutil/typ/typ"
)

type String string
type Bool bool
type Int int
type Int8 int8
type Int16 int16
type Int32 int32
type Int64 int64
type Uint uint
type Uint8 uint8
type Uint16 uint16
type Uint32 uint32
type Uint64 uint64
type Uintptr uintptr
type Float32 float32
type Float64 float64
type Complex64 complex64
type Complex128 complex128

type Array [1]string

type Map map[string]string
type Slice []string
type Chan chan string
type Func func(a, b string) bool

func F() {}

type Struct struct {
	Interface
	a    string
	A    string `json:"a"`
	B    string `json:"b"`
	Bool `json:"bool,omitempty"`
	typ.Part
	Part2 Part `json:",omitempty"`
}

func (Struct) String() string {
	return ""
}

type Part struct {
	C string `json:"c"`
}

func (Part) Value() string {
	return ""
}

func (*Part) PtrValue() string {
	return ""
}

type DeepCompose struct {
	Struct
}

type Interface interface {
	String() string
}

type Enum int

const (
	ENUM__ONE Enum = iota + 1 // one
	ENUM__TWO                 // two
)

func (e *Enum) UnmarshalText(text []byte) error {
	switch string(text) {
	case "ONE":
		*e = ENUM__ONE
	case "TWO":
		*e = ENUM__TWO
	}
	return fmt.Errorf("unknown enum")
}

func (e Enum) MarshalText() ([]byte, error) {
	switch e {
	case ENUM__ONE:
		return []byte("ONE"), nil
	case ENUM__TWO:
		return []byte("TWO"), nil
	}
	return []byte{}, fmt.Errorf("unknown enum")
}

type SomeMixInterface interface {
	encoding.TextMarshaler
	Stringify(ctx context.Context, vs ...any) string
	Add(a, b string) string
	Bytes() []byte
	s() string
}

type AnySlice[V any] []V

func (m AnySlice[V]) Each() {

}

type AnyMap[K comparable, V any] map[K]V

type IntMap = AnyMap[Enum, any]

type AnyStruct[V any] struct {
	Struct
	Name V
}
