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

	t.Run("PrepareProject", func(t *testing.T) {
		t.Run("#CreateProject", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				// create project without user defined config(database/env)
				{
					req := &applet_mgr.CreateProject{}
					req.CreateReq.Name = projectName

					rsp, _, err := client.CreateProject(req)

					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.Name).To(Equal(projectName))
				}
			})
		})
		t.Run("#CreateApplet", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				{
					cwd, err := os.Getwd()
					NewWithT(t).Expect(err).To(BeNil())

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

					NewWithT(t).Expect(err).To(BeNil())
					appletID = rsp.AppletID
				}
			})
		})
		t.Run("#GetResourceID", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				applet, err := applet.GetBySFID(ctx, appletID)
				NewWithT(t).Expect(err).To(BeNil())
				resourceID = applet.ResourceID

			})
		})
	})

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

	// clear project info
	t.Run("ClearProject", func(t *testing.T) {
		t.Run("#DeleteProject", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				// remove project
				{
					req := &applet_mgr.RemoveProject{ProjectName: projectName}
					_, err := client.RemoveProject(req)
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
		})
	})

}
