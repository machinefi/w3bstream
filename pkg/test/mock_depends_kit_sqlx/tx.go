package mock_sqlx

import "github.com/golang/mock/gomock"

func NewTx(c *gomock.Controller) *Tx {
	return &Tx{
		MockDBExecutor: NewMockDBExecutor(c),
		MockTxExecutor: NewMockTxExecutor(c),
	}
}

type Tx struct {
	*MockDBExecutor
	*MockTxExecutor
}
