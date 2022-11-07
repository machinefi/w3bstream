package codegen

import (
	"path"
	"reflect"
	"unicode"
)

func IsEmptyValue(rv reflect.Value) bool {
	if !rv.IsValid() || !rv.CanInterface() {
		return false
	}

	if chk, ok := rv.Interface().(interface{ IsZero() bool }); ok && chk.IsZero() {
		return false
	}

	switch rv.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return rv.IsNil()
	}

	return false
}

func IsValidIdent(s string) bool {
	if len(s) == 0 {
		return false
	}
	if IsReserved(s) {
		return false
	}
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' {
			return false
		}
	}
	return true
}

func IsReserved(s string) bool {
	_, ok := ReservedKeys[s]
	if ok {
		return ok
	}
	_, ok = BuiltinFuncs[s]
	return ok
}

func IsBuiltinFunc(s string) bool {
	_, ok := BuiltinFuncs[s]
	return ok
}

const (
	Bool       BuiltInType = "bool"
	Int        BuiltInType = "int"
	Int8       BuiltInType = "int8"
	Int16      BuiltInType = "int16"
	Int32      BuiltInType = "int32"
	Int64      BuiltInType = "int64"
	Uint       BuiltInType = "uint"
	Uint8      BuiltInType = "uint8"
	Uint16     BuiltInType = "uint16"
	Uint32     BuiltInType = "uint32"
	Uint64     BuiltInType = "uint64"
	Uintptr    BuiltInType = "uintptr"
	Float32    BuiltInType = "float32"
	Float64    BuiltInType = "float64"
	Complex64  BuiltInType = "complex64"
	Complex128 BuiltInType = "complex128"
	String     BuiltInType = "string"
	Byte       BuiltInType = "byte"
	Rune       BuiltInType = "rune"
	Error      BuiltInType = "error"
)

const (
	Iota        SnippetBuiltIn = "iota"
	True        SnippetBuiltIn = "true"
	False       SnippetBuiltIn = "false"
	Nil         SnippetBuiltIn = "nil"
	Break       SnippetBuiltIn = "break"
	Continue    SnippetBuiltIn = "continue"
	Fallthrough SnippetBuiltIn = "fallthrough"
)

var (
	ReservedKeys = map[string]bool{
		"bool":        true,
		"int":         true,
		"int8":        true,
		"int16":       true,
		"int32":       true,
		"int64":       true,
		"uint":        true,
		"uint8":       true,
		"uint16":      true,
		"uint32":      true,
		"uint64":      true,
		"uintptr":     true,
		"float32":     true,
		"float64":     true,
		"complex64":   true,
		"complex128":  true,
		"string":      true,
		"byte":        true,
		"rune":        true,
		"error":       true,
		"iota":        true,
		"true":        true,
		"false":       true,
		"nil":         true,
		"break":       true,
		"continue":    true,
		"fallthrough": true,
	}
	BuiltinFuncs = map[string]bool{
		"append":  true,
		"complex": true,
		"cap":     true,
		"close":   true,
		"copy":    true,
		"delete":  true,
		"imag":    true,
		"len":     true,
		"make":    true,
		"new":     true,
		"panic":   true,
		"print":   true,
		"println": true,
		"real":    true,
		"recover": true,
	}
)

const AnonymousIdent SnippetIdent = "_"

var naming = path.Base

func SetPkgNaming(fn func(string) string) {
	if fn != nil {
		naming = fn
	}
}

var (
	Valuer = ValueWithAlias(naming)
	Typer  = TypeWithAlias(naming)
	Exprer = ExprWithAlias(naming)
)

func Stringify(s Snippet) string { return string(s.Bytes()) }
