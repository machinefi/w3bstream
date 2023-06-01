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
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestAppletAPIs(t *testing.T) {
	var (
		ctx         = requires.Context()
		client      = requires.AuthClient()
		projectName = "test_project_for_applet"
		appletID    types.SFID
	)

	t.Logf("random a project name: %s", projectName)

	t.Run("Applet", func(t *testing.T) {
		t.Run("#CreateApplet", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				// create project without user defined config(database/env)
				{
					req := &applet_mgr.CreateProject{}
					req.CreateReq.Name = projectName

					rsp, _, err := client.CreateProject(req)

					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.Name).To(Equal(projectName))
				}

				// create applet
				{
					cwd, err := os.Getwd()
					NewWithT(t).Expect(err).To(BeNil())

					filename := path.Join(cwd, "../testdata/log.wasm")
					req := &applet_mgr.CreateApplet{
						ProjectName: projectName,
					}
					req.CreateReq.File = transformer.MustNewFileHeader("file", filename, bytes.NewBuffer(code))
					req.CreateReq.Info = applet_mgr.GithubComMachinefiW3BstreamPkgModulesAppletInfo{
						AppletName: "log",
						WasmName:   "log.wasm",
					}

					rsp, _, err := client.CreateApplet(req)
					NewWithT(t).Expect(err).To(BeNil())

					appletID = rsp.AppletID
				}

				// check applet is created
				{
					_, err := applet.GetBySFID(ctx, appletID)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// remove applet
				{
					req := &applet_mgr.RemoveApplet{
						AppletID: appletID,
					}
					_, err := client.RemoveApplet(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// check applet is removed
				{
					_, err := applet.GetBySFID(ctx, appletID)
					requires.CheckError(t, err, status.AppletNotFound)
				}

				// remove project
				{
					req := &applet_mgr.RemoveProject{ProjectName: projectName}
					_, err := client.RemoveProject(req)
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
			t.Run("#InvalidAppletParam", func(t *testing.T) {

				// create project without user defined config(database/env)
				{
					req := &applet_mgr.CreateProject{}
					req.CreateReq.Name = projectName

					rsp, _, err := client.CreateProject(req)

					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.Name).To(Equal(projectName))
				}

				// applet name is empty
				{
					cwd, err := os.Getwd()
					NewWithT(t).Expect(err).To(BeNil())

					filename := path.Join(cwd, "../testdata/log.wasm")
					req := &applet_mgr.CreateApplet{
						ProjectName: projectName,
					}
					req.CreateReq.File = transformer.MustNewFileHeader("file", filename, bytes.NewBuffer(code))
					req.CreateReq.Info = applet_mgr.GithubComMachinefiW3BstreamPkgModulesAppletInfo{
						WasmName: "log.wasm",
					}

					_, _, err = client.CreateApplet(req)
					requires.CheckError(t, err, &statusx.StatusErr{
						Key: "badRequest",
					})
				}

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
