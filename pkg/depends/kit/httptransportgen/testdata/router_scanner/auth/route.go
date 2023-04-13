package auth

import (
	"context"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
)

type NoContent struct {
}

func (NoContent) Output(ctx context.Context) (interface{}, error) {
	return nil, nil
}

// Auth
// auth auth
type Auth struct {
	Headers
	Queries
	Cookies
	Data `in:"body"`
}

func (auth Auth) Output(ctx context.Context) (result interface{}, err error) {
	return auth.Data, nil
}

type Headers struct {
	HInt    int    `in:"header"`
	HString string `in:"header"`
	HBool   bool   `in:"header"`
}

type Queries struct {
	QInt            int      `name:"int" in:"query"`
	QString         string   `name:"string" in:"query"`
	QSlice          []string `name:"slice" in:"query"`
	QBytes          []byte   `name:"bytes,omitempty" in:"query"`
	QBytesOmitEmpty []byte   `name:"bytesOmit,omitempty" in:"query"`
	QBytesKeepEmpty []byte   `name:"bytesKeep" in:"query"`
}

type Cookies struct {
	CString string   `name:"a" in:"cookie"`
	CSlice  []string `name:"slice" in:"cookie"`
}

type Data struct {
	A     string
	B     string
	C     string
	ctype string
}

func (Data) Status() int {
	return http.StatusOK
}

func (d *Data) ContentType() string {
	return d.ctype
}

type FormDataMultipart struct {
	Bytes []byte `name:"bytes"`
	A     []int  `name:"a"`
	C     uint   `name:"c" `
	Data  Data   `name:"data"`

	File  *multipart.FileHeader   `name:"file"`
	Files []*multipart.FileHeader `name:"files"`
}

type RespWithDescribers struct {
}

func (RespWithDescribers) Output(ctx context.Context) (r interface{}, err error) {
	if true {
		r = httpx.Compose(
			httpx.WrapCookies(&http.Cookie{
				Name:    "token",
				Value:   "111",
				Expires: time.Now().Add(24 * time.Hour),
			}),
			httpx.WrapStatusCode(http.StatusOK),
			httpx.WrapMeta(httpx.Metadata("X-A", "XXX")),
		)(nil)
		return
	}

	if true {
		return httpx.Compose(
			httpx.WrapCookies(&http.Cookie{
				Name:    "token",
				Value:   "111",
				Expires: time.Now().Add(24 * time.Hour),
			}),
			httpx.WrapStatusCode(http.StatusOK),
			httpx.WrapContentType(httpx.MIME_JSON),
		)(nil), nil
	}

	resp := httpx.Compose(
		httpx.WrapCookies(&http.Cookie{
			Name:    "token",
			Value:   "111",
			Expires: time.Now().Add(24 * time.Hour),
		}),
		httpx.WrapStatusCode(http.StatusOK),
		httpx.WrapContentType(httpx.MIME_JSON),
	)(nil)
	return resp, nil
}
