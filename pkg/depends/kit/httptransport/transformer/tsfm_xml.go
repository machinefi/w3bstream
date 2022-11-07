package transformer

import (
	"context"
	"encoding/xml"
	"io"
	"net/textproto"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func init() { DefaultFactory.Register(&XML{}) }

type XML struct{}

func (XML) Names() []string { return []string{httpx.MIME_XML, "xml"} }

func (t *XML) String() string { return httpx.MIME_XML }

func (XML) NamedByTag() string { return "xml" }

func (XML) New(context.Context, typesx.Type) (Transformer, error) {
	return &XML{}, nil
}

func (t *XML) EncodeTo(ctx context.Context, w io.Writer, v interface{}) error {
	if rv, ok := v.(reflect.Value); ok {
		v = rv.Interface()
	}

	httpx.MaybeWriteHeader(ctx, w, t.String(), map[string]string{
		"charset": "utf-8",
	})

	return xml.NewEncoder(w).Encode(v)
}

func (XML) DecodeFrom(_ context.Context, r io.Reader, v interface{}, _ ...textproto.MIMEHeader) error {
	if rv, ok := v.(reflect.Value); ok {
		if rv.Kind() != reflect.Ptr && rv.CanAddr() {
			rv = rv.Addr()
		}
		v = rv.Interface()
	}
	d := xml.NewDecoder(r)
	err := d.Decode(v)
	if err != nil {
		// TODO resolve field path by InputOffset()
		return err
	}
	return nil
}
