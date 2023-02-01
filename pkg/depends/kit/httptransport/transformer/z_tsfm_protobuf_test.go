package transformer_test

import (
	"bytes"
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	. "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	testdatapb "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer/testdata"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func TestProtobuf(t *testing.T) {
	data := testdatapb.Event{}

	ct, _ := DefaultFactory.NewTransformer(
		bgctx,
		typesx.FromReflectType(reflect.TypeOf(data)),
		Option{MIME: "protobuf"},
	)

	t.Run("EncodeTo", func(t *testing.T) {
		b := bytes.NewBuffer(nil)
		h := http.Header{}

		err := ct.EncodeTo(context.Background(), WriterWithHeader(b, h), &data)
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(h.Get(httpx.HeaderContentType)).To(Equal("application/x-protobuf"))

		err = ct.EncodeTo(context.Background(), WriterWithHeader(b, h), struct{}{})
		NewWithT(t).Expect(err).To(Equal(ErrEncodeDataNotProtobuf))

		err = ct.EncodeTo(context.Background(), WriterWithHeader(b, h), data)
		NewWithT(t).Expect(err).To(Equal(ErrEncodeDataNotProtobuf))
	})

	t.Run("DecodeFrom", func(t *testing.T) {
		ori := testdatapb.Event{
			Header:  &testdatapb.Header{EventType: "mock_type"},
			Payload: []byte(`{"key":"value"}`),
		}

		raw, err := proto.Marshal(&ori)
		NewWithT(t).Expect(err).To(BeNil())

		b := bytes.NewBuffer(raw)
		decoded := testdatapb.Event{}
		err = ct.DecodeFrom(context.Background(), b, &decoded)
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(decoded.Header.EventType).To(Equal(ori.Header.EventType))
		NewWithT(t).Expect(decoded.Payload).To(Equal(ori.Payload))

		err = ct.DecodeFrom(context.Background(), b, decoded)
		NewWithT(t).Expect(err).To(Equal(ErrDecodeTargetNotProtobuf))

		err = ct.DecodeFrom(context.Background(), b, struct{}{})
		NewWithT(t).Expect(err).To(Equal(ErrDecodeTargetNotProtobuf))
	})

}
