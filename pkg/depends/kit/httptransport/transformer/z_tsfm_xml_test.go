package transformer_test

import (
	"bytes"
	"context"
	"net/http"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	. "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func TestXML(t *testing.T) {
	type TestData struct {
		Data struct {
			Bool        bool
			FirstName   string `xml:"name>first"`
			StructSlice []struct {
				Name string
			}
			StringSlice     []string
			StringAttrSlice []string `xml:"StringAttrSlice,attr"`
			NestedSlice     []struct {
				Names []string
			}
		}
	}

	data := TestData{}
	data.Data.FirstName = "test"
	data.Data.StringSlice = []string{"1", "2", "3"}
	data.Data.StringAttrSlice = []string{"1", "2", "3"}

	ct, _ := DefaultFactory.NewTransformer(
		bgctx,
		typesx.FromReflectType(reflect.TypeOf(data)),
		Option{MIME: "xml"},
	)

	t.Run("EncodeTo", func(t *testing.T) {
		t.Run("RawValue", func(t *testing.T) {
			b := bytes.NewBuffer(nil)
			h := http.Header{}
			err := ct.EncodeTo(bgctx, WriterWithHeader(b, h), data)
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(h.Get(httpx.HeaderContentType)).
				To(Equal("application/xml; charset=utf-8"))
		})

		t.Run("ReflectValue", func(t *testing.T) {
			b := bytes.NewBuffer(nil)
			h := http.Header{}
			err := ct.EncodeTo(bgctx, WriterWithHeader(b, h), reflect.ValueOf(data))
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(h.Get(httpx.HeaderContentType)).
				To(Equal("application/xml; charset=utf-8"))
		})
	})

	t.Run("DecodeAndValidate", func(t *testing.T) {
		t.Run("FAILED", func(t *testing.T) {
			b := bytes.NewBufferString("<")
			err := ct.DecodeFrom(context.Background(), b, &data)
			NewWithT(t).Expect(err).NotTo(BeNil())
		})

		t.Run("SUCCEED", func(t *testing.T) {
			b := bytes.NewBufferString("<TestData></TestData>")
			err := ct.DecodeFrom(context.Background(), b, reflect.ValueOf(&data))
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("FailedWithWrongType", func(t *testing.T) {
			b := bytes.NewBufferString("<TestData><Data><Bool>bool</Bool></Data></TestData>")
			err := ct.DecodeFrom(context.Background(), b, &data)
			NewWithT(t).Expect(err).NotTo(BeNil())
		})
	})
}
