package reflectx

import (
	"bytes"
	"reflect"
)

func Indirect(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return Indirect(v.Elem())
	}
	return v
}

// New a `reflect.Value` with reflect.Type
func New(t reflect.Type) reflect.Value {
	v := reflect.New(t).Elem()
	if t.Kind() == reflect.Ptr {
		v.Set(New(t.Elem()).Addr())
	}
	return v
}

func IsEmptyValue(v interface{}) bool {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(&v).Elem()
	}
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return true
	}
	if rv.Kind() == reflect.Interface && rv.IsNil() {
		return true
	}
	if rv.Kind() == reflect.Invalid {
		return true
	}

	if rv.IsValid() && rv.CanInterface() {
		if chk, ok := rv.Interface().(interface{ IsZero() bool }); ok {
			return chk.IsZero()
		}
	}
	switch rv.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Ptr, reflect.Interface:
		return IsEmptyValue(rv.Elem())
	}
	return false
}

func TypeName(rt reflect.Type) string {
	buf := bytes.NewBuffer(nil)
	for rt.Kind() == reflect.Ptr {
		buf.WriteByte('*')
		rt = rt.Elem()
	}
	if pkg := rt.PkgPath(); pkg != "" {
		buf.WriteString(pkg)
		buf.WriteByte('.')
	}
	if name := rt.Name(); name != "" {
		buf.WriteString(name)
	}
	buf.WriteString(rt.String())
	return buf.String()
}

func DeRef(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// NatureType return the v's nature type, for example:
// Foo nature type is Foo
// *Foo nature type is Foo
// ***Foo nature type is Foo
// reflect.TypeOf(Foo) nature type is Foo
// (interface{})(Foo) nature type is Foo
func NatureType(v interface{}) (rt reflect.Type) {
	if !reflect.ValueOf(v).IsValid() {
		return reflect.TypeOf(nil)
	}

	rt = reflect.TypeOf(v)
	if t, ok := v.(reflect.Type); ok {
		rt = t
	}

	kind := rt.Kind()
	for kind == reflect.Ptr {
		rt = rt.Elem()
		kind = rt.Kind()
	}

	if kind == reflect.Interface {
		return NatureType(reflect.New(rt).Elem().Interface())
	}
	return rt
}

func IsBytes(v interface{}) bool {
	if _, ok := v.([]byte); ok {
		return true
	}
	t := BasicAssertReflectType(v)
	return IsBytesType(t)
}

func IsBytesType(t reflect.Type) bool {
	return t.Kind() == reflect.Slice &&
		t.Elem().Kind() == reflect.Uint8 &&
		t.Elem().PkgPath() == ""
}

func BasicAssertReflectType(v interface{}) reflect.Type {
	var t reflect.Type
	switch x := v.(type) {
	case reflect.Type:
		t = x
	case interface{ Type() reflect.Type }:
		t = x.Type()
	default:
		t = reflect.TypeOf(v)
	}
	return t
}
