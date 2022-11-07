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
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func TestURLEncoded(t *testing.T) {
	queries := `Bool=true` +
		`&Bytes=Ynl0ZXM%3D` +
		`&PtrInt=1` +
		`&StringArray=1&StringArray=&StringArray=3` +
		`&StringSlice=1&StringSlice=2&StringSlice=3` +
		`&Struct=%3CSub%3E%3CName%3E%3C%2FName%3E%3C%2FSub%3E` +
		`&StructSlice=%7B%22Name%22%3A%22name%22%7D%0A` +
		`&first_name=test`

	type Sub struct{ Name string }

	type TestData struct {
		PtrBool     *bool `name:",omitempty"`
		PtrInt      *int
		Bool        bool
		Bytes       []byte
		FirstName   string `name:"first_name"`
		StructSlice []Sub
		StringSlice []string
		StringArray [3]string
		Struct      Sub `mime:"xml"`
	}

	data := TestData{}
	data.FirstName = "test"
	data.Bool = true
	data.Bytes = []byte("bytes")
	data.PtrInt = ptrx.Int(1)
	data.StringSlice = []string{"1", "2", "3"}
	data.StructSlice = []Sub{{Name: "name"}}
	data.StringArray = [3]string{"1", "", "3"}

	tsf, _ := DefaultFactory.NewTransformer(
		bgctx,
		typesx.FromReflectType(reflect.TypeOf(data)),
		Option{MIME: "urlencoded"},
	)

	t.Run("EncodeTo", func(t *testing.T) {
		b := bytes.NewBuffer(nil)
		h := http.Header{}

		err := tsf.EncodeTo(bgctx, WriterWithHeader(b, h), data)

		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(h.Get(httpx.HeaderContentType)).
			To(Equal("application/x-www-form-urlencoded; param=value"))
		NewWithT(t).Expect(b.String()).To(Equal(queries))
	})

	t.Run("DecodeAndValidate", func(t *testing.T) {
		b := bytes.NewBufferString(queries)
		td := TestData{}

		err := tsf.DecodeFrom(context.Background(), b, &td)
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(td).To(Equal(data))
	})
}
