package section_config

import (
	"bytes"
	"go/ast"
	"os"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
	"github.com/machinefi/w3bstream/pkg/depends/x/textx"
)

type Encoder struct {
	Sep    byte
	Values map[string]string
}

func NewEncoder(sep byte) *Encoder { return &Encoder{Sep: sep} }

func (e *Encoder) Marshal(c SectionConfig) ([]byte, error) {
	rv := reflectx.Indirect(reflect.ValueOf(c))

	if rv.Kind() != reflect.Struct {
		panic("input value should be a struct")
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteString(c.GetSection().String())
	buf.WriteRune('\n')

	if err := e.scan(rv, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *Encoder) MarshalToFile(c SectionConfig, path string) error {
	data, err := e.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, os.ModePerm)
}

func (e *Encoder) scan(rv reflect.Value, buf *bytes.Buffer) error {
	kind := rv.Kind()

	switch kind {
	case reflect.Ptr:
		if rv.IsNil() {
			return nil
		}
		return e.scan(rv.Elem(), buf)
	case reflect.Struct:
		rt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			ft, fv := rt.Field(i), rv.Field(i)
			if !ast.IsExported(ft.Name) {
				continue
			}
			if ft.Anonymous {
				if err := e.scan(fv, buf); err != nil {
					return err
				}
			}
			tag, ok := ft.Tag.Lookup("name")
			if !ok {
				continue
			}
			key, _ := reflectx.TagValueAndFlags(tag)
			var (
				val []byte
				err error
			)
			if v, ok := fv.Interface().(types.TextMarshaler); ok {
				val, err = v.MarshalText()
			} else {
				val, err = textx.MarshalText(fv.Interface())
			}
			if err != nil {
				return err
			}
			buf.Write([]byte(key))
			buf.WriteByte(e.Sep)
			buf.Write(val)
			buf.WriteByte('\n')
		}
	default:
		// skip
	}
	return nil
}
