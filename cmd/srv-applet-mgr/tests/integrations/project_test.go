package integrations

import (
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/clients/applet_mgr"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/requires"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestProjectAPIs(t *testing.T) {
	defer requires.Serve()()

	var (
		ctx         = requires.Context()
		client      = requires.AuthClient()
		projectName = "test_project"
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
					_, err := config.GetByRelAndType(ctx, projectID, enums.CONFIG_TYPE__PROJECT_DATABASE)
					requires.CheckError(t, err, status.ConfigNotFound)
				}

				{
					_, err := config.GetByRelAndType(ctx, projectID, enums.CONFIG_TYPE__PROJECT_ENV)
					requires.CheckError(t, err, status.ConfigNotFound)
				}
			})
			t.Run("#InvalidProjectName", func(t *testing.T) {
				// project name is empty
				{
					req := &applet_mgr.CreateProject{}

					_, _, err := client.CreateProject(req)
					requires.CheckError(t, err, &statusx.StatusErr{
						Key: "badRequest",
					})
				}

				{
					req := &applet_mgr.CreateProject{}
					req.CreateReq.Name = strings.Repeat("a", 33)

					_, _, err := client.CreateProject(req)
					requires.CheckError(t, err, &statusx.StatusErr{
						Key: "badRequest",
					})
				}

			})
		})
	})
}
