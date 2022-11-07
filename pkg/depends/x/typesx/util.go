package typesx

import (
	"bytes"
	"fmt"
	"go/types"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
)

var (
	typs   = mapx.New[string, types.Type]()
	pkgs   = mapx.New[string, *types.Package]()
	basics = map[string]types.Type{}

	LoadFiles   = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles
	LoadImports = LoadFiles | packages.NeedImports
	LoadTypes   = LoadImports | packages.NeedTypes | packages.NeedTypesSizes
)

// init basic types
func init() {
	for _, b := range types.Typ {
		basics[types.TypeString(b, nil)] = b
	}
	basics["interface {}"] = types.NewInterfaceType(nil, nil)
	basics["error"] = NewPackage("errors").Scope().Lookup("New").Type().
		Underlying().(*types.Signature).Results().At(0).Type()
}

func NewGoTypeFromReflectType(t reflect.Type) types.Type {
	underlying := func() types.Type {
		k := t.Kind()
		if IsBasicReflectKind(k) {
			return types.Typ[ReflectKindToTypesKind[k]]
		}
		switch k {
		case reflect.Array:
			return types.NewArray(
				NewGoTypeFromReflectType(t.Elem()),
				int64(t.Len()),
			)
		case reflect.Slice:
			return types.NewSlice(NewGoTypeFromReflectType(t.Elem()))
		case reflect.Map:
			return types.NewMap(
				NewGoTypeFromReflectType(t.Key()),
				NewGoTypeFromReflectType(t.Elem()),
			)
		case reflect.Chan:
			return types.NewChan(
				types.ChanDir(t.ChanDir()),
				NewGoTypeFromReflectType(t.Elem()),
			)
		case reflect.Func:
			params := make([]*types.Var, t.NumIn())
			for i := range params {
				v := t.In(i)
				params[i] = types.NewParam(0, NewPackage(v.PkgPath()), "", NewGoTypeFromReflectType(v))
			}
			results := make([]*types.Var, t.NumOut())
			for i := range results {
				v := t.Out(i)
				results[i] = types.NewParam(0, NewPackage(v.PkgPath()), "", NewGoTypeFromReflectType(v))
			}
			return types.NewSignatureType(
				nil,
				nil,
				nil,
				types.NewTuple(params...),
				types.NewTuple(results...),
				t.IsVariadic(),
			)
		case reflect.Interface:
			fns := make([]*types.Func, t.NumMethod())
			for i := range fns {
				f := t.Method(i)
				fns[i] = types.NewFunc(
					0,
					NewPackage(f.PkgPath),
					f.Name,
					NewGoTypeFromReflectType(f.Type).(*types.Signature),
				)
			}
			return types.NewInterfaceType(fns, nil).Complete()
		case reflect.Struct:
			fields := make([]*types.Var, t.NumField())
			tags := make([]string, len(fields))
			for i := range fields {
				f := t.Field(i)
				fields[i] = types.NewField(
					0,
					NewPackage(f.PkgPath),
					f.Name,
					NewGoTypeFromReflectType(f.Type),
					f.Anonymous,
				)
				tags[i] = string(f.Tag)
			}
			return types.NewStruct(fields, tags)
		}
		return nil
	}

	stars := 0

	indirect := func(t types.Type) types.Type {
		for stars > 0 {
			t = types.NewPointer(t)
			stars--
		}
		return t
	}

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		stars++
	}
	name := t.Name()
	path := t.PkgPath()
	if name == "error" && path == "" {
		return nil
	}

	if path != "" {
		return indirect(TypeFor(path + "." + name))
	}
	return indirect(underlying())
}

func NewPackage(path string) *types.Package {
	if path == "" {
		return nil
	}
	if v, ok := pkgs.Load(path); ok {
		return v
	}
	cfg := packages.Config{
		Overlay: make(map[string][]byte),
		Tests:   true,
		Mode:    LoadTypes,
	}
	pkg, err := packages.Load(&cfg, path)
	if err != nil {
		panic(err)
	}
	pkgs.Store(path, pkg[0].Types)
	return pkg[0].Types
}

