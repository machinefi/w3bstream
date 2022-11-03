package textx

import (
	"bytes"
	"encoding"
	"encoding/base64"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

func MarshalText(v interface{}, useBase64 ...bool) ([]byte, error) {
	if rv, ok := v.(reflect.Value); ok {
		for rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return nil, nil
			}
			rv = rv.Elem()
		}

		if rv.CanInterface() {
			v = rv.Interface()
		}
	}

	if marshaler, ok := v.(encoding.TextMarshaler); ok {
		return marshaler.MarshalText()
	}

	if v == nil {
		return nil, nil
	}

	switch x := v.(type) {
	case []byte:
		if len(useBase64) > 0 && useBase64[0] {
			return ToBase64(x), nil
		}
		return x, nil
	case string:
		return []byte(x), nil
	case bool:
		return strconv.AppendBool([]byte{}, x), nil
	case int:
		return strconv.AppendInt([]byte{}, int64(x), 10), nil
	case int8:
		return strconv.AppendInt([]byte{}, int64(x), 10), nil
	case int16:
		return strconv.AppendInt([]byte{}, int64(x), 10), nil
	case int32:
		return strconv.AppendInt([]byte{}, int64(x), 10), nil
	case int64:
		return strconv.AppendInt([]byte{}, x, 10), nil
	case uint:
		return strconv.AppendUint([]byte{}, uint64(x), 10), nil
	case uint8:
		return strconv.AppendUint([]byte{}, uint64(x), 10), nil
	case uint16:
		return strconv.AppendUint([]byte{}, uint64(x), 10), nil
	case uint32:
		return strconv.AppendUint([]byte{}, uint64(x), 10), nil
	case uint64:
		return strconv.AppendUint([]byte{}, x, 10), nil
	case float32:
		return strconv.AppendFloat([]byte{}, float64(x), 'g', -1, 32), nil
	case float64:
		return strconv.AppendFloat([]byte{}, x, 'g', -1, 64), nil
	default:
		rv := reflect.ValueOf(x)

		for rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return nil, nil
			}
			rv = rv.Elem()
		}

		switch rv.Kind() {
		case reflect.Slice:
			elem := rv.Type().Elem()
			if elem.Kind() == reflect.Uint8 && elem.PkgPath() == "" {
				if len(useBase64) > 0 && useBase64[0] {
					return ToBase64(rv.Bytes()), nil
				}
				return rv.Bytes(), nil
			}
		case reflect.String:
			return []byte(rv.String()), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return strconv.AppendInt([]byte{}, rv.Int(), 10), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return strconv.AppendUint([]byte{}, rv.Uint(), 10), nil
		case reflect.Float32:
			return strconv.AppendFloat([]byte{}, rv.Float(), 'g', -1, 32), nil
		case reflect.Float64:
			return strconv.AppendFloat([]byte{}, rv.Float(), 'g', -1, 64), nil
		case reflect.Bool:
			return strconv.AppendBool([]byte{}, rv.Bool()), nil
		}
		return nil, errors.Errorf("unsupported type %T", x)
	}
}

func UnmarshalText(v interface{}, data []byte, useBase64 ...bool) error {
	if rv, ok := v.(reflect.Value); ok {
		if rv.Kind() != reflect.Ptr {
			rv = rv.Addr()
		} else {
			if rv.IsNil() {
				rv.Set(NewReflectValue(rv.Type()))
			}
		}

		if rv.CanInterface() {
			if unmarshaler, ok := rv.Interface().(encoding.TextUnmarshaler); ok {
				if err := unmarshaler.UnmarshalText(data); err != nil {
					return errors.Wrapf(err, "unmarshal text to %T failed", v)
				}
				return nil
			}
		}

		return UnmarshalTextToReflectValue(rv, data)
	}

	if unmarshaler, ok := v.(encoding.TextUnmarshaler); ok {
		if err := unmarshaler.UnmarshalText(data); err != nil {
			return errors.Wrapf(err, "unmarshal text to %T failed", v)
		}
		return nil
	}

	if v == nil {
		return UnmarshalText(reflect.ValueOf(v), data)
	}

	switch x := v.(type) {
	case *[]byte:
		var d []byte
		if len(useBase64) > 0 && useBase64[0] {
			var err error
			d, err = FromBase64(data)
			if err != nil {
				return err
			}
		} else {
			d = make([]byte, len(data))
			copy(d, data)
		}
		*x = d
	case *string:
		*x = string(data)
	case *bool:
		v, err := strconv.ParseBool(string(data))
		if err != nil {
			return errors.Wrapf(err, "unmarshal text")
		}
		*x = v
	case *int:
		i, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		*x = int(i)
	case *int8:
		i, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		*x = int8(i)
	case *int16:
		i, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		*x = int16(i)
	case *int32:
		i, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		*x = int32(i)
	case *int64:
		i, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		*x = i
	case *uint:
		i, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		*x = uint(i)
	case *uint8:
		i, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		*x = uint8(i)
	case *uint16:
		i, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		*x = uint16(i)
	case *uint32:
		i, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		*x = uint32(i)
	case *uint64:
		i, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		*x = i
	case *float32:
		i, err := strconv.ParseFloat(string(data), 32)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		*x = float32(i)
	case *float64:
		i, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		*x = i
	default:
		return UnmarshalTextToReflectValue(reflect.ValueOf(x), data, useBase64...)
	}
	return nil
}

func UnmarshalTextToReflectValue(rv reflect.Value, data []byte, useBase64 ...bool) error {
	if rv.Kind() != reflect.Ptr {
		return errors.Errorf("unmarshal text need ptr value, but got %#v", rv.Interface())
	}

	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(NewReflectValue(rv.Type()))
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Slice:
		elem := rv.Type().Elem()
		if elem.Kind() == reflect.Uint8 && elem.PkgPath() == "" {
			if len(useBase64) > 0 && useBase64[0] {
				d, err := FromBase64(data)
				if err != nil {
					return err
				}
				rv.SetBytes(d)
				return nil
			}
			rv.SetBytes(data)
		}
	case reflect.String:
		rv.SetString(string(data))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intV, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		rv.SetInt(intV)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintV, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		rv.SetUint(uintV)
	case reflect.Float32, reflect.Float64:
		floatV, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		rv.SetFloat(floatV)
	case reflect.Bool:
		boolV, err := strconv.ParseBool(string(data))
		if err != nil {
			return errors.Wrap(err, "unmarshal text")
		}
		rv.SetBool(boolV)
	}
	return nil
}

func NewReflectValue(t reflect.Type) reflect.Value {
	v := reflect.New(t).Elem()
	if t.Kind() == reflect.Ptr {
		v.Set(NewReflectValue(t.Elem()).Addr())
	}
	return v
}

func ToBase64(raw []byte) []byte {
	length := base64.StdEncoding.EncodedLen(len(raw))
	if length <= 1024 {
		d := make([]byte, length)
		base64.StdEncoding.Encode(d, raw)
		return d
	}
	b := bytes.NewBuffer(nil)
	base64.NewDecoder(base64.StdEncoding, b)
	return b.Bytes()
}

func FromBase64(data []byte) ([]byte, error) {
	length := base64.StdEncoding.DecodedLen(len(data))
	d := make([]byte, length)
	n, err := base64.StdEncoding.Decode(d, data)
	if err != nil {
		return nil, err
	}
	return d[:n], nil
}
