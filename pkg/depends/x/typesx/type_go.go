package typesx

import (
	"go/ast"
	"go/types"
	"reflect"
	"strings"
)

type GoType struct {
	Type       types.Type
	scanned    bool
	methods    []*types.Func
	ptrMethods []*types.Func
}

var _ Type = (*GoType)(nil)

func FromGoType(t types.Type) (gt *GoType) {
	if p, ok := t.(*types.TypeParam); ok {
		return &GoType{Type: p.Constraint()}
	}
	return &GoType{Type: t}
}

func (t *GoType) Unwrap() any { return t.Type }

func (t *GoType) Name() string {
	switch x := t.Type.(type) {
	case *types.Named:
		b := strings.Builder{}
		b.WriteString(x.Obj().Name())
		params := x.TypeParams()
		if n := params.Len(); n > 0 {
			b.WriteString("[")
			for i := 0; i < n; i++ {
				if i > 0 {
					b.WriteString(",")
				}
				a := params.At(i).Constraint().(*types.Interface)
				if a.NumEmbeddeds() > 0 {
					b.WriteString(TypeString(FromGoType(a.EmbeddedType(0))))
				} else {
					b.WriteString(TypeString(FromGoType(a)))
				}
			}
			b.WriteString("]")
		}
		return b.String()
	case *types.Basic:
		return x.Name()
	}
	return ""
}

func (t *GoType) PkgPath() string {
	switch x := t.Type.(type) {
	case *types.Named:
		if pkg := x.Obj().Pkg(); pkg != nil {
			return pkg.Path()
		}
		if x.String() == "error" {
			return ""
		}
	case *types.Basic:
		if strings.HasPrefix(x.String(), "unsafe.") {
			return "unsafe" // unsafe.Pointer as basic since 1.17
		}
	}
	return ""
}

func (t *GoType) String() string { return TypeString(t) }

func (t *GoType) Kind() reflect.Kind {
	switch x := t.Type.(type) {
	case *types.Named:
		pkg := x.Obj().Pkg()
		if pkg != nil && pkg.Name() == "unsafe" && x.Obj().Name() == "Pointer" {
			return reflect.UnsafePointer
		}
		return FromGoType(x.Underlying()).Kind()
	case *types.Interface:
		return reflect.Interface
	case *types.Struct:
		return reflect.Struct
	case *types.Pointer:
		return reflect.Pointer
	case *types.Map:
		return reflect.Map
	case *types.Slice:
		return reflect.Slice
	case *types.Array:
		return reflect.Array
	case *types.Chan:
		return reflect.Chan
	case *types.Signature:
		return reflect.Func
	case *types.Basic:
		if k, ok := TypesKindToReflectKind[x.Kind()]; ok {
			return k
		}
	}
	return reflect.Invalid
}

func (t *GoType) Implements(u Type) bool {
	switch x := u.(type) {
	case *GoType:
		return types.Implements(t.Type, x.Type.(*types.Interface))
	case *ReflectType:
		var (
			v   Type = t
			ptr      = false
		)
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
			ptr = true
		}
		if v.PkgPath() == "" || x.PkgPath() == "" {
			return false
		}
		if ptr {
			return types.Implements(
				types.NewPointer(TypeByName(v.PkgPath(), v.Name())),
				NewGoTypeFromReflectType(x.Type).Underlying().(*types.Interface),
			)
		}
		return types.Implements(
			TypeByName(v.PkgPath(), v.Name()),
			NewGoTypeFromReflectType(x.Type).Underlying().(*types.Interface),
		)
	}
	return false
}

func (t *GoType) AssignableTo(u Type) bool {
	return types.AssignableTo(t.Type, u.(*GoType).Type)
}

func (t *GoType) ConvertibleTo(u Type) bool {
	return types.ConvertibleTo(t.Type, u.(*GoType).Type)
}

func (t *GoType) Comparable() bool {
	if t.Kind() == reflect.Struct {
		return true
	}
	return types.Comparable(t.Type)
}