func TypeByName(path string, name string) types.Type {
	if path == "" {
		TypeFor(name)
	}
	return TypeFor(path + "." + name)
}

func PtrTo(t Type) Type {
	switch x := t.(type) {
	case *GoType:
		return FromGoType(types.NewPointer(x.Type))
	case *ReflectType:
		return FromReflectType(reflect.PtrTo(x.Type))
	}
	return nil
}

func TypeString(t Type) string {
	if pkg := t.PkgPath(); pkg != "" {
		return pkg + "." + t.Name()
	}
	k := t.Kind()
	if IsBasicReflectKind(k) {
		return k.String()
	}

	switch k {
	case reflect.Slice:
		return "[]" + t.Elem().String()
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), t.Elem().String())
	case reflect.Chan:
		return "chan " + t.Elem().String()
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", t.Key().String(), t.Elem().String())
	case reflect.Struct:
		b := bytes.NewBuffer(nil)
		b.WriteString("struct {")
		n := t.NumField()
		for i := 0; i < n; i++ {
			b.WriteRune(' ')
			f := t.Field(i)
			if !f.Anonymous() {
				b.WriteString(f.Name())
				b.WriteRune(' ')
			}
			b.WriteString(f.Type().String())
			tag := f.Tag()
			if tag != "" {
				b.WriteRune(' ')
				b.WriteString(strconv.Quote(string(tag)))
			}
			if i == n-1 {
				b.WriteRune(' ')
			} else {
				b.WriteString(";")
			}
		}
		b.WriteString("}")
		return b.String()
	case reflect.Interface:
		if name := t.Name(); name == "error" {
			return name
		}
		b := bytes.NewBuffer(nil)
		b.WriteString("interface {")
		n := t.NumMethod()
		for i := 0; i < n; i++ {
			b.WriteRune(' ')
			m := t.Method(i)
			pkg := m.PkgPath()
			if pkg != "" {
				b.WriteString(NewPackage(pkg).Name())
				b.WriteRune('.')
			}
			b.WriteString(m.Name())
			b.WriteString(m.Type().String()[4:])

			if i == n-1 {
				b.WriteRune(' ')
			} else {
				b.WriteRune(';')
			}
		}
		b.WriteString("}")
		return b.String()
	case reflect.Func:
		b := bytes.NewBuffer(nil)
		b.WriteString("func(")
		{
			n := t.NumIn()
			for i := 0; i < n; i++ {
				p := t.In(i)
				if i == n-1 && t.IsVariadic() {
					b.WriteString("...")
					b.WriteString(p.Elem().String())
				} else {
					b.WriteString(p.String())
				}
				if i < n-1 {
					b.WriteString(", ")
				}
			}
			b.WriteString(")")
		}
		{
			n := t.NumOut()
			if n > 0 {
				b.WriteRune(' ')
			}
			if n > 1 {
				b.WriteString("(")
			}
			for i := 0; i < n; i++ {
				if i > 0 {
					b.WriteString(", ")
				}
				r := t.Out(i)
				b.WriteString(r.String())
			}
			if n > 1 {
				b.WriteString(")")
			}
		}
		return b.String()
	}
	return t.Name()
}

