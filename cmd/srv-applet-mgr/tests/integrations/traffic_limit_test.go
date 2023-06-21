package integrations

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/clients/applet_mgr"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/requires"
	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestTrafficLimitAPIs(t *testing.T) {
	var (
		client      = requires.AuthClient()
		projectName = "test_traffic_limit_project"

		projectID      types.SFID
		trafficLimitID types.SFID
	)

	t.Logf("random a project name: %s, use this name create a project.", projectName)

	{
		req := &applet_mgr.CreateProject{}
		req.CreateReq.Name = projectName

		rsp, _, err := client.CreateProject(req)
		if err != nil {
			panic(err)
		}
		projectID = rsp.ProjectID
	}

	defer func() {
		req := &applet_mgr.RemoveProject{ProjectName: projectName}
		_, err := client.RemoveProject(req)
		if err != nil {
			panic(err)
		}
	}()

	t.Run("TrafficLimit", func(t *testing.T) {
		t.Run("#CreateTrafficLimit", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				// create trafficLimit
				{
					req := &applet_mgr.CreateTrafficLimit{ProjectName: projectName}
					req.CreateReq.Threshold = 2
					req.CreateReq.Duration = base.Duration(10 * time.Minute)
					req.CreateReq.ApiType = enums.TRAFFIC_LIMIT_TYPE__EVENT

					rsp, _, err := client.CreateTrafficLimit(req)
					NewWithT(t).Expect(err).To(BeNil())
					trafficLimitID = rsp.TrafficLimitID
				}

				// get trafficLimit
				{
					req := &applet_mgr.GetTrafficLimit{TrafficLimitID: trafficLimitID}
					rsp, _, err := client.GetTrafficLimit(req)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.ProjectID).To(Equal(projectID))
				}

				// update trafficLimit
				{
					req := &applet_mgr.UpdateTrafficLimit{ProjectName: projectName, TrafficLimitID: trafficLimitID}
					req.UpdateReq.Threshold = 3
					req.UpdateReq.Duration = base.Duration(3 * time.Minute)
					req.UpdateReq.ApiType = enums.TRAFFIC_LIMIT_TYPE__EVENT

					_, _, err := client.UpdateTrafficLimit(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// remove trafficLimit
				{
					req := &applet_mgr.RemoveTrafficLimit{ProjectName: projectName, TrafficLimitID: trafficLimitID}
					_, err := client.RemoveTrafficLimit(req)
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
		})
	})

	t.Run("BatchTrafficLimit", func(t *testing.T) {
		t.Run("#CreateTrafficLimits", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				// prepare data
				{
					req := &applet_mgr.CreateTrafficLimit{ProjectName: projectName}
					req.CreateReq.Threshold = 2
					req.CreateReq.Duration = base.Duration(2 * time.Minute)
					req.CreateReq.ApiType = enums.TRAFFIC_LIMIT_TYPE__EVENT

					_, _, err := client.CreateTrafficLimit(req)
					NewWithT(t).Expect(err).To(BeNil())

					req = &applet_mgr.CreateTrafficLimit{ProjectName: projectName}
					req.CreateReq.Threshold = 3
					req.CreateReq.Duration = base.Duration(3 * time.Minute)
					req.CreateReq.ApiType = enums.TRAFFIC_LIMIT_TYPE__BLOCKCHAIN

					_, _, err = client.CreateTrafficLimit(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// get list trafficLimit
				{
					req := &applet_mgr.ListTrafficLimit{ProjectName: projectName}
					rsp, _, err := client.ListTrafficLimit(req)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(2).To(Equal(int(rsp.Total)))
				}

				// remove batch trafficLimit
				{
					req := &applet_mgr.BatchRemoveTrafficLimit{ProjectName: projectName}
					_, err := client.BatchRemoveTrafficLimit(req)
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
		})
	})
}
