package client_test

import (
	"context"
	"encoding/xml"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/client"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/mock"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
)

type IpInfo struct {
	xml.Name    `xml:"query"`
	Country     string `json:"country" xml:"country"`
	CountryCode string `json:"countryCode" xml:"countryCode"`
}

type GetByJSON struct {
	httpx.MethodGet
}

func (GetByJSON) Path() string {
	return "/me.json"
}

type GetByXML struct {
	httpx.MethodGet
}

func (GetByXML) Path() string {
	return "/me.xml"
}

func TestClient(t *testing.T) {
	if runtime.GOOS == `darwin` {
		return
	}

	cli := &client.Client{
		Protocol: "https",
		Host:     "ip.nf",
		Timeout:  100 * time.Second,
	}
	cli.SetDefault()

	t.Run("direct request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "https://api.github.com", nil)
		_, err := cli.Do(context.Background(), req).Into(nil)
		require.NoError(t, err)
	})

	t.Run("direct request 404", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "https://api.github.com/xxxxn", nil)

		meta, err := cli.Do(context.Background(), req).Into(nil)
		require.Error(t, err)

		t.Log(err)
		t.Log(meta)
	})

	t.Run("request by struct", func(t *testing.T) {
		rsp := IpInfo{}

		meta, err := cli.Do(context.Background(), &GetByJSON{}).Into(&rsp)
		require.NoError(t, err)

		t.Log(rsp)
		t.Log(meta)
	})

	t.Run("request by struct as xml", func(t *testing.T) {
		rsp := IpInfo{}

		meta, err := cli.Do(context.Background(), &GetByXML{}).Into(&rsp)
		require.NoError(t, err)

		t.Log(rsp)
		t.Log(meta)
	})

	t.Run("cancel request", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(1 * time.Millisecond)
			cancel()
		}()

		rsp := IpInfo{}
		_, err := cli.Do(ctx, &GetByJSON{}).Into(&rsp)
		require.Equal(t, "ClientClosedRequest", err.(*statusx.StatusErr).Key)
	})

	t.Run("err request", func(t *testing.T) {
		errcli := &client.Client{
			Timeout: 100 * time.Second,
		}
		errcli.SetDefault()

		{
			rsp := IpInfo{}

			_, err := errcli.Do(
				client.ContextWithClient(
					context.Background(),
					client.GetShortConnClientContext(context.Background(), 10*time.Second),
				),
				&GetByJSON{},
			).Into(&rsp)
			require.Error(t, err)
		}
	})

	t.Run("result pass", func(t *testing.T) {
		req1, _ := http.NewRequest("GET", "https://ip.nf/me.json", nil)
		res := cli.Do(context.Background(), req1)

		req2, _ := http.NewRequest(http.MethodGet, "/", nil)
		rw := mock.NewMockResponseWriter()
		_ = httpx.ResponseFrom(res).WriteTo(rw, req2, nil)

		require.Equal(t, 200, rw.StatusCode)
	})
}