func TypeFor(id string) (t types.Type) {
	if v, ok := typs.Load(id); ok {
		return v.(types.Type)
	}

	defer func() {
		if t == nil {
			t = types.Typ[types.Invalid]
		}
		typs.Store(id, t)
	}()

	if id == "" {
		return
	}

	if basic, ok := basics[id]; ok {
		t = basic
		return
	}

	// map[x]
	l := strings.Index(id, "map[")
	if l == 0 {
		r := strings.Index(id, "]")
		t = types.NewMap(TypeFor(id[4:r]), TypeFor(id[r+1:]))
		return
	}

	// []x [n]x
	l = strings.Index(id, "[")
	if l == 0 {
		r := strings.Index(id, "]")
		if l == r-1 {
			t = types.NewSlice(TypeFor(id[r+1:]))
			return
		}
		n, err := strconv.ParseInt(id[1:r], 10, 64)
		if err != nil {
			// panic(err)
			return // invalid
		}
		t = types.NewArray(TypeFor(id[r+1:]), n)
		return
	} else if l == -1 {
		i := strings.LastIndex(id, ".")
		if i <= 0 {
			return // invalid
		}
		path := id[0:i]
		name := id[i+1:]
		pkg := NewPackage(path)
		if pkg == nil {
			return
		}
		if found := pkg.Scope().Lookup(name); found != nil {
			t = found.Type()
			return
		}
		return
	} else {
		r := strings.Index(id, "]")
		full := id[0:l]
		paramNames := strings.Split(id[l+1:r], ",") // github.com/x/y/z.AnyStruct[int,string]
		if dot := strings.LastIndex(full, "."); dot > 0 {
			path, name := full[0:dot], full[dot+1:]
			if p := NewPackage(path); p != nil {
				if found := p.Scope().Lookup(name); found != nil {
					named := &(*found.(*types.TypeName).Type().(*types.Named))
					paramTypes := named.TypeParams()
					if n := paramTypes.Len(); n > 0 {
						params := make([]*types.TypeParam, n)
						for i := 0; i < n; i++ {
							params[i] = types.NewTypeParam(
								paramTypes.At(i).Obj(),
								TypeFor(paramNames[i]),
							)
						}
						named.SetTypeParams(params)
					}
					return found.Type()
				}
			}
		}
	}
	return types.Typ[types.Invalid]
}

var ReflectKindToTypesKind = map[reflect.Kind]types.BasicKind{
	reflect.Bool:          types.Bool,
	reflect.Int:           types.Int,
	reflect.Int8:          types.Int8,
	reflect.Int16:         types.Int16,
	reflect.Int32:         types.Int32,
	reflect.Int64:         types.Int64,
	reflect.Uint:          types.Uint,
	reflect.Uint8:         types.Uint8,
	reflect.Uint16:        types.Uint16,
	reflect.Uint32:        types.Uint32,
	reflect.Uint64:        types.Uint64,
	reflect.Uintptr:       types.Uintptr,
	reflect.Float32:       types.Float32,
	reflect.Float64:       types.Float64,
	reflect.Complex64:     types.Complex64,
	reflect.Complex128:    types.Complex128,
	reflect.String:        types.String,
	reflect.UnsafePointer: types.UnsafePointer,
}

var TypesKindToReflectKind = map[types.BasicKind]reflect.Kind{
	types.Bool:           reflect.Bool,
	types.Int:            reflect.Int,
	types.Int8:           reflect.Int8,
	types.Int16:          reflect.Int16,
	types.Int32:          reflect.Int32,
	types.Int64:          reflect.Int64,
	types.Uint:           reflect.Uint,
	types.Uint8:          reflect.Uint8,
	types.Uint16:         reflect.Uint16,
	types.Uint32:         reflect.Uint32,
	types.Uint64:         reflect.Uint64,
	types.Uintptr:        reflect.Uintptr,
	types.Float32:        reflect.Float32,
	types.Float64:        reflect.Float64,
	types.Complex64:      reflect.Complex64,
	types.Complex128:     reflect.Complex128,
	types.String:         reflect.String,
	types.UnsafePointer:  reflect.UnsafePointer,
	types.UntypedBool:    reflect.Bool,
	types.UntypedInt:     reflect.Int,
	types.UntypedRune:    reflect.Int32,
	types.UntypedFloat:   reflect.Float32,
	types.UntypedComplex: reflect.Complex64,
	types.UntypedString:  reflect.String,
}

func IsBasicReflectKind(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.String, reflect.UnsafePointer:
		return true
	}
	return false
}

func IsNumericReflectKind(k reflect.Kind) bool {
	return IsIntegerReflectKind(k) || IsFloatReflectKind(k)
}

func IsIntegerReflectKind(k reflect.Kind) bool {
	return IsSignedIntReflectKind(k) || IsUnsignedIntReflectKind(k)
}

func IsSignedIntReflectKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	}
	return false
}

func IsUnsignedIntReflectKind(k reflect.Kind) bool {
	switch k {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
}

func IsFloatReflectKind(k reflect.Kind) bool {
	switch k {
	case reflect.Float32, reflect.Float64:
		return true
	}
	return false
}
