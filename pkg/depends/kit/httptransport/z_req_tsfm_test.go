package httptransport_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	pkgerr "github.com/pkg/errors"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	vldterr "github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/testutil/httptransporttestutil/server/pkg/types"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
)

var regexpContentTypeWithBoundary = regexp.MustCompile(`Content-Type: multipart/form-data; boundary=([A-Za-z0-9]+)`)

func UnifyRequestData(data []byte) []byte {
	data = bytes.Replace(data, []byte("\r\n"), []byte("\n"), -1)
	if regexpContentTypeWithBoundary.Match(data) {
		matches := regexpContentTypeWithBoundary.FindAllSubmatch(data, 1)
		data = bytes.Replace(data, matches[0][1], []byte("boundary1"), -1)
	}
	return data
}

// openapi:strfmt date-time
type Datetime time.Time

func (dt Datetime) IsZero() bool {
	unix := time.Time(dt).Unix()
	return unix == 0 || unix == (time.Time{}).Unix()
}

func (dt Datetime) MarshalText() ([]byte, error) {
	str := time.Time(dt).Format(time.RFC3339)
	return []byte(str), nil
}

func (dt *Datetime) UnmarshalText(data []byte) error {
	if len(data) != 0 {
		return nil
	}
	t, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return err
	}
	*dt = Datetime(t)
	return nil
}

