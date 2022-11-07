package httpx

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func NewRequestInfo(r *http.Request) RequestInfo {
	return &requestInfo{
		recvts:  time.Now(),
		request: r,
	}
}

type WithFromRequestInfo interface {
	FromRequestInfo(ri RequestInfo) error
}

type RequestInfo interface {
	Context() context.Context
	Values(in string, name string) []string
	Header() http.Header
	Body() io.ReadCloser
}

type requestInfo struct {
	request *http.Request
	recvts  time.Time
	query   url.Values
	cookies []*http.Cookie
	params  httprouter.Params
}

func (ri *requestInfo) Header() http.Header { return ri.request.Header }

func (ri *requestInfo) Context() context.Context { return ri.request.Context() }

func (ri *requestInfo) Body() io.ReadCloser { return ri.request.Body }

func (ri *requestInfo) Value(in string, name string) string {
	values := ri.Values(in, name)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (ri *requestInfo) Values(in string, name string) []string {
	switch in {
	case "path":
		v := ri.Param(name)
		if v == "" {
			return []string{}
		}
		return []string{v}
	case "query":
		return ri.QueryValues(name)
	case "cookie":
		return ri.CookieValues(name)
	case "header":
		return ri.HeaderValues(name)
	}
	return []string{}
}

func (ri *requestInfo) Param(name string) string {
	if ri.params == nil {
		params, ok := ri.request.Context().Value(httprouter.ParamsKey).(httprouter.Params)
		if !ok {
			params = httprouter.Params{}
		}
		ri.params = params
	}
	return ri.params.ByName(name)
}

func (ri *requestInfo) QueryValues(name string) []string {
	if ri.query == nil {
		ri.query = ri.request.URL.Query()

		if ri.request.Method == http.MethodGet && len(ri.query) == 0 &&
			ri.request.ContentLength > 0 {
			if strings.HasPrefix(ri.request.Header.Get("Content-Type"), MIME_FORM_URLENCODED) {
				data, err := ioutil.ReadAll(ri.request.Body)
				if err == nil {
					ri.request.Body.Close()

					query, e := url.ParseQuery(string(data))
					if e == nil {
						ri.query = query
					}
				}
			}
		}
	}
	return ri.query[name]
}

func (ri *requestInfo) HeaderValues(name string) []string {
	return ri.request.Header[textproto.CanonicalMIMEHeaderKey(name)]
}

func (ri *requestInfo) CookieValues(name string) []string {
	if ri.cookies == nil {
		ri.cookies = ri.request.Cookies()
	}

	values := make([]string, 0)
	for _, c := range ri.cookies {
		if c.Name == name {
			if c.Expires.IsZero() {
				values = append(values, c.Value)
			} else if c.Expires.After(ri.recvts) {
				values = append(values, c.Value)
			}
		}
	}
	return values
}
