package integrations

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/__test__/clients/applet_mgr"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/__test__/requires"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestProjectAPIs(t *testing.T) {
	defer requires.Serve()()
	var (
		client      = requires.AuthClient()
		projectName = "testdemo"
		projectID   types.SFID
	)

	t.Logf("random a project name: %s", projectName)

	t.Run("Project", func(t *testing.T) {
		t.Run("#CreateProject", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				// create project without user defined config(database/env)
				{
					req := &applet_mgr.CreateProject{}
					req.CreateReq.Name = projectName

					rsp, _, err := client.CreateProject(req)

					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.Name).To(Equal(projectName))
					projectID = rsp.ProjectID
					defer requires.DropTempWasmDatabase(projectID)
				}

				// check project default config
				{
					req := &applet_mgr.GetProjectSchema{ProjectName: projectName}
					rsp, _, err := client.GetProjectSchema(req)

					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.ConfigType()).
						To(Equal(enums.CONFIG_TYPE__PROJECT_DATABASE))
				}

				{
					req := &applet_mgr.GetProjectEnv{ProjectName: projectName}

					rsp, _, err := client.GetProjectEnv(req)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.ConfigType()).
						To(Equal(enums.CONFIG_TYPE__PROJECT_ENV))
				}

				// remove project
				{
					req := &applet_mgr.RemoveProject{ProjectName: projectName}
					_, err := client.RemoveProject(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// check project config is removed
				{
					req := &applet_mgr.GetProjectSchema{ProjectName: projectName}
					_, _, err := client.GetProjectSchema(req)

					requires.CheckError(t, err, status.ProjectNotFound)
				}

				{
					req := &applet_mgr.GetProjectEnv{ProjectName: projectName}
					_, _, err := client.GetProjectEnv(req)

					requires.CheckError(t, err, status.ProjectNotFound)
				}
			})
		})
	})
}
