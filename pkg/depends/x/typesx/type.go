package typesx

import (
	"encoding"
	"go/ast"
	"go/types"
	"reflect"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
)

type Type interface {
	Unwrap() any

	Name() string
	PkgPath() string
	String() string
	Kind() reflect.Kind
	Implements(Type) bool
	AssignableTo(Type) bool
	ConvertibleTo(Type) bool
	Comparable() bool

	Key() Type
	Elem() Type
	Len() int

	NumField() int
	Field(int) StructField
	FieldByName(string) (StructField, bool)
	FieldByNameFunc(func(string) bool) (StructField, bool)

	NumMethod() int
	Method(int) Method
	MethodByName(string) (Method, bool)

	IsVariadic() bool
	NumIn() int
	In(int) Type
	NumOut() int
	Out(int) Type
}

type Method interface {
	PkgPath() string
	Name() string
	Type() Type
}

type StructField interface {
	PkgPath() string
	Name() string
	Tag() reflect.StructTag
	Type() Type
	Anonymous() bool
}

func TryNew(u Type) (reflect.Value, bool) {
	if v, ok := u.(*ReflectType); ok {
		return reflectx.New(v.Type), true
	}
	return reflect.Value{}, false
}

var RtTextMarshaler = FromReflectType(reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem())

func EncodingTextMarshalerTypeReplacer(u Type) (Type, bool) {
	switch x := u.(type) {
	case *GoType:
		return FromGoType(types.Typ[types.String]), x.Implements(RtTextMarshaler)
	case *ReflectType:
		return FromReflectType(reflect.TypeOf("")), x.Implements(RtTextMarshaler)
	}
	return nil, false
}

func DeRef(t Type) Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func FullTypeName(t Type) string {
	b := strings.Builder{}

	for t.Kind() == reflect.Ptr {
		b.WriteByte('*')
		t = t.Elem()
	}

	if name := t.Name(); name != "" {
		if path := t.PkgPath(); path != "" {
			b.WriteString(path)
			b.WriteRune('.')
		}
		b.WriteString(name)
		return b.String()
	}
	b.WriteString(b.String())
	return b.String()
}

func FieldDisplayName(tag reflect.StructTag, key string, name string) (keyTag string, omitempty, exists bool) {
	keyTag, exists = tag.Lookup(key)
	if !exists {
		keyTag = name
		return
	}
	omitempty = strings.Index(keyTag, "omitempty") > 0
	commaIdx := strings.IndexRune(keyTag, ',')
	if keyTag == "" || commaIdx == 0 {
		return name, omitempty, true
	}
	if commaIdx == -1 {
		return
	}
	keyTag = keyTag[0:commaIdx]
	return
}

func EachField(t Type, key string, each func(f StructField, display string, omitempty bool) bool, keepNestedTags ...string) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fname := f.Name()
		ftag := f.Tag()

		display, omitempty, exists := FieldDisplayName(ftag, key, fname)

		if !ast.IsExported(fname) || display == "-" {
			continue
		}
		ftype := DeRef(f.Type())
		if f.Anonymous() {
			k := ftype.Kind()
			if k == reflect.Interface {
				continue
			}
			if k == reflect.Struct {
				if !exists {
					for _, tag := range keepNestedTags {
						if _, ok := ftag.Lookup(tag); ok {
							exists = true
							break
						}
					}
				}
				if !exists {
					EachField(ftype, key, each)
					continue
				}
			}
		}
		if !each(f, display, omitempty) {
			break
		}
	}
}
