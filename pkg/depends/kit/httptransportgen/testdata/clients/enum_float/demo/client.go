// This is a generated source file. DO NOT EDIT
// Source: demo/client.go

package demo

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

type Interface interface {
	Context() context.Context
	WithContext(context.Context) Interface
}

func NewClient(c kit.Client) *Client {
	return &(Client{
		Client: c,
	})
}

type Client struct {
	Client kit.Client
	ctx    context.Context
}

func (c *Client) Context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

func (c *Client) WithContext(ctx context.Context) Interface {
	cc := new(Client)
	cc.Client, cc.ctx = c.Client, ctx
	return cc
}
