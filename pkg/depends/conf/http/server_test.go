package http_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	confhttp "github.com/machinefi/w3bstream/pkg/depends/conf/http"
	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

type GetSome struct {
	httpx.MethodGet
}

func (GetSome) Path() string {
	return "/some"
}

func (GetSome) Output(ctx context.Context) (interface{}, error) {
	html := httpx.NewHTML()

	return httpx.WrapMeta(
		httpx.Metadata("Cache-Control", "no-cache"),
	)(html), nil
}

type GetOther struct {
	httpx.MethodGet
}

func (GetOther) Path() string { return "/other" }

func (GetOther) Output(ctx context.Context) (interface{}, error) {
	client := confhttp.ClientEndpoint{
		Endpoint: types.Endpoint{
			Scheme:   "http",
			Hostname: "0.0.0.0",
			Port:     1234,
		},
	}

	client.SetDefault()
	client.Init()

	_, _ = client.Do(ctx, confhttp.NewRequest(http.MethodGet, "/some")).Into(nil)

	return nil, nil
}

func TestHttp(t *testing.T) {
	ctx, l := logger.NewSpanContext(context.Background(), "TestHttp")
	defer l.End()

	servers := []*confhttp.Server{
		{Protocol: "http", Port: 1234},
		{Protocol: "http", Port: 3456},
		{Protocol: "http+unix", Addr: "/tmp/server3.sock"},
	}
	clients := make([]*confhttp.ClientEndpoint, 0, len(servers))

	router := kit.NewRouter(httptransport.Group("/"))
	router.Register(kit.NewRouter(&GetSome{}))
	router.Register(kit.NewRouter(&GetOther{}))

	for i, srv := range servers {
		srv.SetDefault()
		go func(i int, srv *confhttp.Server) {
			err := srv.Serve(router)
			fmt.Printf("server#%d: %v", i, err)
			time.Sleep(5 * time.Second)
			p, _ := os.FindProcess(os.Getpid())
			_ = p.Signal(os.Interrupt)
		}(i, srv)

		cli := &confhttp.ClientEndpoint{
			Endpoint: types.Endpoint{
				Scheme:   srv.Protocol,
				Hostname: srv.Addr,
				Port:     uint16(srv.Port),
			},
		}
		cli.SetDefault()
		cli.Init()

		clients = append(clients, cli)
	}

	time.Sleep(time.Second)

	printer := func(rsp *http.Response) {
		data, _ := httputil.DumpResponse(rsp, true)
		fmt.Println(string(data))
	}

	for i := 0; i < len(servers); i++ {
		idx := strconv.Itoa(i + 1)
		cli := clients[i]
		srv := servers[i]

		baseURL := srv.Protocol + "://" + srv.Addr
		if srv.Port != 0 {
			baseURL += ":" + strconv.Itoa(srv.Port)
		}
		if srv.Protocol == "http+unix" {
			baseURL = "http://localhost"
		}

		// use confhttp.ClientEndpoint
		t.Run("GetSome#"+idx, func(t *testing.T) {
			meta, err := cli.Do(ctx, confhttp.NewRequest(http.MethodGet, "/some")).Into(nil)
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(http.Header(meta).Get("b3")).NotTo(BeEmpty())
		})
		t.Run("GetOther#"+idx, func(t *testing.T) {
			meta, err := cli.Do(ctx, confhttp.NewRequest(http.MethodGet, "/other")).Into(nil)
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(http.Header(meta).Get("b3")).NotTo(BeEmpty())
		})

		// use http.DefaultClient
		if srv.Protocol == "http+unix" {
			http.DefaultClient.Transport = &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
					return net.Dial("unix", srv.Addr)
				},
			}
		} else {
			http.DefaultClient.Transport = http.DefaultTransport
		}

		t.Run("Head#"+idx, func(t *testing.T) {
			resp, err := http.Head(baseURL)
			NewWithT(t).Expect(err).To(BeNil())
			printer(resp)
		})
		t.Run("Options#"+idx, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodOptions, fmt.Sprintf(baseURL+"/some"), nil)
			req.Header.Add("Origin", "http://localhost:3000")
			req.Header.Add("Access-Control-Request-Method", http.MethodPost)
			req.Header.Set("Access-Control-Request-Headers", "authorization,content-type")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
			resp, err := http.DefaultClient.Do(req)
			NewWithT(t).Expect(err).To(BeNil())
			printer(resp)
		})
	}
	time.Sleep(1 * time.Second)
}
