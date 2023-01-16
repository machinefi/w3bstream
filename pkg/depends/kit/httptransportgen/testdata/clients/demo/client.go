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
	Cookie(req *Cookie, metas ...kit.Metadata) (kit.Metadata, error)
	Create(req *Create, metas ...kit.Metadata) (*Data, kit.Metadata, error)
	DownloadFile(metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxAttachment, kit.Metadata, error)
	FormMultipartWithFile(req *FormMultipartWithFile, metas ...kit.Metadata) (kit.Metadata, error)
	FormMultipartWithFiles(req *FormMultipartWithFiles, metas ...kit.Metadata) (kit.Metadata, error)
	FormURLEncoded(req *FormURLEncoded, metas ...kit.Metadata) (kit.Metadata, error)
	GetByID(req *GetByID, metas ...kit.Metadata) (*Data, kit.Metadata, error)
	HealthCheck(req *HealthCheck, metas ...kit.Metadata) (kit.Metadata, error)
	Proxy(metas ...kit.Metadata) (*IpInfo, kit.Metadata, error)
	ProxyV2(metas ...kit.Metadata) (*IpInfo, kit.Metadata, error)
	Redirect(metas ...kit.Metadata) (kit.Metadata, error)
	RedirectWhenError(metas ...kit.Metadata) (kit.Metadata, error)
	RemoveByID(req *RemoveByID, metas ...kit.Metadata) (kit.Metadata, error)
	ShowImage(metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxImagePNG, kit.Metadata, error)
	UpdateByID(req *UpdateByID, metas ...kit.Metadata) (kit.Metadata, error)
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

func (c *Client) Cookie(req *Cookie, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) Create(req *Create, metas ...kit.Metadata) (*Data, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) DownloadFile(metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxAttachment, kit.Metadata, error) {
	return (&DownloadFile{}).InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) FormMultipartWithFile(req *FormMultipartWithFile, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) FormMultipartWithFiles(req *FormMultipartWithFiles, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) FormURLEncoded(req *FormURLEncoded, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) GetByID(req *GetByID, metas ...kit.Metadata) (*Data, kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) HealthCheck(req *HealthCheck, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) Proxy(metas ...kit.Metadata) (*IpInfo, kit.Metadata, error) {
	return (&Proxy{}).InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ProxyV2(metas ...kit.Metadata) (*IpInfo, kit.Metadata, error) {
	return (&ProxyV2{}).InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) Redirect(metas ...kit.Metadata) (kit.Metadata, error) {
	return (&Redirect{}).InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RedirectWhenError(metas ...kit.Metadata) (kit.Metadata, error) {
	return (&RedirectWhenError{}).InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) RemoveByID(req *RemoveByID, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) ShowImage(metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxImagePNG, kit.Metadata, error) {
	return (&ShowImage{}).InvokeContext(c.Context(), c.Client, metas...)
}

func (c *Client) UpdateByID(req *UpdateByID, metas ...kit.Metadata) (kit.Metadata, error) {
	return req.InvokeContext(c.Context(), c.Client, metas...)
}
