package transformer

import (
	"context"
	"io"
	"reflect"

	pkgerr "github.com/pkg/errors"

	vldterr "github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
)

func NewSuper(tsfm Transformer, opt *CommonOption) *Super {
	return &Super{
		tsfm:         tsfm,
		CommonOption: *opt,
	}
}

type Super struct {
	tsfm Transformer
	CommonOption
}

func (t *Super) EncodeTo(ctx context.Context, w io.Writer, v interface{}) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	if t.Explode {
		rv = reflectx.Indirect(rv)
		// for create slice
		if setter, ok := w.(interface{ SetN(n int) }); ok {
			setter.SetN(rv.Len())
		}

		if next, ok := w.(CanNextWriter); ok {
			errs := vldterr.NewErrorSet()
			for i := 0; i < rv.Len(); i++ {
				nw := next.NextWriter()
				if err := t.tsfm.EncodeTo(ctx, nw, rv.Index(i)); err != nil {
					errs.AddErr(err, i)
				}
			}
			return errs.Err()
		}
		return nil
	}

	// should skip empty value when omitempty
	if !(t.Omitempty && reflectx.IsEmptyValue(rv)) {
		writer, ok := w.(CanNextWriter)
		if ok {
			return t.tsfm.EncodeTo(ctx, writer.NextWriter(), rv)
		}
		return t.tsfm.EncodeTo(ctx, w, rv)
	}

	return nil
}

func (t *Super) DecodeFrom(ctx context.Context, r io.Reader, v interface{}) error {
	if rv, ok := v.(reflect.Value); ok {
		v = rv.Interface()
	}

	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return pkgerr.Errorf("decode target must be ptr value")
	}

	if t.Explode {
		valc := 0
		if with, ok := r.(interface{ Len() int }); ok {
			valc = with.Len()
		}
		if valc == 0 {
			return nil
		}
		if x, ok := v.(*[]string); ok {
			if with, ok := r.(CanInterface); ok {
				if values, ok := with.Interface().([]string); ok {
					*x = values
					return nil
				}
			}
		}

		rv := reflectx.Indirect(reflect.ValueOf(v))

		// make slice (ignore array)
		if rv.Kind() == reflect.Slice {
			rv.Set(reflect.MakeSlice(rv.Type(), valc, valc))
		}

		reader, ok := r.(CanNextReader)
		if !ok {
			return nil
		}

		errs := vldterr.NewErrorSet()
		// ignore when values length greater than array len
		for i := 0; i < rv.Len() && i < valc; i++ {
			if err := t.tsfm.DecodeFrom(
				ctx,
				reader.NextReader(),
				rv.Index(i).Addr(),
			); err != nil {
				errs.AddErr(err, i)
			}
		}
		return errs.Err()
	}

	reader, ok := r.(CanNextReader)
	if ok {
		return t.tsfm.DecodeFrom(ctx, reader.NextReader(), v)
	}
	return t.tsfm.DecodeFrom(ctx, r, v)
}

type CanNextWriter interface {
	NextWriter() io.Writer
}

type CanNextReader interface {
	NextReader() io.Reader
}
