package http

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/client"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
)

type ClientEndpoint struct {
	Endpoint types.Endpoint `env:""`
	Timeout  time.Duration

	client.Client `env:"-"`
}

func (c *ClientEndpoint) Do(ctx context.Context, req interface{}, metas ...kit.Metadata) kit.Result {
	return c.Client.Do(ctx, req, metas...)
}

func (c *ClientEndpoint) LivenessCheck() map[string]string {
	s := map[string]string{}
	s[c.Endpoint.Host()] = "ok"

	_, err := c.Do(context.Background(), NewRequest(http.MethodGet, "/liveness")).Into(&s)
	if err != nil {
		if statusx.FromErr(err).StatusCode() != http.StatusNotFound {
			s[c.Endpoint.Host()] = err.Error()
		}
	}
	return s
}

func NewRequest(method string, path string) *request {
	return &request{method: method, path: path}
}

type request struct {
	method string
	path   string
}

func (req *request) Method() string { return req.method }

func (req *request) Path() string { return req.path }

func (c *ClientEndpoint) SetDefault() {
	c.Client.SetDefault()
	c.Client.Transports = []client.HttpTransport{
		NewLogRoundTripper(logrus.WithContext(context.Background())),
	}
}

func (c *ClientEndpoint) Init() {
	if c.Endpoint.Scheme != "" {
		c.Client.Protocol = c.Endpoint.Scheme
	}
	if c.Endpoint.Hostname != "" {
		c.Client.Host = c.Endpoint.Hostname
	}
	if c.Endpoint.Port != 0 {
		c.Client.Port = c.Endpoint.Port
	}
	if c.Timeout != 0 {
		c.Client.Timeout = c.Timeout
	}
}