func TestRequestTsfm(t *testing.T) {
	factory := NewRequestTsfmFactory(nil, nil)

	type Headers struct {
		HInt    int    `in:"header"`
		HString string `in:"header"`
		HBool   bool   `in:"header"`
	}

	type Queries struct {
		QInt            int       `name:"int"                 in:"query"`
		QEmptyInt       int       `name:"emptyInt,omitempty"  in:"query"`
		QString         string    `name:"string"              in:"query"`
		QSlice          []string  `name:"slice"               in:"query"`
		QBytes          []byte    `name:"bytes,omitempty"     in:"query"`
		StartedAt       *Datetime `name:"startedAt,omitempty" in:"query"`
		QBytesOmitEmpty []byte    `name:"bytesOmit,omitempty" in:"query"`
	}

	type Cookies struct {
		CString string   `name:"a"     in:"cookie"`
		CSlice  []string `name:"slice" in:"cookie"`
	}

	type Data struct {
		A string `json:",omitempty" xml:",omitempty"`
		B string `json:",omitempty" xml:",omitempty"`
		C string `json:",omitempty" xml:",omitempty"`
	}

	type FormDataMultipart struct {
		Bytes []byte `name:"bytes"`
		A     []int  `name:"a"`
		C     uint   `name:"c" `
		Data  Data   `name:"data"`

		File  *multipart.FileHeader   `name:"file"`
		Files []*multipart.FileHeader `name:"files"`
	}

	cases := []struct {
		name   string
		path   string
		expect string
		req    interface{}
	}{
		{
			"FullInParameters",
			"/:id",
			`GET /1?bytes=Ynl0ZXM%3D&int=1&slice=1&slice=2&string=string HTTP/1.1
Content-Type: application/json; charset=utf-8
Cookie: a=xxx; slice=1; slice=2
Hbool: true
Hint: 1
Hstring: string

{}
`,
			&struct {
				Headers
				Queries
				Cookies
				Data `in:"body"`
				ID   string `name:"id" in:"path"`
			}{
				ID: "1",
				Headers: Headers{
					HInt:    1,
					HString: "string",
					HBool:   true,
				},
				Queries: Queries{
					QInt:    1,
					QString: "string",
					QSlice:  []string{"1", "2"},
					QBytes:  []byte("bytes"),
				},
				Cookies: Cookies{
					CString: "xxx",
					CSlice:  []string{"1", "2"},
				},
			},
		},
		{
			"URLEncoded",
			"/",
			`GET / HTTP/1.1
Content-Type: application/x-www-form-urlencoded; param=value

int=1&slice=1&slice=2&string=string`,
			&struct {
				Queries `in:"body" mime:"urlencoded"`
			}{
				Queries: Queries{
					QInt:    1,
					QString: "string",
					QSlice:  []string{"1", "2"},
				},
			},
		},
		{
			"XML",
			"/",
			`GET / HTTP/1.1
Content-Type: application/xml; charset=utf-8

<Data><A>1</A></Data>`,
			&struct {
				Data `in:"body" mime:"xml"`
			}{
				Data: Data{
					A: "1",
				},
			},
		},
		{
			"form-data/multipart",
			"/",
			`GET / HTTP/1.1
Content-Type: multipart/form-data; boundary=5eaf397248958ac38281d1c034e1ad0d4a5f7d986d4c53ac32e8399cbcda

--5eaf397248958ac38281d1c034e1ad0d4a5f7d986d4c53ac32e8399cbcda
Content-Disposition: form-data; name="bytes"
Content-Type: text/plain; charset=utf-8

Ynl0ZXM=
--5eaf397248958ac38281d1c034e1ad0d4a5f7d986d4c53ac32e8399cbcda
Content-Disposition: form-data; name="a"
Content-Type: text/plain; charset=utf-8

-1
--5eaf397248958ac38281d1c034e1ad0d4a5f7d986d4c53ac32e8399cbcda
Content-Disposition: form-data; name="a"
Content-Type: text/plain; charset=utf-8

1
--5eaf397248958ac38281d1c034e1ad0d4a5f7d986d4c53ac32e8399cbcda
Content-Disposition: form-data; name="c"
Content-Type: text/plain; charset=utf-8

1
--5eaf397248958ac38281d1c034e1ad0d4a5f7d986d4c53ac32e8399cbcda
Content-Disposition: form-data; name="data"
Content-Type: application/json; charset=utf-8

{"A":"1"}

--5eaf397248958ac38281d1c034e1ad0d4a5f7d986d4c53ac32e8399cbcda
Content-Disposition: form-data; name="file"; filename="file.text"
Content-Type: application/octet-stream

test
--5eaf397248958ac38281d1c034e1ad0d4a5f7d986d4c53ac32e8399cbcda
Content-Disposition: form-data; name="files"; filename="file1.text"
Content-Type: application/octet-stream

test1
--5eaf397248958ac38281d1c034e1ad0d4a5f7d986d4c53ac32e8399cbcda
Content-Disposition: form-data; name="files"; filename="file2.text"
Content-Type: application/octet-stream

test2
--5eaf397248958ac38281d1c034e1ad0d4a5f7d986d4c53ac32e8399cbcda--
`,
			&struct {
				FormDataMultipart `in:"body" mime:"multipart" boundary:"boundary1"`
			}{
				FormDataMultipart: FormDataMultipart{
					A:     []int{-1, 1},
					C:     1,
					Bytes: []byte("bytes"),
					Data: Data{
						A: "1",
					},
					Files: []*multipart.FileHeader{
						transformer.MustNewFileHeader("files", "file1.text", bytes.NewBufferString("test1")),
						transformer.MustNewFileHeader("files", "file2.text", bytes.NewBufferString("test2")),
					},
					File: transformer.MustNewFileHeader("file", "file.text", bytes.NewBufferString("test")),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for i := 0; i < 5; i++ {
				rt, err := factory.NewRequestTsfm(bgctx, reflect.TypeOf(c.req))
				NewWithT(t).Expect(err).To(BeNil())

				req, err := rt.NewRequest(http.MethodGet, c.path, c.req)
				NewWithT(t).Expect(err).To(BeNil())

				data, _ := httputil.DumpRequest(req, true)
				NewWithT(t).Expect(string(UnifyRequestData(data))).
					To(Equal(string(UnifyRequestData([]byte(c.expect)))))

				rv := reflectx.New(reflect.PtrTo(reflectx.DeRef(reflect.TypeOf(c.req))))
				err2 := rt.DecodeAndValidate(bgctx, httpx.NewRequestInfo(req), rv)
				NewWithT(t).Expect(err2).To(BeNil())
				NewWithT(t).Expect(reflectx.Indirect(rv).Interface()).
					To(Equal(reflectx.Indirect(reflect.ValueOf(c.req)).Interface()))
			}
		})
	}
}

func ExampleNewRequestTsfmFactory() {
	factory := NewRequestTsfmFactory(nil, nil)

	type PlainBody struct {
		A   string `json:"a"                         validate:"@string[2,]"`
		Int int    `json:"int,omitempty" default:"0" validate:"@int[0,]"`
	}

	type Req struct {
		Protocol  types.Protocol `in:"query" name:"protocol,omitempty" default:"HTTP"`
		QString   string         `in:"query" name:"string,omitempty"   default:"s"`
		PlainBody PlainBody      `in:"body"  mime:"plain" validate:"@struct<json>"`
	}

	req := &Req{}
	req.PlainBody.A = "1"

	rt, err := factory.NewRequestTsfm(bgctx, reflect.TypeOf(req))
	if err != nil {
		panic(err)
	}

	statusErr := rt.Params["body"][0].Validator.Validate(req.PlainBody)

	statusErr.(*vldterr.ErrorSet).Each(func(fieldErr *vldterr.FieldError) {
		fmt.Println(fieldErr.Field, strconv.Quote(fieldErr.Error.Error()))
	})
	// Output:
	// a "string length should be larger than 2, but got invalid value 1"
}

func TestRequestTsfm_DecodeFromRequestInfo_WithDefaults(t *testing.T) {
	type Data struct {
		String string `json:"string,omitempty" default:"111" validate:"@string[3,]"`
		Int    int    `json:"int,omitempty"    default:"111" validate:"@int[3,]"`
	}

	type Req struct {
		Protocol types.Protocol `in:"query" name:"protocol,omitempty" default:"HTTP"`
		QInt     int            `in:"query" name:"int,omitempty"      default:"1"`
		QString  string         `in:"query" name:"string,omitempty"   default:"s"`
		List     []Data         `in:"body"`
	}

	factory := NewRequestTsfmFactory(nil, nil)

	rt, err := factory.NewRequestTsfm(bgctx, reflect.TypeOf(&Req{}))
	NewWithT(t).Expect(err).To(BeNil())
	if err != nil {
		return
	}

	req, err := rt.NewRequest(http.MethodGet, "/", &Req{
		List: []Data{
			{
				String: "2222",
			},
			{},
		},
	})
	NewWithT(t).Expect(err).To(BeNil())

	r := &Req{}
	err = rt.DecodeAndValidate(bgctx, httpx.NewRequestInfo(req), r)
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(r).To(Equal(&Req{
		Protocol: types.PROTOCOL__HTTP,
		QInt:     1,
		QString:  "s",
		List: []Data{
			{
				String: "2222",
				Int:    111,
			},
			{
				String: "111",
				Int:    111,
			},
		},
	}))
}

func TestRequestTsfm_DecodeFromRequestInfo_WithEnumValidate(t *testing.T) {
	type Req struct {
		Protocol types.Protocol `name:"protocol,omitempty" validate:"@string{HTTP}" in:"query" default:"HTTP"`
	}

	factory := NewRequestTsfmFactory(nil, nil)

	rt, err := factory.NewRequestTsfm(bgctx, reflect.TypeOf(&Req{}))
	NewWithT(t).Expect(err).To(BeNil())

	req, err := rt.NewRequest(http.MethodGet, "/", &Req{types.PROTOCOL__HTTP})
	NewWithT(t).Expect(err).To(BeNil())

	r := &Req{}
	err = rt.DecodeAndValidate(bgctx, httpx.NewRequestInfo(req), r)
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(r).To(Equal(&Req{types.PROTOCOL__HTTP}))
}

func TestRequestTsfm_DecodeFromRequestInfo_Failed(t *testing.T) {
	factory := NewRequestTsfmFactory(nil, nil)

	type NestedForFailed struct {
		A string `json:"a" validate:"@string[1,]" errMsg:"A wrong"`
		B string `name:"b" validate:"@string[1,]" default:"1" `
		C string `json:"c" validate:"@string[2,]?"`
	}

	type DataForFailed struct {
		A               string `         validate:"@string[1,]"`
		B               string `         validate:"@string[1,]" default:"1" `
		C               string `json:"c" validate:"@string[2,]?"`
		NestedForFailed NestedForFailed
	}

	type ReqForFailed struct {
		ID            string   `in:"path"  name:"id"               validate:"@string[2,]"`
		QString       string   `in:"query" name:"string,omitempty" validate:"@string[2,]" default:"11" `
		QSlice        []string `in:"query" name:"slice,omitempty"  validate:"@slice<@string[2,]>[2,]"`
		DataForFailed `in:"body"`
	}

	rt, err := factory.NewRequestTsfm(bgctx, reflect.TypeOf(&ReqForFailed{}))
	if err != nil {
		return
	}

	req, err := rt.NewRequest(http.MethodGet, "/:id", &ReqForFailed{
		ID:            "1",
		QString:       "!",
		QSlice:        []string{"11", "1"},
		DataForFailed: DataForFailed{C: "1"},
	})
	if err != nil {
		return
	}

	e := rt.DecodeAndValidate(bgctx, httpx.NewRequestInfo(req), &ReqForFailed{})
	if e == nil {
		return
	}

	errFields := e.(*statusx.StatusErr).Fields

	sort.Slice(errFields, func(i, j int) bool {
		return errFields[i].Field < errFields[j].Field
	})

	data, _ := json.MarshalIndent(errFields, "", "  ")

	NewWithT(t).Expect(string(data)).To(Equal(`[
  {
    "field": "A",
    "msg": "missing required field",
    "in": "body"
  },
  {
    "field": "B",
    "msg": "missing required field",
    "in": "body"
  },
  {
    "field": "NestedForFailed.B",
    "msg": "missing required field",
    "in": "body"
  },
  {
    "field": "NestedForFailed.a",
    "msg": "A wrong",
    "in": "body"
  },
  {
    "field": "c",
    "msg": "string length should be larger than 2, but got invalid value 1",
    "in": "body"
  },
  {
    "field": "id",
    "msg": "string length should be larger than 2, but got invalid value 1",
    "in": "path"
  },
  {
    "field": "slice[1]",
    "msg": "string length should be larger than 2, but got invalid value 1",
    "in": "query"
  },
  {
    "field": "string",
    "msg": "string length should be larger than 2, but got invalid value 1",
    "in": "query"
  }
]`))
}

type ReqWithPostValidate struct {
	StartedAt string `in:"query"`
}

func (ReqWithPostValidate) PostValidate(badRequest BadRequestError) {
	badRequest.AddErr(pkgerr.Errorf("ops"), "query", "StartedAt")
}

func ExampleRequestTsfm_DecodeAndValidate_RequestInfo_FailedOfPost() {
	factory := NewRequestTsfmFactory(nil, nil)

	rt, err := factory.NewRequestTsfm(bgctx, reflect.TypeOf(&ReqWithPostValidate{}))
	if err != nil {
		return
	}

	req, err := rt.NewRequest(http.MethodPost, "/:id", &ReqWithPostValidate{})
	if err != nil {
		return
	}

	e := rt.DecodeAndValidate(bgctx, httpx.NewRequestInfo(req), &ReqWithPostValidate{})
	if e == nil {
		return
	}

	errFields := e.(*statusx.StatusErr).Fields

	sort.Slice(errFields, func(i, j int) bool {
		return errFields[i].Field < errFields[j].Field
	})

	for _, ef := range errFields {
		fmt.Println(ef)
	}
	// Output:
	// StartedAt in query - missing required field
	// StartedAt in query - ops
}
