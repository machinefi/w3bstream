package strategy

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
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestStrategy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		db  = mock.NewMockDBExecutor(ctrl)
		idg = confid.MustNewSFIDGenerator()
		ctx = contextx.WithContextCompose(
			types.WithMgrDBExecutorContext(db),
			confid.WithSFIDGeneratorContext(idg),
		)(context.Background())
	)

	t.Run("Get", func(t *testing.T) {
		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil)

			_, err := GetBySFID(ctx, idg.MustGenSFID())
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#StrategyNotFound", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(sqlx.NewSqlError(sqlx.SqlErrTypeNotFound, ""))

			_, err := GetBySFID(ctx, 1)
			NewWithT(t).Expect(err).To(Equal(status.StrategyNotFound))
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).Times(12)
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(6)

			// List
			{
				_, err := List(ctx, &ListReq{})
				NewWithT(t).Expect(err).To(BeNil())
			}

			// ListByCond
			{
				_, err := ListByCond(ctx, &CondArgs{})
				NewWithT(t).Expect(err).To(BeNil())
			}

			// ListDetailByCond
			{
				_, err := ListDetailByCond(ctx, &CondArgs{}, (&ListReq{}).Addition())
				NewWithT(t).Expect(err).To(BeNil())
			}

			// ListDetail
			{
				_, err := ListDetail(ctx, &ListReq{})
				NewWithT(t).Expect(err).To(BeNil())
			}
		})
	})

	t.Run("Remove", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).Times(2)
			db.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(2)

			{
				err := RemoveBySFID(ctx, 1)
				NewWithT(t).Expect(err).To(BeNil())
			}

			// Remove
			{
				err := Remove(ctx, &CondArgs{})
				NewWithT(t).Expect(err).To(BeNil())
			}
		})
	})

	t.Run("Filter", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			_, err := FilterByProjectAndEvent(ctx, 1, "DEFAULT")
			NewWithT(t).Expect(err).To(BeNil())
		})
	})
}
