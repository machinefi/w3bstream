package env

import (
	"go/ast"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
	"github.com/machinefi/w3bstream/pkg/depends/x/textx"
)

type Decoder struct{ vars *Vars }

func NewDecoder(vs *Vars) *Decoder { return &Decoder{vars: vs} }

func (d *Decoder) Decode(v interface{}) error {
	w := NewPathWalker()

	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	return d.scan(w, rv)
}

// scan scan and set value
func (d *Decoder) scan(w *PathWalker, rv reflect.Value) error {
	kind := rv.Kind()

	if kind != reflect.Ptr && rv.CanAddr() {
		if v, ok := rv.Addr().Interface().(types.DefaultSetter); ok {
			v.SetDefault()
		}
	}

	switch kind {
	case reflect.Ptr:
		if rv.IsNil() {
			rv.Set(reflectx.New(rv.Type()))
		}
		return d.scan(w, rv.Elem())
	case reflect.Func, reflect.Interface, reflect.Chan, reflect.Map:
		// skip
	default:
		rt := rv.Type()
		if rt.Implements(types.RTypeTextUnmarshaler) ||
			reflect.PtrTo(rt).Implements(types.RTypeTextUnmarshaler) {
			if v := d.vars.Get(w.String()); v != nil && v.Value != "" {
				return textx.UnmarshalText(rv, []byte(v.Value))
			}
			return nil
		}
		switch kind {
		case reflect.Array, reflect.Slice:
			size := d.vars.Len(w.String())
			if kind == reflect.Slice && rv.IsNil() && size > 0 {
				rv.Set(reflect.MakeSlice(rv.Type(), size, size))
			}
			for i := 0; i < rv.Len(); i++ {
				w.Enter(i)
				if err := d.scan(w, rv.Index(i)); err != nil {
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
				if err := d.scan(w, rv.Field(i)); err != nil {
					return err
				}
				if !inline {
					w.Exit()
				}
			}
		default:
			if v := d.vars.Get(w.String()); v != nil && v.Value != "" {
				return textx.UnmarshalText(rv, []byte(v.Value))
			}
		}
	}
	return nil
}
