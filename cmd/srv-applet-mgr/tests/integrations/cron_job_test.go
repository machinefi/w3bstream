package integrations

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/clients/applet_mgr"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/requires"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/cronjob"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestCronJobAPIs(t *testing.T) {
	var (
		ctx         = requires.Context()
		client      = requires.AuthClient()
		projectName = "test_project_for_cron_job"
		cronJobID   types.SFID
		projectID   types.SFID
	)

	t.Logf("random a project name: %s", projectName)

	t.Run("CronJob", func(t *testing.T) {
		t.Run("#CreateCronJob", func(t *testing.T) {
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

				// create cron job
				{
					req := applet_mgr.CreateCronJob{
						ProjectID: projectID,
					}
					req.CreateReq.EventType = "default"
					req.CreateReq.CronExpressions = "*/3 * * * *"
					rsp, _, err := client.CreateCronJob(&req)

					NewWithT(t).Expect(err).To(BeNil())
					cronJobID = rsp.CronJobID
				}

				// check cron job is created
				{
					_, err := cronjob.GetBySFID(ctx, cronJobID)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// remove cron job
				{
					req := &applet_mgr.RemoveCronJob{
						CronJobID: cronJobID,
					}
					_, err := client.RemoveCronJob(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// check cron job is removed
				{
					_, err := cronjob.GetBySFID(ctx, cronJobID)
					requires.CheckError(t, err, status.CronJobNotFound)
				}

				// remove project
				{
					req := &applet_mgr.RemoveProject{ProjectName: projectName}
					_, err := client.RemoveProject(req)
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
			t.Run("#InvalidCronJobParam", func(t *testing.T) {

				// create project without user defined config(database/env)
				{
					req := &applet_mgr.CreateProject{}
					req.CreateReq.Name = projectName

					rsp, _, err := client.CreateProject(req)

					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.Name).To(Equal(projectName))
					projectID = rsp.ProjectID
				}

				// cron expressions is empty
				{
					req := applet_mgr.CreateCronJob{
						ProjectID: projectID,
					}

					_, _, err := client.CreateCronJob(&req)
					requires.CheckError(t, err, &statusx.StatusErr{
						Key: "badRequest",
					})
				}

				// cron expressions is illegal
				{
					req := applet_mgr.CreateCronJob{
						ProjectID: projectID,
					}
					req.CreateReq.CronExpressions = "*/3 * * * * *"

					_, _, err := client.CreateCronJob(&req)
					requires.CheckError(t, err, &statusx.StatusErr{
						Key: "InvalidCronExpressions",
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
