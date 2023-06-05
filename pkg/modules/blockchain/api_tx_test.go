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

func TestChainTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		db  = mock.NewMockDBExecutor(ctrl)
		ctx = contextx.WithContextCompose(
			types.WithMonitorDBExecutorContext(db),
			confid.WithSFIDGeneratorContext(confid.MustNewSFIDGenerator()),
		)(context.Background())
		req = CreateChainTxReq{
			ProjectName: "test_project",
			ChainTxInfo: models.ChainTxInfo{},
		}
	)

	t.Run("Create", func(t *testing.T) {

		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil)
			db.EXPECT().Exec(gomock.Any()).Return(nil, nil)

			_, err := CreateChainTx(ctx, &req)
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#ChainIDNotExist", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(sqlx.NewSqlError(sqlx.SqlErrTypeNotFound, ""))

			_, err := CreateChainTx(ctx, &req)
			NewWithT(t).Expect(err).To(Equal(status.BlockchainNotFound))
		})

		t.Run("#ChainTxConflict", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil)
			db.EXPECT().Exec(gomock.Any()).Return(nil, sqlx.NewSqlError(sqlx.SqlErrTypeConflict, ""))

			_, err := CreateChainTx(ctx, &req)
			NewWithT(t).Expect(err).To(Equal(status.ChainTxConflict))
		})
	})

	t.Run("Get", func(t *testing.T) {

		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil)

			_, err := GetChainTxBySFID(ctx, 1)
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#ChainTxNotFound", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(sqlx.NewSqlError(sqlx.SqlErrTypeNotFound, ""))

			_, err := GetChainTxBySFID(ctx, 1)
			NewWithT(t).Expect(err).To(Equal(status.ChainTxNotFound))
		})
	})

	t.Run("Remove", func(t *testing.T) {

		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
			db.EXPECT().Exec(gomock.Any()).Return(nil, nil)

			err := RemoveChainTxBySFID(ctx, 1)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})
}
