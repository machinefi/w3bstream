package transformer_test

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"reflect"
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	. "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func TestMultipart(t *testing.T) {
	parts := `--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="PtrBool"
Content-Type: text/plain; charset=utf-8

true
--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="PtrInt"
Content-Type: text/plain; charset=utf-8

1
--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="Bool"
Content-Type: text/plain; charset=utf-8

true
--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="bytes"
Content-Type: text/plain; charset=utf-8

Ynl0ZXM=
--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="first_name"
Content-Type: text/plain; charset=utf-8

test
--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="StructSlice"
Content-Type: application/json; charset=utf-8

{"Name":"name"}

--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="StringSlice"
Content-Type: text/plain; charset=utf-8

1
--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="StringSlice"
Content-Type: text/plain; charset=utf-8

2
--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="StringSlice"
Content-Type: text/plain; charset=utf-8

3
--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="StringArray"
Content-Type: text/plain; charset=utf-8

1
--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="StringArray"
Content-Type: text/plain; charset=utf-8


--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="StringArray"
Content-Type: text/plain; charset=utf-8

3
--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="Struct"
Content-Type: application/xml; charset=utf-8

<Sub><Name></Name></Sub>
--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="Files"; filename="file0.txt"
Content-Type: application/octet-stream

text0
--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="Files"; filename="file1.txt"
Content-Type: application/octet-stream

text1
--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6
Content-Disposition: form-data; name="File"; filename="file.txt"
Content-Type: application/octet-stream

text
--99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6--`

	type Sub struct {
		Name string
	}

	type TestData struct {
		PtrBoolEmpty *bool `name:",omitempty"`
		PtrBool      *bool `name:",omitempty"`
		PtrInt       *int
		Bool         bool
		Bytes        []byte `name:"bytes"`
		FirstName    string `name:"first_name,omitempty"`
		StructSlice  []Sub
		StringSlice  []string
		StringArray  [3]string
		Struct       Sub                     `mime:"xml"`
		Files        []*multipart.FileHeader `name:",omitempty"`
		File         *multipart.FileHeader   `name:",omitempty"`
	}

	data := TestData{}
	data.PtrBool = ptrx.Bool(true)
	data.FirstName = "test"
	data.Bool = true
	data.Bytes = []byte("bytes")
	data.PtrInt = ptrx.Int(1)
	data.StringSlice = []string{"1", "2", "3"}
	data.StructSlice = []Sub{{Name: "name"}}
	data.StringArray = [3]string{"1", "", "3"}

	data.File, _ = NewFileHeader("File", "file.txt", bytes.NewBufferString("text"))

	data.Files = []*multipart.FileHeader{
		MustNewFileHeader("Files", "file0.txt", bytes.NewBufferString("text0")),
		MustNewFileHeader("Files", "file1.txt", bytes.NewBufferString("text1")),
	}

	ct, _ := DefaultFactory.NewTransformer(
		bgctx,
		typesx.FromReflectType(reflect.TypeOf(data)),
		Option{MIME: "multipart"},
	)

	t.Run("EncodeTo", func(t *testing.T) {
		b := bytes.NewBuffer(nil)
		h := http.Header{}
		err := ct.EncodeTo(bgctx, WriterWithHeader(b, h), data)
		NewWithT(t).Expect(err).To(BeNil())
		_, params, _ := mime.ParseMediaType(h.Get(httpx.HeaderContentType))

		gen := toParts(b, params["boundary"])
		expect := toParts(bytes.NewBufferString(replaceBoundaryMultipart(parts, params["boundary"])), params["boundary"])

		NewWithT(t).Expect(len(gen)).To(Equal(len(expect)))

		for i := range gen {
			NewWithT(t).Expect(gen[i].FormName()).To(Equal(expect[i].FormName()))
			NewWithT(t).Expect(gen[i].FileName()).To(Equal(expect[i].FileName()))
			NewWithT(t).Expect(gen[i].Header).To(Equal(expect[i].Header))
		}
	})

	t.Run("DecodeAndValidate", func(t *testing.T) {
		b := bytes.NewBufferString(parts)
		td := TestData{}

		err := ct.DecodeFrom(bgctx, b, &td, textproto.MIMEHeader{
			httpx.HeaderContentType: []string{
				mime.FormatMediaType(
					ct.Names()[0],
					map[string]string{"boundary": boundary},
				),
			},
		})

		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(td).To(Equal(data))
	})
}

var boundary = "99bb5d156e61cf661d01fc370479b62a3451759d25d14711fd7e9db170f6"

func replaceBoundaryMultipart(data string, generatedBoundary string) string {
	return strings.Replace(data, boundary, generatedBoundary, -1)
}

func toParts(r io.Reader, b string) (parts []*multipart.Part) {
	gen := multipart.NewReader(r, b)

	for {
		rp, err := gen.NextPart()
		if err != nil {
			break
		}
		data := bytes.NewBuffer(nil)
		_, _ = io.Copy(data, rp)
		rp.Header["Content"] = []string{data.String()}
		_ = rp.Close()
		parts = append(parts, rp)
	}

	return
}
