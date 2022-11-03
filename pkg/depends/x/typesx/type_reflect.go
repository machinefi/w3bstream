package typesx

import "reflect"

type ReflectType struct{ reflect.Type }

func FromReflectType(t reflect.Type) *ReflectType { return &ReflectType{Type: t} }

func (t *ReflectType) Unwrap() any { return t.Type }

// func(t*ReflectType)	Name() string
// func(t*ReflectType)	PkgPath() string

func (t *ReflectType) String() string { return TypeString(t) }

// func(t*ReflectType)	Kind() reflect.Kind

func (t *ReflectType) Implements(u Type) bool {
	switch x := u.(type) {
	case *ReflectType:
		return t.Type.Implements(x.Type)
		// case *GoType:
		// 	if t.PkgPath() == "" {
		// 		return false
		// 	}
		// 	if i, ok := x.Type.(*types.Interface); ok {
		// 		return types.Implements(NewGoTypeFromReflectType(t.Type), i)
		// 	}
	}
	return false
}

func (t *ReflectType) AssignableTo(u Type) bool {
	if rt, ok := u.(*ReflectType); ok {
		return t.Type.AssignableTo(rt.Type)
	}
	return false
}

func (t *ReflectType) ConvertibleTo(u Type) bool {
	if rt, ok := u.(*ReflectType); ok {
		return t.Type.ConvertibleTo(rt.Type)
	}
	return false
}

// func (t *ReflectType) Comparable() bool

func (t *ReflectType) Key() Type { return FromReflectType(t.Type.Key()) }

func (t *ReflectType) Elem() Type { return FromReflectType(t.Type.Elem()) }

// func(t *ReflectType)	Len() int
// func(t *ReflectType)	NumField() int

func (t *ReflectType) Field(i int) StructField {
	return &ReflectStructField{StructField: t.Type.Field(i)}
}

func (t *ReflectType) FieldByName(name string) (StructField, bool) {
	if f, ok := t.Type.FieldByName(name); ok {
		return &ReflectStructField{StructField: f}, true
	}
	return nil, false
}

func (t *ReflectType) FieldByNameFunc(match func(string) bool) (StructField, bool) {
	if f, ok := t.Type.FieldByNameFunc(match); ok {
		return &ReflectStructField{StructField: f}, true
	}
	return nil, false
}

// func (t *ReflectType) NumMethod() int

func (t *ReflectType) Method(i int) Method { return &ReflectMethod{Method: t.Type.Method(i)} }

func (t *ReflectType) MethodByName(name string) (Method, bool) {
	if m, ok := t.Type.MethodByName(name); ok {
		return &ReflectMethod{m}, true
	}
	return nil, false
}

// func (t *ReflectType) IsVariadic() bool

// func (t *ReflectType) NumIn() int

func (t *ReflectType) In(i int) Type { return FromReflectType(t.Type.In(i)) }

// func (t *ReflectType) NumOut() int

func (t *ReflectType) Out(i int) Type { return FromReflectType(t.Type.Out(i)) }

type ReflectMethod struct{ Method reflect.Method }

func (m *ReflectMethod) PkgPath() string { return m.Method.PkgPath }

func (m *ReflectMethod) Name() string { return m.Method.Name }

func (m *ReflectMethod) Type() Type { return FromReflectType(m.Method.Type) }

type ReflectStructField struct{ StructField reflect.StructField }

func (f *ReflectStructField) PkgPath() string { return f.StructField.PkgPath }

func (f *ReflectStructField) Name() string { return f.StructField.Name }

func (f *ReflectStructField) Tag() reflect.StructTag { return f.StructField.Tag }

func (f *ReflectStructField) Type() Type { return FromReflectType(f.StructField.Type) }

func (f *ReflectStructField) Anonymous() bool { return f.StructField.Anonymous }
