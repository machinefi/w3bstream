package resource

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

func TestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		db  = mock.NewMockDBExecutor(ctrl)
		idg = confid.MustNewSFIDGenerator()
		ctx = contextx.WithContextCompose(
			types.WithMgrDBExecutorContext(db),
			confid.WithSFIDGeneratorContext(idg),
			types.WithAccountContext(&models.Account{
				RelAccount: models.RelAccount{AccountID: idg.MustGenSFID()},
			}),
		)(context.Background())
	)

	t.Run("Get", func(t *testing.T) {
		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).Times(2)
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(2)

			// GetBySFID
			{
				_, err := GetBySFID(ctx, idg.MustGenSFID())
				NewWithT(t).Expect(err).To(BeNil())
			}

			// GetByMd5
			{
				_, err := GetByMd5(ctx, "resource_md5")
				NewWithT(t).Expect(err).To(BeNil())
			}
		})

		t.Run("#ResourceNotFound", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).Times(2)
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(sqlx.NewSqlError(sqlx.SqlErrTypeNotFound, "")).Times(2)

			// GetBySFID
			{
				_, err := GetBySFID(ctx, 1)
				NewWithT(t).Expect(err).To(Equal(status.ResourceNotFound))
			}

			// GetByMd5
			{
				_, err := GetByMd5(ctx, "resource_md5")
				NewWithT(t).Expect(err).To(Equal(status.ResourceNotFound))
			}
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).Times(4)
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(2)

			_, err := List(ctx, &ListReq{})
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("GetOwnerByAccountAndSFID", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil)

			_, err := GetOwnerByAccountAndSFID(ctx, types.MustAccountFromContext(ctx).AccountID, 1)
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#ResourceNotFound", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(sqlx.NewSqlError(sqlx.SqlErrTypeNotFound, ""))

			_, err := GetOwnerByAccountAndSFID(ctx, types.MustAccountFromContext(ctx).AccountID, 1)
			NewWithT(t).Expect(err).To(Equal(status.ResourcePermNotFound))
		})
	})

	t.Run("RemoveOwnershipBySFID", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			db.EXPECT().Exec(gomock.Any()).Return(nil, nil)

			err := RemoveOwnershipBySFID(ctx, 1)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})
}
