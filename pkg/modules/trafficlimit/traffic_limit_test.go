package trafficlimit

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
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestTrafficLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		db  = mock.NewMockDBExecutor(ctrl)
		idg = confid.MustNewSFIDGenerator()
		ctx = contextx.WithContextCompose(
			types.WithMgrDBExecutorContext(db),
			confid.WithSFIDGeneratorContext(idg),
			types.WithProjectContext(&models.Project{
				RelProject:  models.RelProject{ProjectID: idg.MustGenSFID()},
				ProjectName: models.ProjectName{Name: "test_project_for_traffic_limit_unit_test"},
			}),
		)(context.Background())
	)

	t.Run("Get", func(t *testing.T) {
		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).Times(2)
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(2)

			// Get
			{
				_, err := GetBySFID(ctx, idg.MustGenSFID())
				NewWithT(t).Expect(err).To(BeNil())
			}

			// GetByProjectAndTypeMustDB
			{
				_, err := GetByProjectAndTypeMustDB(ctx, 1, enums.TRAFFIC_LIMIT_TYPE__EVENT)
				NewWithT(t).Expect(err).To(BeNil())
			}
		})

		t.Run("#TrafficLimitNotFound", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).Times(2)
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(sqlx.NewSqlError(sqlx.SqlErrTypeNotFound, "")).Times(2)

			// Get
			{
				_, err := GetBySFID(ctx, 1)
				NewWithT(t).Expect(err).To(Equal(status.TrafficLimitNotFound))
			}

			// GetByProjectAndTypeMustDB
			{
				_, err := GetByProjectAndTypeMustDB(ctx, 1, enums.TRAFFIC_LIMIT_TYPE__EVENT)
				NewWithT(t).Expect(err).To(Equal(status.TrafficLimitNotFound))
			}
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).Times(6)
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(5)

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

			{
				_, err := ListDetail(ctx, &ListReq{})
				NewWithT(t).Expect(err).To(BeNil())
			}
		})
	})
}
