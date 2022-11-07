package transformer_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

type Sub struct {
	A string `name:"a" in:"query"`
}

type PtrSub struct {
	B []string `name:"b" in:"query"`
}

type P struct {
	Sub
	*PtrSub
	C *string `name:"c" in:"query"`
}

func TestParameters(t *testing.T) {
	params := make([]*Param, 0)

	p := P{}
	p.A = "a"
	p.PtrSub = &PtrSub{B: []string{"b"}}
	p.C = ptrx.String("c")

	EachParameter(bgctx, typesx.FromReflectType(reflect.TypeOf(p)),
		func(p *Param) bool {
			params = append(params, p)
			return true
		},
	)

	rv := reflect.ValueOf(&p)

	NewWithT(t).Expect(params).To(HaveLen(3))
	NewWithT(t).Expect(params[0].FieldValue(rv).Interface()).To(Equal(p.A))
	NewWithT(t).Expect(params[1].FieldValue(rv).Interface()).To(Equal(p.B))
	NewWithT(t).Expect(params[2].FieldValue(rv).Interface()).To(Equal(p.C))
}

func BenchmarkParameter_FieldValue(b *testing.B) {
	p := P{}
	p.A = "a"
	p.PtrSub = &PtrSub{
		B: []string{"b"},
	}
	p.C = ptrx.String("c")

	rv := reflect.ValueOf(&p).Elem()

	params := make([]*Param, 0)

	EachParameter(bgctx, typesx.FromReflectType(reflect.TypeOf(p)),
		func(p *Param) bool {
			params = append(params, p)
			return true
		},
	)

	b.Run("useCache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for i := range params {
				_ = params[i].FieldValue(rv).Addr()
			}
		}
	})

	b.Run("WalkDirect", func(b *testing.B) {
		var walk func(rv reflect.Value)

		walk = func(rv reflect.Value) {
			tpe := rv.Type()

			for i := 0; i < rv.NumField(); i++ {
				ft := tpe.Field(i)
				f := rv.Field(i)

				if ft.Anonymous && ft.Type.Kind() == reflect.Struct {
					walk(f)
				}
			}
		}

		for i := 0; i < b.N; i++ {
			walk(rv)
		}
	})
}
