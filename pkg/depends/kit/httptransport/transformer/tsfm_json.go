package transformer

import (
	"context"
	"encoding/json"
	"io"
	"net/textproto"
	"reflect"
	"strconv"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	vldterr "github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func init() { DefaultFactory.Register(&JSON{}) }

type JSON struct{}

func (JSON) Names() []string { return []string{httpx.MIME_JSON, "json"} }

func (JSON) NamedByTag() string { return "json" }

func (t *JSON) String() string { return httpx.MIME_JSON }

func (JSON) New(context.Context, typesx.Type) (Transformer, error) { return &JSON{}, nil }

func (t *JSON) EncodeTo(ctx context.Context, w io.Writer, v interface{}) error {
	if rv, ok := v.(reflect.Value); ok {
		v = rv.Interface()
	}

	httpx.MaybeWriteHeader(ctx, w, t.String(), map[string]string{
		"charset": "utf-8",
	})

	return json.NewEncoder(w).Encode(v)
}

func (JSON) DecodeFrom(ctx context.Context, r io.Reader, v interface{}, _ ...textproto.MIMEHeader) error {
	if rv, ok := v.(reflect.Value); ok {
		if rv.Kind() != reflect.Ptr && rv.CanAddr() {
			rv = rv.Addr()
		}
		v = rv.Interface()
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(v); err != nil {
		return WrapLocationDecoderError(dec, err)
	}
	return nil
}

func WrapLocationDecoderError(dec *json.Decoder, err error) error {
	switch e := err.(type) {
	case *json.UnmarshalTypeError:
		r := reflect.ValueOf(dec).Elem()
		errs := vldterr.NewErrorSet()
		errs.AddErr(e, location(r.Field(1 /* .buf */).Bytes(), int(e.Offset)))
		return errs.Err()
	case *json.SyntaxError:
		return e
	default:
		r := reflect.ValueOf(dec).Elem()
		// json.Decoder.d.off
		offset := r.Field(2).Field(1).Int()
		if offset > 0 {
			errs := vldterr.NewErrorSet()
			// json.Decoder.buf
			errs.AddErr(e, location(r.Field(1).Bytes(), int(offset-1)))
			return errs.Err()
		}
		return e
	}
}

func location(data []byte, offset int) string {
	i := 0
	arrayPaths := map[string]bool{}
	arrayIdxSet := map[string]int{}
	pw := &PathWalker{}

	markObjectKey := func() {
		jsonKey, l := nextString(data[i:])
		i += l

		if i < int(offset) && len(jsonKey) > 0 {
			key, _ := strconv.Unquote(string(jsonKey))
			pw.Enter(key)
		}
	}

	markArrayIdx := func(path string) {
		if arrayPaths[path] {
			arrayIdxSet[path]++
		} else {
			arrayPaths[path] = true
		}
		pw.Enter(arrayIdxSet[path])
	}

	for i < offset {
		i += nextToken(data[i:])
		char := data[i]

		switch char {
		case '"':
			_, l := nextString(data[i:])
			i += l
		case '[', '{':
			i++

			if char == '[' {
				markArrayIdx(pw.String())
			} else {
				markObjectKey()
			}
		case '}', ']', ',':
			i++
			pw.Exit()

			if char == ',' {
				path := pw.String()

				if _, ok := arrayPaths[path]; ok {
					markArrayIdx(path)
				} else {
					markObjectKey()
				}
			}
		default:
			i++
		}
	}

	return pw.String()
}

func nextToken(data []byte) int {
	for i, c := range data {
		switch c {
		case ' ', '\n', '\r', '\t':
			continue
		default:
			return i
		}
	}
	return -1
}

func nextString(data []byte) (finalData []byte, l int) {
	quoteStartAt := -1
	for i, c := range data {
		switch c {
		case '"':
			if i > 0 && string(data[i-1]) == "\\" {
				continue
			}
			if quoteStartAt >= 0 {
				return data[quoteStartAt : i+1], i + 1
			} else {
				quoteStartAt = i
			}
		default:
			continue
		}
	}
	return nil, 0
}
