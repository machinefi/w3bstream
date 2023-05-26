package integrations

import (
	"bytes"
	"context"
	_ "embed"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/clients/applet_mgr"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/requires"
	client2 "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/client"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/types"
)

//go:embed testdata/log.wasm
var code []byte

var (
	clientEvent           *applet_mgr.Client
	projectNameEventBench = "test_event_benchmark"
	publisherToken        string
)

func required() func() {
	var (
		client    = requires.AuthClient()
		projectID types.SFID
	)

	{
		req := &applet_mgr.CreateProject{}
		req.CreateReq.Name = projectNameEventBench

		rsp, _, err := client.CreateProject(req)
		if err != nil {
			panic(err)
		}
		projectID = rsp.ProjectID
	}

	{
		cwd, _ := os.Getwd()
		filename := path.Join(cwd, "../testdata/log.wasm")
		req := &applet_mgr.CreateApplet{
			ProjectName: projectNameEventBench,
		}
		req.CreateReq.File = transformer.MustNewFileHeader("file", filename, bytes.NewBuffer(code))
		req.CreateReq.Info = applet_mgr.GithubComMachinefiW3BstreamPkgModulesAppletInfo{
			AppletName: "log",
			WasmName:   "log.wasm",
		}

		_, _, err := client.CreateApplet(req)
		if err != nil {
			panic(err)
		}
	}

	{
		req := &applet_mgr.CreatePublisher{
			ProjectName: projectNameEventBench,
		}
		req.CreateReq.Name = "test_publisher"
		req.CreateReq.Key = "mn_test_publisher"

		rsp, _, err := client.CreatePublisher(req)
		if err != nil {
			panic(err)
		}
		publisherToken = rsp.Token
	}

	clientEvent = requires.ClientEvent()
	clientEvent.WithContext(client2.ContextWithDftTransport(context.Background(),
		&http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 100 * time.Second,
			}).DialContext,
			DisableKeepAlives:     false,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	))

	return func() {
		if projectID != 0 {
			_, _ = client.RemoveProject(&applet_mgr.RemoveProject{
				ProjectName: projectNameEventBench,
			})
		}
	}
}

var onceRequire = &sync.Once{}

func BenchmarkEventHandling(b *testing.B) {
	var stop func()
	onceRequire.Do(func() {
		stop = required()
	})
	defer stop()

	channel := strings.Join([]string{"aid", requires.AccountID.String(), projectNameEventBench}, "_")
	failed := 0
	b.N = 10
	for i := 0; i < b.N; i++ {
		req := &applet_mgr.HandleEvent{
			Channel:      channel,
			AuthInHeader: "Bearer " + publisherToken,
			EventID:      uuid.NewString(),
			Timestamp:    time.Now().UTC().UnixMicro(),
			Payload:      *bytes.NewBufferString("log content: " + uuid.NewString()),
		}
		_, _, err := clientEvent.HandleEvent(req)
		if err != nil {
			b.Log(i, err)
			failed++
		}
	}
	b.Logf("summary %d/%d\n", failed, b.N)
}
