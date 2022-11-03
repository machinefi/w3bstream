package textx_test

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	. "github.com/machinefi/w3bstream/pkg/depends/x/textx"
)

type Duration time.Duration

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

func (d *Duration) UnmarshalText(data []byte) error {
	dur, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

type NamedString string
type NamedInt int

var (
	v = struct {
		NamedString  NamedString
		NamedInt     NamedInt
		Duration     Duration
		PtrDuration  *Duration
		String       string
		PtrString    *string
		PtrPtrString **string
		Int          int
		PtrInt       *int
		Uint         uint
		PtrUint      *uint
		Float        float32
		PtrFloat     *float32
		Bool         bool
		PtrBool      *bool
		Bytes        []byte
	}{}

	rv = reflect.ValueOf(&v).Elem()
	d  = Duration(2 * time.Second)
)

var cases = []struct {
	name   string
	v      interface{}
	text   string
	expect interface{}
}{
	{
		"Ptr String",
		rv.FieldByName("PtrString"),
		"string",
		ptrx.String("string"),
	},
	{
		"Ptr Ptr String",
		rv.FieldByName("PtrPtrString"),
		"string",
		func() **string {
			s := ptrx.String("string")
			return &s
		}(),
	},
	{
		"Ptr String raw value",
		&v.String,
		"ptr",
		ptrx.String("ptr"),
	},
	{
		"Named String",
		rv.FieldByName("NamedString"),
		"string",
		NamedString("string"),
	},
	{
		"Duration",
		rv.FieldByName("Duration"),
		"2s",
		Duration(2 * time.Second),
	},
	{
		"Ptr Duration",
		rv.FieldByName("PtrDuration"),
		"2s",
		&d,
	},
	{
		"Int",
		rv.FieldByName("Int"),
		"1",
		1,
	},
	{
		"Named Int",
		rv.FieldByName("NamedInt"),
		"11",
		NamedInt(11),
	},
	{
		"PtrInt",
		rv.FieldByName("PtrInt"),
		"1",
		ptrx.Int(1),
	},
	{
		"Uint",
		rv.FieldByName("Uint"),
		"1",
		uint(1),
	},
	{
		"Int raw value",
		rv.FieldByName("Int").Addr().Interface(),
		"1",
		ptrx.Int(1),
	},
	{
		"PtrUint",
		rv.FieldByName("PtrUint"),
		"1",
		ptrx.Uint(1),
	},
	{
		"Float",
		rv.FieldByName("Float"),
		"1",
		float32(1),
	},
	{
		"PtrFloat",
		rv.FieldByName("PtrFloat"),
		"1.1",
		ptrx.Float32(1.1),
	},
	{
		"Bool",
		rv.FieldByName("Bool"),
		"true",
		true,
	},
	{
		"PtrBool",
		rv.FieldByName("PtrBool"),
		"true",
		ptrx.Bool(true),
	},
	{
		"Bytes",
		rv.FieldByName("Bytes"),
		"111",
		[]byte("111"),
	},
	{
		"Bytes direct",
		&v.Bytes,
		"111",
		func() *[]byte {
			b := []byte("111")
			return &b
		}(),
	},
}

func BenchmarkPtrFloat(b *testing.B) {
	v.PtrFloat = ptrx.Float32(1.1)
	//rv := reflect.ValueOf(v.PtrFloat).Elem()

	b.Run("append", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			//f := rv.Float()
			//_, _ = MarshalText(v.PtrFloat)
			d := make([]byte, 0)
			strconv.AppendFloat(d, float64(*v.PtrFloat), 'f', -1, 32)
		}

		//fmt.Println(string(d))
	})

	b.Run("format", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			//f := rv.Float()
			//_, _ = MarshalText(v.PtrFloat)
			_ = []byte(strconv.FormatFloat(float64(*v.PtrFloat), 'f', -1, 32))
		}
		//fmt.Println(string(d))
	})
}

func BenchmarkUnmarshalTextAndMarshalText(b *testing.B) {
	for i := range cases {
		c := cases[i]

		b.Run(fmt.Sprintf("UnmarshalText %s", c.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = UnmarshalText(c.v, []byte(c.text))
			}
		})

		b.Run(fmt.Sprintf("MarshalText %s", c.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = MarshalText(c.v)
			}
		})
	}
}

func TestUnmarshalTextAndMarshalText(t *testing.T) {
	for _, c := range cases {
		t.Run(fmt.Sprintf("UnmarshalText %s", c.name), func(t *testing.T) {
			err := UnmarshalText(c.v, []byte(c.text))

			NewWithT(t).Expect(err).To(BeNil())

			if rv, ok := c.v.(reflect.Value); ok {
				NewWithT(t).Expect(c.expect).To(Equal(rv.Interface()))
			} else {
				NewWithT(t).Expect(c.expect).To(Equal(c.v))
			}
		})
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("MarshalText %s", c.name), func(t *testing.T) {
			text, err := MarshalText(c.v)
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(c.text).To(Equal(string(text)))
		})
	}

	v2 := struct {
		PtrString *string
		Slice     []string
	}{}

	rv2 := reflect.ValueOf(v2)

	{
		_, err := MarshalText(rv2.FieldByName("Slice"))
		NewWithT(t).Expect(err).NotTo(BeNil())
	}

	{
		_, err := MarshalText(rv2.FieldByName("PtrString"))
		NewWithT(t).Expect(err).To(BeNil())
	}
}
