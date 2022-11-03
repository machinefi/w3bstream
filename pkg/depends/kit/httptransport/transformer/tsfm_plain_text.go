package transformer

import (
	"context"
	"io"
	"net/textproto"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/x/textx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func init() { DefaultFactory.Register(&PlainText{}) }

type PlainText struct{}

func (t *PlainText) String() string { return httpx.MIME_PLAIN_TEXT }

func (PlainText) Names() []string {
	return []string{httpx.MIME_PLAIN_TEXT, "plain", "text", "txt"}
}

func (PlainText) New(context.Context, typesx.Type) (Transformer, error) { return &PlainText{}, nil }

func (t *PlainText) EncodeTo(ctx context.Context, w io.Writer, v interface{}) error {
	httpx.MaybeWriteHeader(ctx, w, t.String(), map[string]string{
		"charset": "utf-8",
	})

	data, err := textx.MarshalText(v, true)
	if err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

func (t *PlainText) DecodeFrom(_ context.Context, r io.Reader, v interface{}, _ ...textproto.MIMEHeader) error {
	switch x := r.(type) {
	case CanString:
		raw := x.String()
		if x, ok := v.(*string); ok {
			*x = raw
			return nil
		}
		return textx.UnmarshalText(v, []byte(raw), true)
	case CanInterface:
		if raw, ok := x.Interface().(string); ok {
			if x, ok := v.(*string); ok {
				*x = raw
				return nil
			}
			return textx.UnmarshalText(v, []byte(raw), true)
		}
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return textx.UnmarshalText(v, data, true)
}