func (t *GoType) Key() Type {
	switch x := t.Type.(type) {
	case *types.Named:
		return FromGoType(ConstraintUnderlying(x.TypeParams(), x.Underlying())).Key()
	case interface{ Key() types.Type }:
		return FromGoType(x.Key())
	}
	return nil
}

func (t *GoType) Elem() Type {
	switch x := t.Type.(type) {
	case *types.Named:
		return FromGoType(ConstraintUnderlying(x.TypeParams(), x.Underlying())).Elem()
	case interface{ Elem() types.Type }:
		return FromGoType(x.Elem())
	}
	return nil
}

func (t *GoType) Len() int {
	switch x := t.Type.(type) {
	case *types.Named:
		return FromGoType(x.Underlying()).Len()
	case *types.Array:
		return int(x.Len())
	}
	return 0
}

func (t *GoType) NumField() int {
	switch x := t.Type.(type) {
	case *types.Pointer:
		return FromGoType(x.Elem()).NumField()
	case *types.Named:
		return FromGoType(x.Underlying()).NumField()
	case *types.Struct:
		return x.NumFields()
	}
	return 0
}

func (t *GoType) Field(i int) StructField {
	switch x := t.Type.(type) {
	case *types.Named:
		return FromGoType(ConstraintUnderlying(x.TypeParams(), x.Underlying())).Field(i)
	case *types.Struct:
		return &GoStructField{Var: x.Field(i), TagValue: x.Tag(i)}
	}
	return nil
}

func (t *GoType) FieldByName(name string) (StructField, bool) {
	return t.FieldByNameFunc(func(s string) bool { return name == s })
}

func (t *GoType) FieldByNameFunc(match func(string) bool) (StructField, bool) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if match(f.Name()) {
			return f, true
		}
		if f.Anonymous() {
			if _f, ok := f.Type().FieldByNameFunc(match); ok {
				return _f, ok
			}
		}
	}
	return nil, false
}

func (t *GoType) scanMethods() {
	if t.scanned {
		return
	}
	t.scanned = true
	methods := MethodsOfGoType(t.Type)
	for _, m := range methods {
		if !m.Ptr {
			t.methods = append(t.methods, m.Func)
		}
		t.ptrMethods = append(t.ptrMethods, m.Func)
	}
}

func (t *GoType) NumMethod() int {
	if t.Kind() == reflect.Interface {
		switch x := t.Type.(type) {
		case *types.Named:
			return x.Underlying().(*types.Interface).NumMethods()
		case *types.Interface:
			return x.NumMethods()
		}
	}

	t.scanMethods()

	switch t.Type.(type) {
	case *types.Pointer:
		return len(t.ptrMethods)
	default:
		return len(t.methods)
	}
}

func (t *GoType) Method(i int) Method {
	if t.Kind() == reflect.Interface {
		switch x := t.Type.(type) {
		case *types.Named:
			return &GoMethod{Func: x.Underlying().(*types.Interface).Method(i)}
		case *types.Interface:
			return &GoMethod{Func: x.Method(i)}
		}
	}

	t.scanMethods()

	switch t.Type.(type) {
	case *types.Pointer:
		if t.ptrMethods != nil {
			return &GoMethod{Recv: t, Func: t.ptrMethods[i]}
		}
	default:
		if t.methods != nil {
			return &GoMethod{Recv: t, Func: t.methods[i]}
		}
	}
	return nil
}

func (t *GoType) MethodByName(name string) (Method, bool) {
	for i := 0; i < t.NumMethod(); i++ {
		if m := t.Method(i); m.Name() == name {
			return m, true
		}
	}
	return nil, false
}

func (t *GoType) IsVariadic() bool {
	if sig, ok := t.Type.(*types.Signature); ok {
		return sig.Variadic()
	}
	return false
}

func (t *GoType) NumIn() int {
	switch x := t.Type.(type) {
	case *types.Named:
		return FromGoType(x.Underlying()).NumIn()
	case *types.Signature:
		return x.Params().Len()
	}
	return 0
}

func (t *GoType) In(i int) Type {
	switch x := t.Type.(type) {
	case *types.Named:
		return FromGoType(x.Underlying()).In(i)
	case *types.Signature:
		return FromGoType(x.Params().At(i).Type())
	}
	return nil
}

