package blockchain

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/mock"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestChainHeight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		db         = mock.NewMockDBExecutor(ctrl)
		ethClients = &types.ETHClientConfig{
			Clients: map[uint32]string{4690: "https://babel-api.testnet.iotex.io"},
		}
		ctx = contextx.WithContextCompose(
			types.WithMonitorDBExecutorContext(db),
			confid.WithSFIDGeneratorContext(confid.MustNewSFIDGenerator()),
			types.WithETHClientConfigContext(ethClients),
		)(context.Background())
		req = CreateChainHeightReq{
			ProjectName: "test_project",
			ChainHeightInfo: models.ChainHeightInfo{
				ChainID: 4690,
			},
		}
	)

	t.Run("Create", func(t *testing.T) {

		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			db.EXPECT().Exec(gomock.Any()).Return(nil, nil)

			_, err := CreateChainHeight(ctx, &req)
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#ChainIDNotExist", func(t *testing.T) {
			req := req
			req.ChainID = 1
			_, err := CreateChainHeight(ctx, &req)
			NewWithT(t).Expect(err).To(Equal(status.BlockchainNotFound))
		})

		t.Run("#ChainHeightConflict", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			db.EXPECT().Exec(gomock.Any()).Return(nil, sqlx.NewSqlError(sqlx.SqlErrTypeConflict, ""))

			_, err := CreateChainHeight(ctx, &req)
			NewWithT(t).Expect(err).To(Equal(status.ChainHeightConflict))
		})
	})

	t.Run("Get", func(t *testing.T) {

		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil)

			_, err := GetChainHeightBySFID(ctx, 1)
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#ChainHeightNotFound", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(sqlx.NewSqlError(sqlx.SqlErrTypeNotFound, ""))

			_, err := GetChainHeightBySFID(ctx, 1)
			NewWithT(t).Expect(err).To(Equal(status.ChainHeightNotFound))
		})
	})

	t.Run("Remove", func(t *testing.T) {

		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			db.EXPECT().Exec(gomock.Any()).Return(nil, nil)

			err := RemoveChainHeightBySFID(ctx, 1)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})
}
