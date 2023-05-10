package integrations

import (
	"bytes"
	_ "embed"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	// . "github.com/onsi/gomega"

	"github.com/google/uuid"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/clients/applet_mgr"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/requires"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
)

//go:embed testdata/log.wasm
var code []byte

func BenchmarkEventHandling(b *testing.B) {
	defer requires.Serve()()

	var (
		client         = requires.AuthClient()
		clientEvent    *applet_mgr.Client
		projectName    = "testdemo"
		publisherToken string
	)

	{
		req := &applet_mgr.CreateProject{}
		req.CreateReq.Name = projectName

		_, _, err := client.CreateProject(req)
		if err != nil {
			b.Log(err)
			return
		}
	}
	defer func() {
		_, _ = client.RemoveProject(&applet_mgr.RemoveProject{
			ProjectName: projectName,
		})
	}()

	{
		cwd, _ := os.Getwd()
		filename := path.Join(cwd, "../testdata/log.wasm")
		req := &applet_mgr.CreateApplet{
			ProjectName: projectName,
		}
		req.CreateReq.File = transformer.MustNewFileHeader("file", filename, bytes.NewBuffer(code))
		req.CreateReq.Info = applet_mgr.GithubComMachinefiW3BstreamPkgModulesAppletInfo{
			AppletName: "log",
			WasmName:   "log.wasm",
		}

		_, _, err := client.CreateApplet(req)
		if err != nil {
			b.Log(err)
			return
		}
	}

	{
		req := &applet_mgr.CreatePublisher{
			ProjectName: projectName,
		}
		req.CreateReq.Name = "test_publisher"
		req.CreateReq.Key = "mn_test_publisher"

		rsp, _, err := client.CreatePublisher(req)
		if err != nil {
			b.Log(err)
			return
		}
		publisherToken = rsp.Token
	}
	clientEvent = requires.ClientEvent()

	b.N = 1
	channel := strings.Join([]string{"aid", requires.AccountID.String(), projectName}, "_")
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
		}
		time.Sleep(time.Second)
	}
}
