package blockchain

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/models"
)

func TestContractParseTopic(t *testing.T) {
	c := contract{}

	t.Run("Empty String", func(t *testing.T) {
		r := c.parseTopic("")
		NewWithT(t).Expect(r).To(BeNil())
	})

	t.Run("Success", func(t *testing.T) {
		str := "0xb9fb64ccf647f3e7ba45742b97b6b8e464a822c67817276accb7b1f905d292a2"
		r := c.parseTopic(str)
		NewWithT(t).Expect(*r).To(Equal(common.HexToHash(str)))
	})
}

func TestContractGetTopic(t *testing.T) {
	c := contract{}
	tStr := "0xb9fb64ccf647f3e7ba45742b97b6b8e464a822c67817276accb7b1f905d292a2"
	tHash := *c.parseTopic(tStr)

	t.Run("Empty Topic", func(t *testing.T) {
		rs, mrs := c.getTopic(nil)
		NewWithT(t).Expect(len(rs)).To(Equal(0))
		NewWithT(t).Expect(len(mrs)).To(Equal(0))
	})

	t.Run("Have Topic 0", func(t *testing.T) {
		ms := []*models.ContractLog{{
			ContractLogData: models.ContractLogData{
				ContractLogInfo: models.ContractLogInfo{
					Topic0: tStr,
				},
			},
		}}
		rs, mrs := c.getTopic(ms)
		NewWithT(t).Expect(len(rs)).To(Equal(1))
		NewWithT(t).Expect(len(mrs)).To(Equal(1))
		NewWithT(t).Expect(rs[0][0]).To(Equal(tHash))
	})

	t.Run("Have Topic 1", func(t *testing.T) {
		ms := []*models.ContractLog{{
			ContractLogData: models.ContractLogData{
				ContractLogInfo: models.ContractLogInfo{
					Topic1: tStr,
				},
			},
		}}
		rs, mrs := c.getTopic(ms)
		NewWithT(t).Expect(len(rs)).To(Equal(2))
		NewWithT(t).Expect(len(mrs)).To(Equal(1))
		NewWithT(t).Expect(rs[1][0]).To(Equal(tHash))
	})

	t.Run("Have Topic 2", func(t *testing.T) {
		ms := []*models.ContractLog{{
			ContractLogData: models.ContractLogData{
				ContractLogInfo: models.ContractLogInfo{
					Topic2: tStr,
				},
			},
		}}
		rs, mrs := c.getTopic(ms)
		NewWithT(t).Expect(len(rs)).To(Equal(3))
		NewWithT(t).Expect(len(mrs)).To(Equal(1))
		NewWithT(t).Expect(rs[2][0]).To(Equal(tHash))
	})

	t.Run("Have Topic 3", func(t *testing.T) {
		ms := []*models.ContractLog{{
			ContractLogData: models.ContractLogData{
				ContractLogInfo: models.ContractLogInfo{
					Topic3: tStr,
				},
			},
		}}
		rs, mrs := c.getTopic(ms)
		NewWithT(t).Expect(len(rs)).To(Equal(4))
		NewWithT(t).Expect(len(mrs)).To(Equal(1))
		NewWithT(t).Expect(rs[3][0]).To(Equal(tHash))
	})

	t.Run("Have all Topics", func(t *testing.T) {
		ms := []*models.ContractLog{{
			ContractLogData: models.ContractLogData{
				ContractLogInfo: models.ContractLogInfo{
					Topic0: tStr,
					Topic1: tStr,
					Topic2: tStr,
					Topic3: tStr,
				},
			},
		}}
		rs, mrs := c.getTopic(ms)
		NewWithT(t).Expect(len(rs)).To(Equal(4))
		NewWithT(t).Expect(len(mrs)).To(Equal(1))
		NewWithT(t).Expect(rs[0][0]).To(Equal(tHash))
		NewWithT(t).Expect(rs[1][0]).To(Equal(tHash))
		NewWithT(t).Expect(rs[2][0]).To(Equal(tHash))
		NewWithT(t).Expect(rs[3][0]).To(Equal(tHash))
	})
}

