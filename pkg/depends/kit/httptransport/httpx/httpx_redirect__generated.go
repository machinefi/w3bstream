// This is a generated source file. DO NOT EDIT
// Source: httpx/httpx_redirect__generated.go

package httpx

import (
	"net/http"
	"net/url"
)

func RedirectWithStatusMultipleChoices(u *url.URL) *StatusMultipleChoices {
	return &StatusMultipleChoices{Response: &Response{Location: u}}
}

type StatusMultipleChoices struct{ *Response }

func (StatusMultipleChoices) StatusCode() int { return http.StatusMultipleChoices }

func (r StatusMultipleChoices) Location() *url.URL { return r.Response.Location }

func RedirectWithStatusMovedPermanently(u *url.URL) *StatusMovedPermanently {
	return &StatusMovedPermanently{Response: &Response{Location: u}}
}

type StatusMovedPermanently struct{ *Response }

func (StatusMovedPermanently) StatusCode() int { return http.StatusMovedPermanently }

func (r StatusMovedPermanently) Location() *url.URL { return r.Response.Location }

func RedirectWithStatusFound(u *url.URL) *StatusFound {
	return &StatusFound{Response: &Response{Location: u}}
}

type StatusFound struct{ *Response }

func (StatusFound) StatusCode() int { return http.StatusFound }

func (r StatusFound) Location() *url.URL { return r.Response.Location }

func RedirectWithStatusSeeOther(u *url.URL) *StatusSeeOther {
	return &StatusSeeOther{Response: &Response{Location: u}}
}

type StatusSeeOther struct{ *Response }

func (StatusSeeOther) StatusCode() int { return http.StatusSeeOther }

func (r StatusSeeOther) Location() *url.URL { return r.Response.Location }

func RedirectWithStatusNotModified(u *url.URL) *StatusNotModified {
	return &StatusNotModified{Response: &Response{Location: u}}
}

type StatusNotModified struct{ *Response }

func (StatusNotModified) StatusCode() int { return http.StatusNotModified }

func (r StatusNotModified) Location() *url.URL { return r.Response.Location }

func RedirectWithStatusUseProxy(u *url.URL) *StatusUseProxy {
	return &StatusUseProxy{Response: &Response{Location: u}}
}

type StatusUseProxy struct{ *Response }

func (StatusUseProxy) StatusCode() int { return http.StatusUseProxy }

func (r StatusUseProxy) Location() *url.URL { return r.Response.Location }

func RedirectWithStatusTemporaryRedirect(u *url.URL) *StatusTemporaryRedirect {
	return &StatusTemporaryRedirect{Response: &Response{Location: u}}
}

type StatusTemporaryRedirect struct{ *Response }

func (StatusTemporaryRedirect) StatusCode() int { return http.StatusTemporaryRedirect }

func (r StatusTemporaryRedirect) Location() *url.URL { return r.Response.Location }

func RedirectWithStatusPermanentRedirect(u *url.URL) *StatusPermanentRedirect {
	return &StatusPermanentRedirect{Response: &Response{Location: u}}
}

type StatusPermanentRedirect struct{ *Response }

func (StatusPermanentRedirect) StatusCode() int { return http.StatusPermanentRedirect }

func (r StatusPermanentRedirect) Location() *url.URL { return r.Response.Location }
