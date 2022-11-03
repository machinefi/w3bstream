package validator_test

import (
	"context"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

var (
	rtInt                      = reflect.TypeOf(int(1))
	rtInt8                     = reflect.TypeOf(int8(1))
	rtInt16                    = reflect.TypeOf(int16(1))
	rtInt32                    = reflect.TypeOf(int32(1))
	rtInt64                    = reflect.TypeOf(int64(1))
	rtUint                     = reflect.TypeOf(uint(1))
	rtUint8                    = reflect.TypeOf(uint8(1))
	rtUint16                   = reflect.TypeOf(uint16(1))
	rtUint32                   = reflect.TypeOf(uint32(1))
	rtUint64                   = reflect.TypeOf(uint64(1))
	rtFloat32                  = reflect.TypeOf(float32(1))
	rtFloat64                  = reflect.TypeOf(float64(1))
	rtMapStringString          = reflect.TypeOf(map[string]string{})
	rtMapStringMapStringString = reflect.TypeOf(map[string]map[string]string{})
	rtString                   = reflect.TypeOf("")
	rtArrayString              = reflect.TypeOf([1]string{})
	rtSliceString              = reflect.TypeOf([]string{})
	rtSliceFloat64             = reflect.TypeOf([]float64{})

	rtSomeStruct     = reflect.TypeOf(&SomeStruct{})
	rttSomeStruct    = typesx.FromReflectType(rtSomeStruct.Elem())
	rttSomeStructPtr = typesx.FromReflectType(rtSomeStruct)

	rttInt                      = typesx.FromReflectType(rtInt)
	rttInt8                     = typesx.FromReflectType(rtInt8)
	rttInt16                    = typesx.FromReflectType(rtInt16)
	rttInt32                    = typesx.FromReflectType(rtInt32)
	rttInt64                    = typesx.FromReflectType(rtInt64)
	rttFloat32                  = typesx.FromReflectType(rtFloat32)
	rttFloat64                  = typesx.FromReflectType(rtFloat64)
	rttSliceFloat64             = typesx.FromReflectType(rtSliceFloat64)
	rttMapStringString          = typesx.FromReflectType(rtMapStringString)
	rttMapStringMapStringString = typesx.FromReflectType(rtMapStringMapStringString)
	rttString                   = typesx.FromReflectType(rtString)

	bg  = context.Background()
	ctx = validator.ContextWithFactory(bg, validator.DefaultFactory)
)

type StringType string

type Named string

type SubPtrStruct struct {
	PtrInt   *int     `validate:"@int[1,]"`
	PtrFloat *float32 `validate:"@float[1,]"`
	PtrUint  *uint    `validate:"@uint[1,]"`
}

type SubStruct struct {
	Int   int     `json:"int"   validate:"@int[1,]"`
	Float float32 `json:"float" validate:"@float[1,]"`
	Uint  uint    `json:"uint"  validate:"@uint[1,]"`
}

type SomeStruct struct {
	JustRequired string
	CanEmpty     *string              `validate:"@string[0,]?"`
	String       string               `validate:"@string[1,]"`
	Named        Named                `validate:"@string[2,]"`
	PtrString    *string              `validate:"@string[3,]" default:"123"`
	SomeStringer *SomeTextMarshaler   `validate:"@string[20,]"`
	Slice        []string             `validate:"@slice<@string[1,]>"`
	SliceStruct  []SubStruct          `validate:"@slice"`
	Map          map[string]string    `validate:"@map<@string[2,],@string[1,]>"`
	MapStruct    map[string]SubStruct `validate:"@map<@string[2,],>"`
	Struct       SubStruct
	SubStruct
	*SubPtrStruct
}

type SomeTextMarshaler struct {
}

func (*SomeTextMarshaler) MarshalText() ([]byte, error) {
	return []byte("SomeTextMarshaler"), nil
}
