package httpx

import (
	"context"
	"io"
	"mime"
	"net/http"
	"net/textproto"
	"net/url"

	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
)

type Encode func(context.Context, io.Writer, interface{}) error

type ResolveEncode func(*Response) (Encode, error)

type ResponseWrapper func(v interface{}) *Response

func Compose(rws ...ResponseWrapper) ResponseWrapper {
	return func(v interface{}) *Response {
		rsp := ResponseFrom(v)
		for i := len(rws) - 1; i >= 0; i-- {
			rsp = rws[i](rsp)
		}
		return rsp
	}
}

func WrapStatusCode(statusCode int) ResponseWrapper {
	return func(v interface{}) *Response {
		rsp := ResponseFrom(v)
		rsp.StatusCode = statusCode
		return rsp
	}
}

func WrapCookies(cookies ...*http.Cookie) ResponseWrapper {
	return func(v interface{}) *Response {
		rsp := ResponseFrom(v)
		rsp.Cookies = cookies
		return rsp
	}
}

func WrapSchema(s interface{}) ResponseWrapper {
	return func(v interface{}) *Response {
		rsp := ResponseFrom(v)
		return rsp
	}
}

func WrapContentType(ct string) ResponseWrapper {
	return func(v interface{}) *Response {
		rsp := ResponseFrom(v)
		rsp.ContentType = ct
		return rsp
	}
}

func WrapMeta(metas ...kit.Metadata) ResponseWrapper {
	return func(v interface{}) *Response {
		rsp := ResponseFrom(v)
		rsp.Meta = kit.FromMetas(metas...)
		return rsp
	}
}

func Metadata(k string, vs ...string) kit.Metadata {
	return kit.Metadata{k: vs}
}

func ResponseFrom(v interface{}) *Response {
	if r, ok := v.(*Response); ok {
		return r
	}

	rsp := &Response{}

	if redirectDescriber, ok := v.(RedirectDescriber); ok {
		rsp.Location = redirectDescriber.Location()
		rsp.StatusCode = redirectDescriber.StatusCode()
		return rsp
	}

	if e, ok := v.(error); ok {
		if e != nil {
			se, ok := statusx.IsStatusErr(e)
			if !ok {
				if e == context.Canceled {
					// https://httpstatuses.com/499
					se = statusx.Wrap(e, 499, "ContextCanceled")
				} else {
					se = statusx.Wrap(e, http.StatusInternalServerError, "UnknownError")
				}
			}
			v = se
		}
	}

	rsp.Value = v

	if with, ok := v.(kit.MetadataCarrier); ok {
		rsp.Meta = with.Meta()
	}

	if with, ok := v.(WithCookies); ok {
		rsp.Cookies = with.Cookies()
	}

	if with, ok := v.(WithContentType); ok {
		rsp.ContentType = with.ContentType()
	}

	if with, ok := v.(WithStatusCode); ok {
		rsp.StatusCode = with.StatusCode()
	}

	return rsp
}

type Response struct {
	Value       interface{} // Value of body
	Meta        kit.Metadata
	Cookies     []*http.Cookie
	Location    *url.URL
	ContentType string
	StatusCode  int
}

func (r *Response) Unwrap() error {
	if err, ok := r.Value.(error); ok {
		return err
	}
	return nil
}

func (r *Response) Error() string {
	if err := r.Unwrap(); err != nil {
		return err.Error()
	}
	return "response error"
}

func (r *Response) WriteTo(rw http.ResponseWriter, req *http.Request, resolve ResolveEncode) error {
	defer func() { r.Value = nil }()
	if upgrader, ok := r.Value.(Upgrader); ok {
		return upgrader.Upgrade(rw, req)
	}
	if r.StatusCode == 0 {
		if r.Value == nil {
			r.StatusCode = http.StatusNoContent
		} else {
			if req.Method == http.MethodPost {
				r.StatusCode = http.StatusCreated
			} else {
				r.StatusCode = http.StatusOK
			}
		}
	}
	if r.Meta != nil {
		header := rw.Header()
		for k, vs := range r.Meta {
			header[textproto.CanonicalMIMEHeaderKey(k)] = vs
		}
	}
	if r.Cookies != nil {
		for _, cookie := range r.Cookies {
			if cookie != nil {
				http.SetCookie(rw, cookie)
			}
		}
	}
	if r.Location != nil {
		http.Redirect(rw, req, r.Location.String(), r.StatusCode)
		return nil
	}
	if r.StatusCode == http.StatusNoContent {
		rw.WriteHeader(r.StatusCode)
		return nil
	}
	if r.ContentType != "" {
		rw.Header().Set(HeaderContentType, r.ContentType)
	}
	switch v := r.Value.(type) {
	case kit.Result:
		rw.WriteHeader(r.StatusCode)
		_, err := v.Into(rw)
		return err
	case io.Reader:
		rw.WriteHeader(r.StatusCode)
		defer func() {
			if c, ok := v.(io.Closer); ok {
				c.Close()
			}
		}()
		_, err := io.Copy(rw, v)
		return err
	default:
		enc, err := resolve(r)
		if err != nil {
			return err
		}
		return enc(
			ContextWithStatusCode(req.Context(), r.StatusCode),
			rw, r.Value,
		)
	}
}

type ResponseWriteError interface {
	WriteError(err error) (int, error)
}

type Upgrader interface {
	Upgrade(http.ResponseWriter, *http.Request) error
}

func MaybeWriteHeader(ctx context.Context, w io.Writer, ct string, param map[string]string) {
	if rw, ok := w.(WithHeader); ok {
		if len(param) == 0 {
			rw.Header().Set(HeaderContentType, ct)
		} else {
			rw.Header().Set(
				HeaderContentType,
				mime.FormatMediaType(ct, param),
			)
		}
	}
	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(StatusCodeFromContext(ctx))
	}
}
