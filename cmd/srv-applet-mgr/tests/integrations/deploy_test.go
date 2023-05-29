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
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestDeployAPIs(t *testing.T) {
	var (
		ctx         = requires.Context()
		client      = requires.AuthClient()
		projectName = "test_project_for_deploy"
		appletID    types.SFID
		instanceID  types.SFID
	)

	t.Logf("random a project name: %s", projectName)

	t.Run("Deploy", func(t *testing.T) {
		t.Run("#CreateInstance", func(t *testing.T) {
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
					instanceID = rsp.Instance.InstanceID
				}

				// check instance is created
				{
					_, err := deploy.GetByAppletSFID(ctx, appletID)
					NewWithT(t).Expect(err).To(BeNil())

					_, err = deploy.GetBySFID(ctx, instanceID)
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

				// check instance is removed
				{
					_, err := deploy.GetBySFID(ctx, instanceID)
					requires.CheckError(t, err, status.InstanceNotFound)
				}

				// remove project
				{
					req := &applet_mgr.RemoveProject{ProjectName: projectName}
					_, err := client.RemoveProject(req)
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
		})
		t.Run("#ControlInstance", func(t *testing.T) {
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
					instanceID = rsp.Instance.InstanceID
				}

				// check instance is started
				{
					req := &applet_mgr.GetInstanceByInstanceID{
						InstanceID: instanceID,
					}
					rsp, _, err := client.GetInstanceByInstanceID(req)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.State).
						To(Equal(enums.INSTANCE_STATE__STARTED))
				}

				// put instance to stoped
				{
					req := &applet_mgr.ControlInstance{
						InstanceID: instanceID,
						Cmd:        enums.DEPLOY_CMD__HUNGUP,
					}
					_, err := client.ControlInstance(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// check instance is stopped
				{
					req := &applet_mgr.GetInstanceByInstanceID{
						InstanceID: instanceID,
					}
					rsp, _, err := client.GetInstanceByInstanceID(req)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.State).
						To(Equal(enums.INSTANCE_STATE__STOPPED))
				}

				// remove applet
				{
					req := &applet_mgr.RemoveApplet{
						AppletID: appletID,
					}
					_, err := client.RemoveApplet(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// check instance is removed
				{
					_, err := deploy.GetBySFID(ctx, instanceID)
					requires.CheckError(t, err, status.InstanceNotFound)
				}

				// remove project
				{
					req := &applet_mgr.RemoveProject{ProjectName: projectName}
					_, err := client.RemoveProject(req)
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
			t.Run("#InvalidDeployParam", func(t *testing.T) {

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
					instanceID = rsp.Instance.InstanceID
				}

				// control cmd is empty
				{
					req := &applet_mgr.ControlInstance{
						InstanceID: instanceID,
					}
					_, err := client.ControlInstance(req)
					requires.CheckError(t, err, &statusx.StatusErr{
						Key: "badRequest",
					})
				}

				// control cmd is illegal
				{
					req := &applet_mgr.ControlInstance{
						InstanceID: instanceID,
						Cmd:        enums.DeployCmd(1),
					}
					_, err := client.ControlInstance(req)
					requires.CheckError(t, err, &statusx.StatusErr{
						Key: "badRequest",
					})
				}

				// remove applet
				{
					req := &applet_mgr.RemoveApplet{
						AppletID: appletID,
					}
					_, err := client.RemoveApplet(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// check instance is removed
				{
					_, err := deploy.GetBySFID(ctx, instanceID)
					requires.CheckError(t, err, status.InstanceNotFound)
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
