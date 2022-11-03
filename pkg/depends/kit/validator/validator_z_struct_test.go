package validator_test

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/validator"
	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func TestStruct_New(t *testing.T) {
	v := NewStructValidator("json")
	_, err := v.New(ctx, &Rule{Type: rttSomeStruct})

	NewWithT(t).Expect(err).To(BeNil())
}

func TestStruct_NewFailed(t *testing.T) {
	type Named string

	type Struct struct {
		Int    int  `validate:"@int[1,"`
		PtrInt *int `validate:"@uint[2,"`
	}

	type SubStruct struct {
		Float    float32  `validate:"@string"`
		PtrFloat *float32 `validate:"@unknown"`
	}

	type SomeStruct struct {
		String string   `validate:"@string[1,"`
		Named  Named    `validate:"@int"`
		Slice  []string `validate:"@slice<@int>"`
		SubStruct
		Struct Struct
	}

	v := NewStructValidator("json")
	rt := typesx.FromReflectType(reflect.TypeOf(&SomeStruct{}).Elem())

	_, err := v.New(ctx, &Rule{Type: rt})
	NewWithT(t).Expect(err).NotTo(BeNil())
	t.Logf("\n%v", err)

	_, err = v.New(ctx, &Rule{Type: rttString})
	NewWithT(t).Expect(err).NotTo(BeNil())
	t.Logf("\n%v", err)
}

func ExampleNewStructValidator() {
	v := NewStructValidator("json")

	sv, err := v.New(ctx, &Rule{
		Type: typesx.FromReflectType(reflect.TypeOf(&SomeStruct{}).Elem()),
	})
	if err != nil {
		return
	}

	s := SomeStruct{
		Slice:       []string{"", ""},
		SliceStruct: []SubStruct{{Int: 0}},
		Map:         map[string]string{"1": "", "11": "", "12": ""},
		MapStruct:   map[string]SubStruct{"222": SubStruct{}},
	}

	err = sv.Validate(s)
	var (
		errs     = map[string]string{}
		keyPaths = make([]string, 0)
		flatten  = err.(*errors.ErrorSet).Flatten()
	)

	flatten.Each(
		func(ferr *errors.FieldError) {
			errs[ferr.Field.String()] = strconv.Quote(ferr.Error.Error())
			keyPaths = append(keyPaths, ferr.Field.String())
		},
	)

	sort.Strings(keyPaths)

	for i := range keyPaths {
		k := keyPaths[i]
		fmt.Println(k, errs[k])
	}

	// Output:
	// JustRequired "missing required field"
	// Map.1 "missing required field"
	// Map.1/key "string length should be larger than 2, but got invalid value 1"
	// Map.11 "missing required field"
	// Map.12 "missing required field"
	// MapStruct.222.float "missing required field"
	// MapStruct.222.int "missing required field"
	// MapStruct.222.uint "missing required field"
	// Named "missing required field"
	// PtrFloat "missing required field"
	// PtrInt "missing required field"
	// PtrString "missing required field"
	// PtrUint "missing required field"
	// SliceStruct[0].float "missing required field"
	// SliceStruct[0].int "missing required field"
	// SliceStruct[0].uint "missing required field"
	// Slice[0] "missing required field"
	// Slice[1] "missing required field"
	// SomeStringer "missing required field"
	// String "missing required field"
	// Struct.float "missing required field"
	// Struct.int "missing required field"
	// Struct.uint "missing required field"
	// float "missing required field"
	// int "missing required field"
	// uint "missing required field"
}
