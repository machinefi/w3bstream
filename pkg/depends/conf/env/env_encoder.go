package env

import (
	"go/ast"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
	"github.com/machinefi/w3bstream/pkg/depends/x/textx"
)

type Encoder struct {
	vars *Vars
	sets map[string]map[string]bool
}

func NewEncoder(vars *Vars) *Encoder { return &Encoder{vars: vars} }

func (e *Encoder) Encode(v interface{}) ([]byte, error) {
	w := NewPathWalker()
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	if err := e.scan(w, rv); err != nil {
		return nil, err
	}
	return e.vars.Bytes(), nil
}

func (e *Encoder) SecurityEncode(v interface{}) ([]byte, error) {
	w := NewPathWalker()
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	if err := e.scan(w, rv); err != nil {
		return nil, err
	}
	return e.vars.MaskBytes(), nil
}

func (e *Encoder) flag(key string, flags map[string]bool) {
	if e.sets == nil {
		e.sets = make(map[string]map[string]bool)
	}
	e.sets[key] = flags
}

func (e *Encoder) scan(w *PathWalker, rv reflect.Value) error {
	kind := rv.Kind()

	setter := func(rv reflect.Value) error {
		key := w.String()
		v := &Var{Name: key}
		if e.sets != nil {
			if flags, ok := e.sets[key]; ok {
				v.SetMeta(flags)
			} else {
				name := StringifyPath(w.path[0 : len(w.path)-1]...)
				if flags, ok := e.sets[name]; ok {
					v.SetMeta(flags)
				}
			}
		}
		if s, ok := rv.Interface().(types.SecurityString); ok {
			v.Mask = s.SecurityString()
		}
		text, err := textx.MarshalText(rv)
		if err != nil {
			return err
		}
		v.Value = string(text)
		e.vars.Set(v)
		return nil
	}

	switch kind {
	case reflect.Ptr:
		if rv.IsNil() {
			return nil
		}
		return e.scan(w, rv.Elem())
	case reflect.Func, reflect.Interface, reflect.Chan, reflect.Map: // skip
	default:
		rt := rv.Type()
		if _, ok := rv.Interface().(types.TextMarshaler); ok {
			return setter(rv)
		}
		switch kind {
		case reflect.Array, reflect.Slice:
			for i := 0; i < rv.Len(); i++ {
				w.Enter(i)
				if err := e.scan(w, rv.Index(i)); err != nil {
					return err
				}
				w.Exit()
			}
		case reflect.Struct:
			for i := 0; i < rv.NumField(); i++ {
				var flags map[string]bool
				ft := rt.Field(i)
				name := ft.Name

				if !ast.IsExported(name) {
					continue
				}
				if tag, ok := ft.Tag.Lookup("env"); ok {
					key, _flags := reflectx.TagValueAndFlags(tag)
					if key == "-" {
						continue
					}
					if key != "" {
						name = key
					}
					flags = _flags
				}
				inline := flags == nil && ft.Anonymous &&
					reflectx.DeRef(ft.Type).Kind() == reflect.Struct
				if !inline {
					w.Enter(name)
				}
				if flags != nil {
					e.flag(w.String(), flags)
				}
				if err := e.scan(w, rv.Field(i)); err != nil {
					return err
				}
				if !inline {
					w.Exit()
				}
			}
		default:
			return setter(rv)
		}
	}
	return nil
}
