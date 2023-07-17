package publisher_test

import (
	"context"
	"database/sql/driver"
	"runtime"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	confjwt "github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	mock_sqlx "github.com/machinefi/w3bstream/pkg/test/mock_depends_kit_sqlx"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestPublisher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		db = &struct {
			*mock_sqlx.MockDBExecutor
			*mock_sqlx.MockTxExecutor
		}{
			MockDBExecutor: mock_sqlx.NewMockDBExecutor(ctrl),
			MockTxExecutor: mock_sqlx.NewMockTxExecutor(ctrl),
		}
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

	db.MockDBExecutor.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
	db.MockDBExecutor.EXPECT().Context().Return(ctx).AnyTimes()
	db.MockTxExecutor.EXPECT().IsTx().Return(true).AnyTimes()

	t.Run("Create", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}

		patch := gomonkey.NewPatches()
		defer patch.Reset()

		anyErr := errors.New("any")

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#CreateAccessKeyFailed", func(t *testing.T) {
				patch = patch.
					ApplyFunc(
						access_key.Create,
						func(_ context.Context, _ *access_key.CreateReq) (*access_key.CreateRsp, error) {
							return nil, anyErr
						},
					)
				_, err := publisher.Create(ctx, &publisher.CreateReq{})
				NewWithT(t).Expect(err).To(Equal(anyErr))
			})
			t.Run("#CreatePublisherFailed", func(t *testing.T) {
				patch = patch.
					ApplyFunc(
						access_key.Create,
						func(_ context.Context, _ *access_key.CreateReq) (*access_key.CreateRsp, error) {
							return &access_key.CreateRsp{
								Name:         "test",
								IdentityType: enums.ACCESS_KEY_IDENTITY_TYPE__PUBLISHER,
								IdentityID:   100,
								AccessKey:    "w3b_any",
							}, nil
						},
					)
				t.Run("#PublisherConflict", func(t *testing.T) {
					db.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrConflict).Times(1)
					_, err := publisher.Create(ctx, &publisher.CreateReq{})
					mock_sqlx.ExpectError(t, err, status.PublisherConflict)
				})
				t.Run("#DatabaseError", func(t *testing.T) {
					db.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrDatabase).Times(1)
					_, err := publisher.Create(ctx, &publisher.CreateReq{})
					mock_sqlx.ExpectError(t, err, status.DatabaseError)
				})
			})
		})

		t.Run("#Success", func(t *testing.T) {
			db.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)
			_, err := publisher.Create(ctx, &publisher.CreateReq{})
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("#Success", func(t *testing.T) {
			db.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(2)

			// GetBySFID
			{
				_, err := publisher.GetBySFID(ctx, idg.MustGenSFID())
				NewWithT(t).Expect(err).To(BeNil())
			}

			// GetByProjectAndKey
			{
				_, err := publisher.GetByProjectAndKey(ctx, idg.MustGenSFID(), "publisher_key")
				NewWithT(t).Expect(err).To(BeNil())
			}
		})

		t.Run("#PublisherNotFound", func(t *testing.T) {
			db.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(sqlx.NewSqlError(sqlx.SqlErrTypeNotFound, "")).Times(2)

			// GetBySFID
			{
				_, err := publisher.GetBySFID(ctx, 1)
				NewWithT(t).Expect(err).To(Equal(status.PublisherNotFound))
			}

			// GetByProjectAndKey
			{
				_, err := publisher.GetByProjectAndKey(ctx, idg.MustGenSFID(), "publisher_key")
				NewWithT(t).Expect(err).To(Equal(status.PublisherNotFound))
			}
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			db.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(5)

			// List
			{
				_, err := publisher.List(ctx, &publisher.ListReq{})
				NewWithT(t).Expect(err).To(BeNil())
			}

			// ListByCond
			{
				_, err := publisher.ListByCond(ctx, &publisher.CondArgs{})
				NewWithT(t).Expect(err).To(BeNil())
			}

			{
				_, err := publisher.ListDetail(ctx, &publisher.ListReq{})
				NewWithT(t).Expect(err).To(BeNil())
			}
		})
	})

	t.Run("Remove", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			db.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(driver.RowsAffected(1), nil)

			err := publisher.Remove(ctx, types.MustAccountFromContext(ctx), &publisher.CondArgs{})
			NewWithT(t).Expect(err).To(BeNil())
		})
	})
}
