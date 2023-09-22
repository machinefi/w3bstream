package integrations

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/clients/applet_mgr"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/requires"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestChainHeightAPIs(t *testing.T) {
	var (
		ctx           = requires.Context()
		client        = requires.AuthClient()
		projectName   = "test_project"
		chainHeightID types.SFID
		baseReq       = applet_mgr.CreateChainHeight{}
	)

	baseReq.ProjectName = projectName
	baseReq.CreateChainHeightReq.ChainID = 4690
	baseReq.CreateChainHeightReq.EventType = "DEFAULT"
	baseReq.CreateChainHeightReq.Height = 1000

	t.Logf("random a project name: %s", projectName)

	t.Run("ChainHeight", func(t *testing.T) {
		t.Run("#CreateChainHeight", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				// create project without user defined config(database/env)
				{
					req := &applet_mgr.CreateProject{}
					req.CreateReq.Name = projectName

					rsp, _, err := client.CreateProject(req)

					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.Name).To(Equal(projectName))
				}

				// init blockchain config
				{
					err := blockchain.InitChainDB(ctx)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// create chain height monitor
				{
					req := baseReq
					rsp, _, err := client.CreateChainHeight(&req)

					NewWithT(t).Expect(err).To(BeNil())
					chainHeightID = rsp.ChainHeightID
				}

				// check chain height is created
				{
					_, err := blockchain.GetChainHeightBySFID(ctx, chainHeightID)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// remove chain height
				{
					req := &applet_mgr.RemoveChainHeight{
						ProjectName:   projectName,
						ChainHeightID: chainHeightID,
					}
					_, err := client.RemoveChainHeight(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// check chain height is removed
				{
					_, err := blockchain.GetChainHeightBySFID(ctx, chainHeightID)
					requires.CheckError(t, err, status.ChainHeightNotFound)
				}

				// remove project
				{
					req := &applet_mgr.RemoveProject{ProjectName: projectName}
					_, err := client.RemoveProject(req)
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
			t.Run("#InvalidChainHeightParam", func(t *testing.T) {
				// chain height is empty
				{
					req := baseReq
					req.CreateChainHeightReq.Height = 0

					_, _, err := client.CreateChainHeight(&req)
					requires.CheckError(t, err, &statusx.StatusErr{
						Key: "badRequest",
					})
				}

				// chain id is empty
				{
					req := baseReq
					req.CreateChainHeightReq.ChainID = 0
					_, _, err := client.CreateChainHeight(&req)
					requires.CheckError(t, err, &statusx.StatusErr{
						Key: "ProjectNotFound",
					})
				}
			})
		})
	})
}
