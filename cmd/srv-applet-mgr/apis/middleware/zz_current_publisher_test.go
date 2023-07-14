package middleware_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	confjwt "github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	mock_sqlx "github.com/machinefi/w3bstream/pkg/test/mock_depends_kit_sqlx"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestCurrentPublisher(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	pub := &middleware.ContextPublisherAuth{}

	t.Run("#ContextKey", func(t *testing.T) {
		NewWithT(t).Expect(pub.ContextKey()).To(Equal("middleware.ContextAccountAuth"))
	})

	t.Run("#Output", func(t *testing.T) {
		conf := &confjwt.Jwt{
			Issuer:  "test_context_account_auth",
			ExpIn:   base.Duration(time.Minute),
			SignKey: "__test__",
		}

		d := mock_sqlx.NewMockDBExecutor(ctl)
		d.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()

		ctx := contextx.WithContextCompose(
			confjwt.WithConfContext(conf),
			types.WithMgrDBExecutorContext(d),
		)(context.Background())
		key := (&confjwt.Auth{}).ContextKey()
		errFrom := func(from string) error { return errors.New(from) }

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#ParseJwtAuthContentFailed", func(t *testing.T) {
				ctx := context.WithValue(ctx, key, "wrong_auth_content")
				_, err := pub.Output(ctx)
				mock_sqlx.ExpectError(t, err, status.InvalidAuthPublisherID)
				t.Log(err)
			})
			t.Run("#AccountModelQueryFailed", func(t *testing.T) {
				d.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(errFrom(t.Name())).Times(1)
				ctx := context.WithValue(ctx, key, &models.AccessKey{
					AccessKeyInfo: models.AccessKeyInfo{
						IdentityID:   100,
						IdentityType: enums.ACCESS_KEY_IDENTITY_TYPE__ACCOUNT,
					},
				})
				_, err := pub.Output(ctx)
				mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
				t.Log(err)
			})
			t.Run("#PublisherModelQueryFailed", func(t *testing.T) {
				d.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(errFrom(t.Name())).Times(1)
				ctx := context.WithValue(ctx, key, &models.AccessKey{
					AccessKeyInfo: models.AccessKeyInfo{
						IdentityID:   100,
						IdentityType: enums.ACCESS_KEY_IDENTITY_TYPE__PUBLISHER,
					},
				})
				_, err := pub.Output(ctx)
				mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
				t.Log(err)
			})
		})
		t.Run("#Success", func(t *testing.T) {
			t.Run("#AsPublisher", func(t *testing.T) {
				d.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				ctx := context.WithValue(ctx, key, "100")
				v, err := pub.Output(ctx)
				NewWithT(t).Expect(err).To(BeNil())

				_, ok := v.(*middleware.CurrentPublisher)
				NewWithT(t).Expect(ok).To(BeTrue())
			})
			t.Run("#AsAccount", func(t *testing.T) {
				d.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				ctx = context.WithValue(ctx, key, &models.AccessKey{
					AccessKeyInfo: models.AccessKeyInfo{
						IdentityID:   100,
						IdentityType: enums.ACCESS_KEY_IDENTITY_TYPE__ACCOUNT,
					},
				})
				v, err := pub.Output(ctx)
				NewWithT(t).Expect(err).To(BeNil())

				_, ok := v.(*middleware.CurrentAccount)
				NewWithT(t).Expect(ok).To(BeTrue())
			})
		})
	})
}

func TestMustPublisher(t *testing.T) {
	key := (&middleware.ContextPublisherAuth{}).ContextKey()
	ctx := context.WithValue(context.Background(), key, &middleware.CurrentPublisher{})

	_ = middleware.MustPublisher(ctx)
	_, ok := middleware.PublisherFromContext(ctx)
	NewWithT(t).Expect(ok).To(BeTrue())
}
