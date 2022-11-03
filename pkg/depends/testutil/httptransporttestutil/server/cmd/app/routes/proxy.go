package routes

import (
	"context"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/client"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var ProxyRouter = kit.NewRouter(httptransport.Group("/proxy"))

var (
	c = &client.Client{
		Host:    "ip-api.com",
		Timeout: 100 * time.Second,
	}
)

func init() {
	c.SetDefault()

	RootRouter.Register(ProxyRouter)

	ProxyRouter.Register(kit.NewRouter(&Proxy{}))
	ProxyRouter.Register(kit.NewRouter(&ProxyV2{}))
}

type Proxy struct {
	httpx.MethodGet
}

func (Proxy) Output(ctx context.Context) (interface{}, error) {
	resp := &IpInfo{}
	_, err := c.Do(ctx, &GetByJSON{}).Into(resp)
	return resp, err
}

type ProxyV2 struct {
	httpx.MethodGet `basePath:"/demo/v2"`
}

func (ProxyV2) Output(ctx context.Context) (interface{}, error) {
	result := c.Do(ctx, &GetByJSON{})

	return httpx.WrapSchema(&IpInfo{})(result), nil
}

type GetByJSON struct {
	httpx.MethodGet
}

func (GetByJSON) Path() string {
	return "/json"
}

type IpInfo struct {
	Country     string `json:"country"     xml:"country"`
	CountryCode string `json:"countryCode" xml:"countryCode"`
}