func (t *GoType) NumOut() int {
	switch x := t.Type.(type) {
	case *types.Named:
		return FromGoType(x.Underlying()).NumOut()
	case *types.Signature:
		return x.Results().Len()
	}
	return 0
}

func (t *GoType) Out(i int) Type {
	switch x := t.Type.(type) {
	case *types.Named:
		return FromGoType(x.Underlying()).Out(i)
	case *types.Signature:
		return FromGoType(x.Results().At(i).Type())
	}
	return nil
}

type GoMethod struct {
	Ptr  bool
	Recv *GoType
	Func *types.Func
}

func MethodsOfGoType(t types.Type) (methods []*GoMethod) {
	switch x := t.(type) {
	case *types.Named:
		for i := 0; i < x.NumMethods(); i++ {
			info := &GoMethod{Func: x.Method(i)}
			if _, ok := info.Func.Type().(*types.Signature).Recv().Type().(*types.Pointer); ok {
				info.Ptr = true
			}
			methods = append(methods, info)
		}
		methods = append(methods, MethodsOfGoType(x.Underlying())...)
	case *types.Pointer:
		methods = append(methods, MethodsOfGoType(x.Elem())...)
	case *types.Struct:
		for i := 0; i < x.NumFields(); i++ {
			f := x.Field(i)
			if f.Anonymous() {
				methods = append(methods, MethodsOfGoType(f.Type())...)
			}
		}
	}
	return
}

func (m *GoMethod) PkgPath() string {
	if ast.IsExported(m.Name()) {
		return ""
	}
	if pkg := m.Func.Pkg(); pkg != nil {
		return pkg.Path()
	}
	return ""
}

func (m *GoMethod) Name() string { return m.Func.Name() }

func (m *GoMethod) Type() Type {
	sig := m.Func.Type().(*types.Signature)

	if m.Recv == nil {
		return FromGoType(sig)
	}

	vars := make([]*types.Var, sig.Params().Len()+1)
	vars[0] = types.NewVar(0, nil, "", m.Recv.Type)
	for i := 0; i < sig.Params().Len(); i++ {
		vars[i+1] = sig.Params().At(i)
	}
	return FromGoType(types.NewSignatureType(
		nil,
		nil,
		nil,
		types.NewTuple(vars...),
		sig.Results(),
		sig.Variadic(),
	))
}

type GoStructField struct {
	*types.Var
	TagValue string
}

func (f *GoStructField) PkgPath() string {
	if ast.IsExported(f.Name()) {
		return ""
	}
	if pkg := f.Pkg(); pkg != nil {
		return pkg.Path()
	}
	return ""
}

func (f *GoStructField) Tag() reflect.StructTag { return reflect.StructTag(f.TagValue) }

func (f *GoStructField) Type() Type { return FromGoType(f.Var.Type()) }

func ConstraintUnderlying(params *types.TypeParamList, underlying types.Type) types.Type {
	if params.Len() == 0 {
		return underlying
	}
	switch t := underlying.(type) {
	case *types.TypeParam:
		a := params.At(t.Index()).Constraint().(*types.Interface)
		if a.NumEmbeddeds() > 0 {
			return a.EmbeddedType(0)
		}
		return a
	case *types.Map:
		return types.NewMap(
			ConstraintUnderlying(params, t.Key()),
			ConstraintUnderlying(params, t.Elem()),
		)
	case *types.Slice:
		return types.NewSlice(ConstraintUnderlying(params, t.Elem()))
	case *types.Array:
		return types.NewArray(ConstraintUnderlying(params, t.Elem()), t.Len())
	case *types.Struct:
		n := t.NumFields()
		tags, fields := make([]string, n), make([]*types.Var, n)
		for i := 0; i < n; i++ {
			tags[i] = t.Tag(i)
			f := t.Field(i)
			fields[i] = types.NewField(
				f.Pos(),
				f.Pkg(),
				f.Name(),
				ConstraintUnderlying(params, f.Type()),
				f.Embedded(),
			)
		}
		return types.NewStruct(fields, tags)
	}
	return underlying
}
