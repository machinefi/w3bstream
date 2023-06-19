package cronjob

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	mock_sqlx "github.com/machinefi/w3bstream/pkg/test/mock_depends_kit_sqlx"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestCronJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		db  = mock_sqlx.NewMockDBExecutor(ctrl)
		prj = &models.Project{
			RelProject: models.RelProject{ProjectID: 1},
		}
		ctx = contextx.WithContextCompose(
			types.WithMgrDBExecutorContext(db),
			confid.WithSFIDGeneratorContext(confid.MustNewSFIDGenerator()),
			types.WithProjectContext(prj),
		)(context.Background())
	)

	db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()

	t.Run("Create", func(t *testing.T) {
		req := CreateReq{
			CronJobInfo: models.CronJobInfo{
				CronExpressions: "* * * * *",
			},
		}

		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().Exec(gomock.Any()).Return(nil, nil)

			_, err := Create(ctx, &req)
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#InvalidCronExpressions", func(t *testing.T) {
			req := req
			req.CronExpressions = "*"
			_, err := Create(ctx, &req)
			NewWithT(t).Expect(err.Error()).To(ContainSubstring(status.InvalidCronExpressions.Error()))
		})

		t.Run("#CronJobConflict", func(t *testing.T) {
			db.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrConflict)

			_, err := Create(ctx, &req)
			NewWithT(t).Expect(err).To(Equal(status.CronJobConflict))
		})
	})

	t.Run("Get", func(t *testing.T) {

		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil)

			_, err := GetBySFID(ctx, 1)
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#CronJobNotFound", func(t *testing.T) {
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound)

			_, err := GetBySFID(ctx, 1)
			NewWithT(t).Expect(err).To(Equal(status.CronJobNotFound))
		})
	})

	t.Run("List", func(t *testing.T) {
		req := ListReq{}

		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(2)

			_, err := List(ctx, &req)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("Remove", func(t *testing.T) {

		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().Exec(gomock.Any()).Return(nil, nil)

			err := RemoveBySFID(ctx, 1)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})
}
