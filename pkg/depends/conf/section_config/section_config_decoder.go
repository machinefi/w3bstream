package section_config

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"io"
	"os"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
	"github.com/machinefi/w3bstream/pkg/depends/x/textx"
)

type Decoder struct {
	Section
	Sep    byte
	Values map[string]string
}

func NewDecoder(sep byte) *Decoder {
	return &Decoder{
		Sep:    sep,
		Values: make(map[string]string),
	}
}

func (d *Decoder) Unmarshal(c SectionConfig, data []byte) error {
	rv := reflect.ValueOf(c)

	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		panic("input is nil")
	}

	rv = reflectx.Indirect(rv)

	if rv.Kind() != reflect.Struct {
		panic("input should be a struct")
	}

	scanner := bufio.NewReader(bytes.NewBuffer(data))

	var (
		line []byte
		err  error
	)
	for {
		line, _, err = scanner.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		if line[0] == '[' {
			if line[len(line)-1] != ']' {
				panic(fmt.Sprintf("unexpected section header: %s", string(line)))
			}
			line = line[1 : len(line)-1]
			parts := bytes.SplitN(line, []byte{':'}, 2)
			d.Name = string(parts[0])
			if len(parts) == 2 {
				d.Value = string(parts[1])
			}
			if d.Name != "" {
				c.SetSection(d.Name, d.Value)
			}
			break
		} else {
			panic("should meet section header at beginning")
		}
	}

	for {
		line, _, err = scanner.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		kv := bytes.SplitN(line, []byte{d.Sep}, 2)
		kv[0] = bytes.TrimSpace(kv[0])
		d.Values[string(kv[0])] = ""
		if len(kv) == 2 {
			d.Values[string(kv[0])] = string(bytes.TrimSpace(kv[1]))
		}
	}

	return d.scan(rv)
}

func (d *Decoder) UnmarshalFromFile(v SectionConfig, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return d.Unmarshal(v, data)
}

func (d *Decoder) scan(rv reflect.Value) error {
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
		return d.scan(rv.Elem())
	case reflect.Struct:
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			ft, fv := rt.Field(i), rv.Field(i)
			if !ast.IsExported(ft.Name) {
				continue
			}
			if ft.Anonymous {
				if err := d.scan(fv); err != nil {
					return err
				}
			}
			var (
				text []byte
				err  error
			)
			tag, ok := ft.Tag.Lookup("name")
			if !ok {
				continue
			}
			key, _ := reflectx.TagValueAndFlags(tag)
			if key == "-" {
				continue
			}
			text = []byte(d.Values[key])
			if len(text) == 0 {
				continue
			}
			if v, ok := fv.Addr().Interface().(types.TextUnmarshaler); ok {
				err = v.UnmarshalText(text)
			} else {
				err = textx.UnmarshalText(fv, text)
			}
			if err != nil {
				return err
			}
		}
	default:
		// skip
	}
	return nil
}
