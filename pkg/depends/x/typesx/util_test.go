package typesx_test

import (
	"encoding"
	"go/types"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"unsafe"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/testutil/typesxtestutil/typ"
	typ2 "github.com/machinefi/w3bstream/pkg/depends/testutil/typesxtestutil/typ/typ"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	. "github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

var pkgid string

func init() {
	_, current, _, _ := runtime.Caller(0)
	pkgid = must.String(pkgx.PkgIdByPath(filepath.Dir(current)))
}

func TestTypeFor(t *testing.T) {
	cases := []struct {
		name, id string
	}{
		{"string", "string"},
		{"int", "int"},
		{"map[int]int", "map[int]int"},
		{"[]int", "[]int"},
		{"[2]int", "[2]int"},
		{"error", "error"},
		{"typesx.GoType", pkgid + ".GoType"},
		{"typesx.ReflectType", pkgid + ".ReflectType"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gt := FromGoType(TypeFor(c.id))
			NewWithT(t).Expect(gt.String()).To(Equal(c.id))
		})
	}
}

func TestTypes(t *testing.T) {
	fn := func(a, b string) bool {
		return true
	}

	values := []any{
		typ.AnyStruct[string]{Name: "x"},
		typ.AnySlice[string]{},
		typ.AnyMap[int, string]{},
		typ.IntMap{},
		typ.DeepCompose{},
		func() *typ.Enum { v := typ.ENUM__ONE; return &v }(),
		typ.ENUM__ONE,
		reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem(),
		reflect.TypeOf((*typ.SomeMixInterface)(nil)).Elem(),
		unsafe.Pointer(t),
		make(typ.Chan),
		make(chan string, 100),
		typ.F,
		typ.Func(fn),
		fn,
		typ.String(""),
		"",
		typ.Bool(true),
		true,
		typ.Int(0),
		ptrx.Int(1),
		int(0),
		typ.Int8(0),
		int8(0),
		typ.Int16(0),
		int16(0),
		typ.Int32(0),
		int32(0),
		typ.Int64(0),
		int64(0),
		typ.Uint(0),
		uint(0),
		typ.Uintptr(0),
		uintptr(0),
		typ.Uint8(0),
		uint8(0),
		typ.Uint16(0),
		uint16(0),
		typ.Uint32(0),
		uint32(0),
		typ.Uint64(0),
		uint64(0),
		typ.Float32(0),
		float32(0),
		typ.Float64(0),
		float64(0),
		typ.Complex64(0),
		complex64(0),
		typ.Complex128(0),
		complex128(0),
		typ.Array{},
		[1]string{},
		typ.Slice{},
		[]string{},
		typ.Map{},
		map[string]string{},
		typ.Struct{},
		struct{}{},
		struct {
			typ.Part
			Part2  typ2.Part
			a      string
			A      string `json:"a"`
			Struct struct {
				B string
			}
		}{},
	}

	for i := range values {
		check(t, values[i])
	}
}

