package publisher

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	confjwt "github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/mock"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestPublisher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		db  = mock.NewMockDBExecutor(ctrl)
		idg = confid.MustNewSFIDGenerator()
		ctx = contextx.WithContextCompose(
			conflog.WithLoggerContext(conflog.Std()),
			types.WithMgrDBExecutorContext(db),
			confid.WithSFIDGeneratorContext(idg),
			types.WithAccountContext(&models.Account{
				RelAccount: models.RelAccount{AccountID: idg.MustGenSFID()},
			}),
			types.WithProjectContext(&models.Project{
				RelProject:  models.RelProject{ProjectID: idg.MustGenSFID()},
				ProjectName: models.ProjectName{Name: "test_project_for_publisher_unit_test"},
			}),
			confjwt.WithConfContext(&confjwt.Jwt{
				Issuer:  "w3bstream_test",
				SignKey: "xxxx",
			}),
		)(context.Background())
	)

	t.Run("Create", func(t *testing.T) {

		req := CreateReq{
			Name: "unit_test_publisher",
			Key:  "unit_test_publisher_01",
		}
		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			db.EXPECT().Exec(gomock.Any()).Return(nil, nil)

			_, err := Create(ctx, &req)
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#PublisherConflict", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			db.EXPECT().Exec(gomock.Any()).Return(nil, sqlx.NewSqlError(sqlx.SqlErrTypeConflict, ""))

			_, err := Create(ctx, &req)
			NewWithT(t).Expect(err).To(Equal(status.PublisherConflict))
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("#Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).Times(2)
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(2)

			// GetBySFID
			{
				_, err := GetBySFID(ctx, idg.MustGenSFID())
				NewWithT(t).Expect(err).To(BeNil())
			}

			// GetByProjectAndKey
			{
				_, err := GetByProjectAndKey(ctx, idg.MustGenSFID(), "publisher_key")
				NewWithT(t).Expect(err).To(BeNil())
			}
		})

		t.Run("#PublisherNotFound", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{}).Times(2)
			db.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(sqlx.NewSqlError(sqlx.SqlErrTypeNotFound, "")).Times(2)

			// GetBySFID
			{
				_, err := GetBySFID(ctx, 1)
				NewWithT(t).Expect(err).To(Equal(status.PublisherNotFound))
			}

			// GetByProjectAndKey
			{
				_, err := GetByProjectAndKey(ctx, idg.MustGenSFID(), "publisher_key")
				NewWithT(t).Expect(err).To(Equal(status.PublisherNotFound))
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

	t.Run("Remove", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			db.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			db.EXPECT().Exec(gomock.Any()).Return(driver.RowsAffected(1), nil)

			err := Remove(ctx, types.MustAccountFromContext(ctx), &CondArgs{})
			NewWithT(t).Expect(err).To(BeNil())
		})
	})
}
