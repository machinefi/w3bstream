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

func TestChainTxAPIs(t *testing.T) {
	defer requires.Serve()()

	var (
		ctx         = requires.Context()
		client      = requires.AuthClient()
		projectName = "test_project"
		chainTxID   types.SFID
		baseReq     = applet_mgr.CreateChainTx{}
	)

	baseReq.ProjectName = projectName
	baseReq.CreateChainTxReq.ChainID = 4690
	baseReq.CreateChainTxReq.EventType = "DEFAULT"
	baseReq.CreateChainTxReq.TxAddress = "85cfc682ccbf276cc80a1242e69e1d1a8ab9295465c6fc9e4d3f433363ec3ccd"

	t.Logf("random a project name: %s", projectName)

	t.Run("ChainTx", func(t *testing.T) {
		t.Run("#CreateChainTx", func(t *testing.T) {
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

				// create chain tx monitor
				{
					req := baseReq
					rsp, _, err := client.CreateChainTx(&req)

					NewWithT(t).Expect(err).To(BeNil())
					chainTxID = rsp.ChainTxID
				}

				// check chain tx is created
				{
					_, err := blockchain.GetChainTxBySFID(ctx, chainTxID)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// remove chain tx
				{
					req := &applet_mgr.RemoveChainTx{
						ProjectName: projectName,
						ChainTxID:   chainTxID,
					}
					_, err := client.RemoveChainTx(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// check chain tx is removed
				{
					_, err := blockchain.GetChainTxBySFID(ctx, chainTxID)
					requires.CheckError(t, err, status.ChainTxNotFound)
				}

				// remove project
				{
					req := &applet_mgr.RemoveProject{ProjectName: projectName}
					_, err := client.RemoveProject(req)
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
			t.Run("#InvalidChainTxParam", func(t *testing.T) {
				// chain tx is empty
				{
					req := baseReq
					req.CreateChainTxReq.TxAddress = ""

					_, _, err := client.CreateChainTx(&req)
					requires.CheckError(t, err, &statusx.StatusErr{
						Key: "badRequest",
					})
				}

				// chain id is empty
				{
					req := baseReq
					req.CreateChainTxReq.ChainID = 0
					_, _, err := client.CreateChainTx(&req)
					requires.CheckError(t, err, &statusx.StatusErr{
						Key: "badRequest",
					})
				}
			})
		})
	})
}
