package transformer

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/textproto"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func init() { DefaultFactory.Register(&OctetStream{}) }

type OctetStream struct{}

func (t *OctetStream) String() string { return httpx.MIME_OCTET_STREAM }

func (OctetStream) Names() []string {
	return []string{httpx.MIME_OCTET_STREAM, "stream", "octet-stream"}
}

func (OctetStream) New(context.Context, typesx.Type) (Transformer, error) { return &OctetStream{}, nil }

func (t *OctetStream) EncodeTo(ctx context.Context, w io.Writer, v interface{}) error {
	rv, ok := v.(reflect.Value)
	if ok {
		v = rv.Interface()
	}

	switch x := v.(type) {
	case io.Reader:
		httpx.MaybeWriteHeader(ctx, w, t.Names()[0], nil)
		if _, err := io.Copy(w, x); err != nil {
			return err
		}
	case bytes.Buffer:
		// TODO draft for fix bytes.Buffer
		httpx.MaybeWriteHeader(ctx, w, t.Names()[0], nil)
		if _, err := io.Copy(w, &x); err != nil {
			return err
		}
	case *multipart.FileHeader:
		file, err := x.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		if rw, ok := w.(httpx.WithHeader); ok {
			for k := range x.Header {
				rw.Header()[k] = x.Header[k]
			}
		}

		if _, err := io.Copy(w, file); err != nil {
			return err
		}
	}

	return nil
}

func (OctetStream) DecodeFrom(_ context.Context, r io.Reader, v interface{}, _ ...textproto.MIMEHeader) error {
	rv, ok := v.(reflect.Value)
	if ok {
		v = rv.Interface()
	}

	switch x := v.(type) {
	case io.Writer:
		if _, err := io.Copy(x, r); err != nil {
			return err
		}
	case *multipart.FileHeader:
		if with, ok := r.(CanInterface); ok {
			if fh, ok := with.Interface().(*multipart.FileHeader); ok {
				*x = *fh
			}
		}
	case **multipart.FileHeader:
		if with, ok := r.(CanInterface); ok {
			if fh, ok := with.Interface().(*multipart.FileHeader); ok {
				*x = fh
			}
		}
	}

	return nil
}
