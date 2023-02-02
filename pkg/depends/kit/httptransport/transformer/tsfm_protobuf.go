package transformer

import (
	"context"
	"io"
	"net/textproto"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/validator"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

var (
	ErrEncodeDataNotProtobuf   = errors.New("encode data must be `proto.Message`")
	ErrDecodeTargetNotProtobuf = errors.New("decode target must be `proto.Message`")
)

func init() { DefaultFactory.Register(&Protobuf{}) }

type Protobuf struct{}

func (Protobuf) Names() []string { return []string{httpx.MIME_PROTOBUF, "protobuf", "x-protobuf"} }

func (Protobuf) NamedByTag() string { return "protobuf" }

func (t *Protobuf) String() string { return httpx.MIME_PROTOBUF }

func (Protobuf) New(context.Context, typesx.Type) (Transformer, error) { return &Protobuf{}, nil }

func (t *Protobuf) EncodeTo(ctx context.Context, w io.Writer, v interface{}) error {
	if rv, ok := v.(reflect.Value); ok {
		v = rv.Interface()
	}
	httpx.MaybeWriteHeader(ctx, w, t.String(), map[string]string{})
	pv, ok := v.(proto.Message)
	if !ok {
		return ErrEncodeDataNotProtobuf
	}
	data, err := proto.Marshal(pv)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (t *Protobuf) DecodeFrom(ctx context.Context, r io.Reader, v interface{}, _ ...textproto.MIMEHeader) error {
	if rv, ok := v.(reflect.Value); ok {
		if rv.Kind() != reflect.Ptr && rv.CanAddr() {
			rv = rv.Addr()
		}
		v = rv.Interface()
	}

	pv, ok := v.(proto.Message)
	if !ok {
		return ErrDecodeTargetNotProtobuf
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return proto.Unmarshal(data, pv)
}

// NewValidator returns empty validator to implements interface `MayValidate` to skip protobuf struct validation
func (t *Protobuf) NewValidator(_ context.Context, _ typesx.Type) (validator.Validator, error) {
	return nil, nil
}
