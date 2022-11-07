package client

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/textproto"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/http2"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/client/roundtrippers"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

type HttpTransport func(rt http.RoundTripper) http.RoundTripper

type Client struct {
	Protocol           string
	Host               string
	Port               uint16
	Timeout            time.Duration
	RequestTsfmFactory *httptransport.RequestTsfmFactory
	Transports         []HttpTransport
	NewError           func(resp *http.Response) error
}

func (c *Client) SetDefault() {
	if c.RequestTsfmFactory == nil {
		c.RequestTsfmFactory = httptransport.NewRequestTsfmFactory(nil, nil)
		c.RequestTsfmFactory.SetDefault()
	}
	if c.Transports == nil {
		c.Transports = []HttpTransport{roundtrippers.NewLogRoundTripper()}
	}
	if c.NewError == nil {
		c.NewError = func(resp *http.Response) error {
			return &statusx.StatusErr{
				Code:    resp.StatusCode * 1e6,
				Msg:     resp.Status,
				Sources: []string{resp.Request.Host},
			}
		}
	}
}

type keyClient struct{}

func ContextWithClient(ctx context.Context, c *http.Client) context.Context {
	return contextx.WithValue(ctx, keyClient{}, c)
}

func ClientFromContext(ctx context.Context) *http.Client {
	if ctx == nil {
		return nil
	}
	if c, ok := ctx.Value(keyClient{}).(*http.Client); ok {
		return c
	}
	return nil
}

type keyDftTransport struct{}

func ContextWithDftTransport(ctx context.Context, t *http.Transport) context.Context {
	return contextx.WithValue(ctx, keyDftTransport{}, t)
}

func DftTransportFromContext(ctx context.Context) *http.Transport {
	if ctx == nil {
		return nil
	}
	if t, ok := ctx.Value(keyDftTransport{}).(*http.Transport); ok {
		return t
	}
	return nil
}

func (c *Client) Do(ctx context.Context, req interface{}, metas ...kit.Metadata) kit.Result {
	request, ok := req.(*http.Request)
	if !ok {
		request2, err := c.newRequest(ctx, req, metas...)
		if err != nil {
			return &Result{
				Err:      statusx.Wrap(err, http.StatusInternalServerError, "RequestFailed"),
				NewError: c.NewError,
				Tsfm:     c.RequestTsfmFactory.Tsfm,
			}
		}
		request = request2
	}

	httpClient := ClientFromContext(ctx)
	if httpClient == nil {
		httpClient = GetShortConnClientContext(ctx, c.Timeout, c.Transports...)
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		if errors.Unwrap(err) == context.Canceled {
			return &Result{
				Err:      statusx.Wrap(err, 499, "ClientClosedRequest"),
				NewError: c.NewError,
				Tsfm:     c.RequestTsfmFactory.Tsfm,
			}
		}

		return &Result{
			Err:      statusx.Wrap(err, http.StatusInternalServerError, "RequestFailed"),
			NewError: c.NewError,
			Tsfm:     c.RequestTsfmFactory.Tsfm,
		}
	}
	return &Result{
		NewError: c.NewError,
		Tsfm:     c.RequestTsfmFactory.Tsfm,
		Response: resp,
	}
}

func (c *Client) toUrl(path string) string {
	protocol := c.Protocol
	if protocol == "" {
		protocol = "http"
	}
	url := fmt.Sprintf("%s://%s", protocol, c.Host)
	if c.Port > 0 {
		url = fmt.Sprintf("%s:%d", url, c.Port)
	}
	return url + path
}

func (c *Client) newRequest(ctx context.Context, req interface{}, metas ...kit.Metadata) (*http.Request, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	method := ""
	path := ""

	if methodDescriber, ok := req.(httptransport.MethodDescriber); ok {
		method = methodDescriber.Method()
	}

	if pathDescriber, ok := req.(httptransport.PathDescriber); ok {
		path = pathDescriber.Path()
	}

	request, err := c.RequestTsfmFactory.NewRequestWithContext(ctx, method, c.toUrl(path), req)
	if err != nil {
		return nil, statusx.Wrap(err, http.StatusBadRequest, "RequestTransformFailed")
	}

	request = request.WithContext(ctx)

	for k, vs := range kit.FromMetas(metas...) {
		for _, v := range vs {
			request.Header.Add(k, v)
		}
	}

	return request, nil
}

type Result struct {
	Tsfm     transformer.Factory
	Response *http.Response
	NewError func(resp *http.Response) error
	Err      error
}

func (r *Result) StatusCode() int {
	if r.Response != nil {
		return r.Response.StatusCode
	}
	return 0
}

func (r *Result) Meta() kit.Metadata {
	if r.Response != nil {
		return kit.Metadata(r.Response.Header)
	}
	return kit.Metadata{}
}

func (r *Result) Into(body interface{}) (kit.Metadata, error) {
	defer func() {
		if r.Response != nil && r.Response.Body != nil {
			r.Response.Body.Close()
		}
	}()

	if r.Err != nil {
		return nil, r.Err
	}

	meta := kit.Metadata(r.Response.Header)

	if !isOk(r.Response.StatusCode) {
		body = r.NewError(r.Response)
	}

	if body == nil {
		return meta, nil
	}

	decode := func(body interface{}) error {
		contentType := meta.Get(httpx.HeaderContentType)

		if contentType != "" {
			contentType, _, _ = mime.ParseMediaType(contentType)
		}

		rv := reflect.ValueOf(body)

		tsfm, err := r.Tsfm.NewTransformer(
			context.Background(),
			typesx.FromReflectType(rv.Type()),
			transformer.Option{MIME: contentType},
		)

		if err != nil {
			return statusx.Wrap(err, http.StatusInternalServerError, "ReadFailed")
		}
		if e := tsfm.DecodeFrom(
			context.Background(),
			r.Response.Body,
			rv,
			textproto.MIMEHeader(r.Response.Header),
		); e != nil {
			return statusx.Wrap(e, http.StatusInternalServerError, "DecodeFailed")
		}
		return nil
	}

	switch v := body.(type) {
	case error:
		// to unmarshal status error
		if err := decode(v); err != nil {
			return meta, err
		}
		return meta, v
	case io.Writer:
		if _, err := io.Copy(v, r.Response.Body); err != nil {
			return meta, statusx.Wrap(err, http.StatusInternalServerError, "WriteFailed")
		}
	default:
		if err := decode(body); err != nil {
			return meta, err
		}
	}

	return meta, nil
}

func isOk(code int) bool {
	return code >= http.StatusOK && code < http.StatusMultipleChoices
}

func GetShortConnClientContext(
	ctx context.Context,
	timeout time.Duration,
	transports ...HttpTransport,
) *http.Client {
	t := DftTransportFromContext(ctx)

	if t != nil {
		t = t.Clone()
	} else {
		t = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 0,
			}).DialContext,
			DisableKeepAlives:     true,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	}

	if err := http2.ConfigureTransport(t); err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: t,
	}

	for i := range transports {
		httpTransport := transports[i]
		client.Transport = httpTransport(client.Transport)
	}

	return client
}
