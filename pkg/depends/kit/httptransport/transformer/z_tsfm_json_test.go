package transformer_test

import (
	"bytes"
	"context"
	"net/http"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	. "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	vldterr "github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

type S string

func (s *S) UnmarshalText(data []byte) error {
	return errors.Errorf("err")
}

func TestJSON(t *testing.T) {
	data := struct {
		Data struct {
			S           S    `json:"s,omitempty"`
			Bool        bool `json:"bool"`
			StructSlice []struct {
				Name string `json:"name"`
			} `json:"structSlice"`
			StringSlice []string `json:"stringSlice"`
			NestedSlice []struct {
				Names []string `json:"names"`
			} `json:"nestedSlice"`
		} `json:"data"`
	}{}

	ct, _ := DefaultFactory.NewTransformer(
		bgctx,
		typesx.FromReflectType(reflect.TypeOf(data)), Option{},
	)

	t.Run("EncodeTo", func(t *testing.T) {
		b := bytes.NewBuffer(nil)
		h := http.Header{}

		err := ct.EncodeTo(context.Background(), WriterWithHeader(b, h), data)
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(h.Get(httpx.HeaderContentType)).
			To(Equal("application/json; charset=utf-8"))
	})

	t.Run("EncodeWithReflectValue", func(t *testing.T) {
		b := bytes.NewBuffer(nil)
		h := http.Header{}

		err := ct.EncodeTo(context.Background(), WriterWithHeader(b, h), reflect.ValueOf(data))
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(h.Get(httpx.HeaderContentType)).
			To(Equal("application/json; charset=utf-8"))
	})

	t.Run("DecodeAndValidate failed", func(t *testing.T) {
		b := bytes.NewBufferString(`{`)
		err := ct.DecodeFrom(context.Background(), b, &data)
		NewWithT(t).Expect(err).NotTo(BeNil())
	})

	t.Run("DecodeAndValidate success", func(t *testing.T) {
		b := bytes.NewBufferString(`{}`)
		err := ct.DecodeFrom(context.Background(), b, reflect.ValueOf(&data))
		NewWithT(t).Expect(err).To(BeNil())
	})

	t.Run("DecodeAndValidateFailedWithLocation", func(t *testing.T) {
		cases := []struct {
			json     string
			location string
		}{{
			`{
	"data": {
		"s": "111",
		"bool": true
	}
}`, "data.s",
		},
			{
				`
{
 	"data": {
		"bool": ""
	}
}
`, "data.bool",
			},
			{
				`
{
		"data": {
			"structSlice": [
				{"name":"{"},
				{"name":"1"},
				{"name": { "test": 1 }},
				{"name":"1"}
			]
		}
}`,
				"data.structSlice[2].name",
			},
			{
				`
		{
			"data": {
				"stringSlice":["1","2",3]
			}
		}`,
				"data.stringSlice[2]",
			},
			{
				`
		{
			"data": {
				"stringSlice":["1","2",3]
			}
		}`,
				"data.stringSlice[2]",
			},
			{
				`
		{
			"data": {
				"bool": true,
				"nestedSlice": [
					{ "names": ["1","2","3"] },
			        { "names": ["1","\"2", 3] }
				]
			}
		}
		`, "data.nestedSlice[1].names[2]",
			},
		}

		for _, c := range cases {
			b := bytes.NewBufferString(c.json)
			err := ct.DecodeFrom(context.Background(), b, &data)

			err.(*vldterr.ErrorSet).Each(func(fe *vldterr.FieldError) {
				NewWithT(t).Expect(fe.Field.String()).To(Equal(c.location))
			})
		}
	})
}