func check(t *testing.T, v any) {
	_rt, ok := v.(reflect.Type)
	if !ok {
		_rt = reflect.TypeOf(v)
	}
	gt := FromGoType(NewGoTypeFromReflectType(_rt))
	rt := FromReflectType(_rt)

	t.Run(FullTypeName(rt), func(t *testing.T) {
		w := NewWithT(t)
		w.Expect(rt.String()).To(Equal(gt.String()))
		w.Expect(rt.Kind().String()).To(Equal(gt.Kind().String()))
		w.Expect(rt.Name()).To(Equal(gt.Name()))
		w.Expect(rt.PkgPath()).To(Equal(gt.PkgPath()))
		w.Expect(rt.Comparable()).To(Equal(gt.Comparable()))

		urt := FromReflectType(reflect.TypeOf(""))
		ugt := FromGoType(types.Typ[types.String])
		w.Expect(rt.AssignableTo(urt)).To(Equal(gt.AssignableTo(ugt)))
		w.Expect(rt.ConvertibleTo(urt)).To(Equal(gt.ConvertibleTo(ugt)))

		w.Expect(rt.NumMethod()).To(Equal(gt.NumMethod()))
		for i := 0; i < rt.NumMethod(); i++ {
			rm := rt.Method(i)
			gm, exists := gt.MethodByName(rm.Name())

			t.Run(
				"M_"+rm.Name()+"_"+rm.Type().String(),
				func(t *testing.T) {
					w.Expect(exists).To(BeTrue())
					w.Expect(rm.Name()).To(Equal(gm.Name()))
					w.Expect(rm.PkgPath()).To(Equal(gm.PkgPath()))
					w.Expect(rm.Type().String()).To(Equal(gm.Type().String()))
				})
		}

		{
			_, rok := rt.MethodByName("String")
			_, gok := gt.MethodByName("String")
			w.Expect(rok).To(Equal(gok))
		}
		{
			rplacer, rok := EncodingTextMarshalerTypeReplacer(rt)
			gplacer, gok := EncodingTextMarshalerTypeReplacer(gt)
			w.Expect(rok).To(Equal(gok))
			w.Expect(rplacer.String()).To(Equal(gplacer.String()))
		}

		k := rt.Kind()
		if k == reflect.Func {
			w.Expect(rt.NumIn()).To(Equal(gt.NumIn()))
			w.Expect(rt.NumOut()).To(Equal(gt.NumOut()))
			for i := 0; i < rt.NumIn(); i++ {
				rp, gp := rt.In(i), gt.In(i)
				w.Expect(rp.String()).To(Equal(gp.String()))
			}
			for i := 0; i < rt.NumOut(); i++ {
				rp, gp := rt.Out(i), gt.Out(i)
				w.Expect(rp.String()).To(Equal(gp.String()))
			}
		}
		if k == reflect.Ptr {
			drr := DeRef(rt).(*ReflectType).String()
			drg := DeRef(gt).(*GoType).String()
			w.Expect(drr).To(Equal(drg))
		}
		if k == reflect.Array || k == reflect.Slice || k == reflect.Map {
			w.Expect(FullTypeName(rt.Elem())).To(Equal(FullTypeName(gt.Elem())))
			if k == reflect.Map {
				w.Expect(FullTypeName(rt.Key())).To(Equal(FullTypeName(gt.Key())))
			}
			if k == reflect.Array {
				w.Expect(rt.Len()).To(Equal(gt.Len()))
			}
		}
		if k == reflect.Struct {
			w.Expect(rt.NumField()).To(Equal(gt.NumField()))
			fields := rt.NumField()
			if fields > 0 {
				for i := 0; i < fields; i++ {
					rf, gf := rt.Field(i), gt.Field(i)
					t.Run("F_"+rf.Name(), func(t *testing.T) {
						w.Expect(rf.Anonymous()).To(Equal(gf.Anonymous()))
						w.Expect(rf.Tag()).To(Equal(gf.Tag()))
						w.Expect(rf.Name()).To(Equal(gf.Name()))
						w.Expect(rf.PkgPath()).To(Equal(gf.PkgPath()))
						w.Expect(FullTypeName(rf.Type())).To(Equal(FullTypeName(gf.Type())))
					})
				}
				rf, _ := rt.FieldByName("A")
				gf, _ := rt.FieldByName("A")
				w.Expect(rf.Anonymous()).To(Equal(gf.Anonymous()))
				w.Expect(rf.Tag()).To(Equal(gf.Tag()))
				w.Expect(rf.Name()).To(Equal(gf.Name()))
				w.Expect(rf.PkgPath()).To(Equal(gf.PkgPath()))
				w.Expect(FullTypeName(rf.Type())).To(Equal(FullTypeName(gf.Type())))

				_, rok := rt.FieldByName("_")
				w.Expect(rok).To(BeFalse())
				_, gok := gt.FieldByName("_")
				w.Expect(gok).To(BeFalse())

				rf, _ = rt.FieldByNameFunc(func(s string) bool { return s == "A" })
				gf, _ = gt.FieldByNameFunc(func(s string) bool { return s == "A" })
				w.Expect(rf.Anonymous()).To(Equal(gf.Anonymous()))
				w.Expect(rf.Tag()).To(Equal(gf.Tag()))
				w.Expect(rf.Name()).To(Equal(gf.Name()))
				w.Expect(rf.PkgPath()).To(Equal(gf.PkgPath()))
				w.Expect(FullTypeName(rf.Type())).To(Equal(FullTypeName(gf.Type())))

				_, rok = rt.FieldByNameFunc(func(s string) bool { return false })
				_, gok = gt.FieldByNameFunc(func(s string) bool { return false })
				w.Expect(rok).To(BeFalse())
				w.Expect(gok).To(BeFalse())
			}
		}
	})
}

func TestTryNew(t *testing.T) {
	{
		_, ok := TryNew(FromReflectType(reflect.TypeOf(typ.Struct{})))
		NewWithT(t).Expect(ok).To(BeTrue())
	}
	{
		_, ok := TryNew(FromGoType(NewGoTypeFromReflectType(reflect.TypeOf(typ.Struct{}))))
		NewWithT(t).Expect(ok).To(BeFalse())
	}
}

func TestEachField(t *testing.T) {
	expects := []string{"a", "b", "bool", "c", "Part2"}
	{
		rt := FromReflectType(reflect.TypeOf(typ.Struct{}))
		names := make([]string, 0)
		EachField(rt, "json", func(f StructField, display string, omitempty bool) bool {
			names = append(names, display)
			return true
		})
		NewWithT(t).Expect(expects).To(Equal(names))
	}
	{
		gt := FromGoType(NewGoTypeFromReflectType(reflect.TypeOf(typ.Struct{})))
		names := make([]string, 0)
		EachField(gt, "json", func(f StructField, display string, omitempty bool) bool {
			names = append(names, display)
			return true
		})
		NewWithT(t).Expect(expects).To(Equal(names))
	}
}
