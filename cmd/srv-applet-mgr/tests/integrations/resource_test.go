package integrations

import (
	"bytes"
	"os"
	"path"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/clients/applet_mgr"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/requires"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestResourceAPIs(t *testing.T) {
	var (
		ctx         = requires.Context()
		client      = requires.AuthClient()
		projectName = "test_resource_project"

		appletID   types.SFID
		resourceID types.SFID
	)

	t.Logf("random a project name: %s, use this name create a project and an applet.", projectName)

	{
		req := &applet_mgr.CreateProject{}
		req.CreateReq.Name = projectName

		_, _, err := client.CreateProject(req)
		if err != nil {
			panic(err)
		}
	}

	{
		cwd, err := os.Getwd()
		filename := path.Join(cwd, "../testdata/log.wasm")
		appletName := "testApplet"
		wasmName := "test.log"

		req := &applet_mgr.CreateApplet{ProjectName: projectName}
		req.CreateReq.File = transformer.MustNewFileHeader("file", filename, bytes.NewBuffer(code))
		req.CreateReq.Info = applet.Info{
			AppletName: appletName,
			WasmName:   wasmName,
		}

		rsp, _, err := client.CreateApplet(req)
		if err != nil {
			panic(err)
		}
		appletID = rsp.AppletID
	}

	{
		applet, err := applet.GetBySFID(ctx, appletID)
		if err != nil {
			panic(err)
		}
		resourceID = applet.ResourceID
	}

	defer func() {
		req := &applet_mgr.RemoveProject{ProjectName: projectName}
		_, err := client.RemoveProject(req)
		if err != nil {
			panic(err)
		}
	}()

	t.Logf("start test resource api.")

	t.Run("Resource", func(t *testing.T) {
		t.Run("#GetResuorce", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				// get resource
				{
					req := &applet_mgr.DownloadResource{ResourceID: resourceID}
					_, _, err := client.DownloadResource(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// get list strategies
				{
					req := &applet_mgr.ListResources{}
					_, _, err := client.ListResources(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// remove resource
				{
					req := &applet_mgr.RemoveResource{ResourceID: resourceID}
					_, err := client.RemoveResource(req)
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
		})
	})

}
