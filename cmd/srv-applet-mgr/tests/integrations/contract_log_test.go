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

func TestContractLogAPIs(t *testing.T) {
	var (
		ctx           = requires.Context()
		client        = requires.AuthClient()
		projectName   = "test_project"
		contractLogID types.SFID
		baseReq       = applet_mgr.CreateContractLog{}
	)

	baseReq.ProjectName = projectName
	baseReq.CreateContractLogReq.BlockStart = 20299310
	baseReq.CreateContractLogReq.ChainID = 4690
	baseReq.CreateContractLogReq.EventType = "DEFAULT"
	baseReq.CreateContractLogReq.ContractAddress = "0x1AA325E5144f763a520867c56FC77cC1411430d0"
	baseReq.CreateContractLogReq.Topic0 = "0x9ffdf0136249d99680088653555755221714868b4f7ca1ff7d8523e3bef1dc4a"

	t.Logf("random a project name: %s", projectName)

	t.Run("ContractLog", func(t *testing.T) {
		t.Run("#CreateContractLog", func(t *testing.T) {
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

				// create contract log monitor
				{
					req := baseReq
					rsp, _, err := client.CreateContractLog(&req)

					NewWithT(t).Expect(err).To(BeNil())
					contractLogID = rsp.ContractLogID
				}

				// check contract log is created
				{
					_, err := blockchain.GetContractLogBySFID(ctx, contractLogID)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// remove contract log
				{
					req := &applet_mgr.RemoveContractLog{
						ProjectName:   projectName,
						ContractLogID: contractLogID,
					}
					_, err := client.RemoveContractLog(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// check contract log is removed
				{
					_, err := blockchain.GetContractLogBySFID(ctx, contractLogID)
					requires.CheckError(t, err, status.ContractLogNotFound)
				}

				// remove project
				{
					req := &applet_mgr.RemoveProject{ProjectName: projectName}
					_, err := client.RemoveProject(req)
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
			t.Run("#InvalidContractLogParam", func(t *testing.T) {
				// contract address is empty
				{
					req := baseReq
					req.CreateContractLogReq.ContractAddress = ""

					_, _, err := client.CreateContractLog(&req)
					requires.CheckError(t, err, &statusx.StatusErr{
						Key: "badRequest",
					})
				}

				// block start is empty
				{
					req := baseReq
					req.CreateContractLogReq.BlockStart = 0
					_, _, err := client.CreateContractLog(&req)
					requires.CheckError(t, err, &statusx.StatusErr{
						Key: "badRequest",
					})
				}

				// chain id is empty
				{
					req := baseReq
					req.CreateContractLogReq.ChainID = 0
					_, _, err := client.CreateContractLog(&req)
					requires.CheckError(t, err, &statusx.StatusErr{
						Key: "badRequest",
					})
				}
			})
		})
	})
}