func TestContractGetAddresses(t *testing.T) {
	c := contract{}
	aStr := "0x1AA325E5144f763a520867c56FC77cC1411430d0"
	aHash := common.HexToAddress(aStr)

	t.Run("Empty Address", func(t *testing.T) {
		rs, mrs := c.getAddresses(nil)
		NewWithT(t).Expect(len(rs)).To(Equal(0))
		NewWithT(t).Expect(len(mrs)).To(Equal(0))
	})

	t.Run("Have Address", func(t *testing.T) {
		ms := []*models.ContractLog{{
			ContractLogData: models.ContractLogData{
				ContractLogInfo: models.ContractLogInfo{
					ContractAddress: aStr,
				},
			},
		}}
		rs, mrs := c.getAddresses(ms)
		NewWithT(t).Expect(len(rs)).To(Equal(1))
		NewWithT(t).Expect(len(mrs)).To(Equal(1))
		NewWithT(t).Expect(rs[0]).To(Equal(aHash))
	})
}

func TestContractGetExpectedContractLogs(t *testing.T) {
	c := contract{}
	tStr := "0xb9fb64ccf647f3e7ba45742b97b6b8e464a822c67817276accb7b1f905d292a2"
	aStr := "0x1AA325E5144f763a520867c56FC77cC1411430d0"
	tHash := *c.parseTopic(tStr)
	aHash := common.HexToAddress(aStr)

	m := &models.ContractLog{
		ContractLogData: models.ContractLogData{
			ContractLogInfo: models.ContractLogInfo{
				Topic0:          tStr,
				ContractAddress: aStr,
			},
		},
	}
	ms := []*models.ContractLog{m}

	_, mts := c.getTopic(ms)
	_, mas := c.getAddresses(ms)

	t.Run("Match Failed", func(t *testing.T) {
		res := c.getExpectedContractLogs(&ethtypes.Log{}, mas, mts)
		NewWithT(t).Expect(len(res)).To(Equal(0))
	})

	t.Run("Success", func(t *testing.T) {
		log := &ethtypes.Log{
			Address: aHash,
			Topics:  []common.Hash{tHash},
		}

		res := c.getExpectedContractLogs(log, mas, mts)
		NewWithT(t).Expect(len(res)).To(Equal(1))
		NewWithT(t).Expect(res[0]).To(Equal(m))
	})
}

func TestContractPruneListChainGroups(t *testing.T) {
	c := contract{}
	ms := []*models.ContractLog{{
		ContractLogData: models.ContractLogData{
			ContractLogInfo: models.ContractLogInfo{
				BlockCurrent: 100,
			},
		},
	},
		{
			ContractLogData: models.ContractLogData{
				ContractLogInfo: models.ContractLogInfo{
					BlockCurrent: 100,
				},
			},
		},
	}
	gs := []*listChainGroup{{cs: ms}}

	t.Run("Have the same block current", func(t *testing.T) {
		c.pruneListChainGroups(gs)
		NewWithT(t).Expect(gs).To(Equal(gs))
	})

	t.Run("Don't have the same block current", func(t *testing.T) {
		ngs := gs
		ngs[0].cs = append(ngs[0].cs, &models.ContractLog{
			ContractLogData: models.ContractLogData{
				ContractLogInfo: models.ContractLogInfo{
					BlockCurrent: 200,
				},
			},
		},
		)
		c.pruneListChainGroups(gs)
		NewWithT(t).Expect(ngs[0].cs).To(Equal(gs[0].cs))
		NewWithT(t).Expect(ngs[0].toBlock).To(Equal(uint64(199)))
	})
}

func TestContractGroupContractLog(t *testing.T) {
	c := contract{}
	ms := []models.ContractLog{
		{
			ContractLogData: models.ContractLogData{
				ProjectName: "test1",
				ContractLogInfo: models.ContractLogInfo{
					ChainID:      1,
					BlockCurrent: 100,
				},
			},
		},
		{
			ContractLogData: models.ContractLogData{
				ProjectName: "test1",
				ContractLogInfo: models.ContractLogInfo{
					ChainID:      2,
					BlockCurrent: 100,
				},
			},
		},
	}

	gs := c.groupContractLog(ms)
	NewWithT(t).Expect(len(gs)).To(Equal(int(2)))
	NewWithT(t).Expect(gs[0].cs[0].ChainID).To(Equal(uint64(1)))
}
