package middleware_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	confjwt "github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	mock_sqlx "github.com/machinefi/w3bstream/pkg/test/mock_depends_kit_sqlx"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestContextAccountAuth(t *testing.T) {
	t.Run("ContextKey", func(t *testing.T) {
		caa := &middleware.ContextAccountAuth{}
		NewWithT(t).Expect(caa.ContextKey()).To(Equal("middleware.ContextAccountAuth"))
	})

	t.Run("Output", func(t *testing.T) {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		conf := &confjwt.Jwt{
			Issuer:  "test_context_account_auth",
			ExpIn:   base.Duration(time.Minute),
			SignKey: "__test__",
		}

		d := mock_sqlx.NewMockDBExecutor(ctl)
		d.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
		d.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrDatabase).Times(1)
		d.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		ctx := contextx.WithContextCompose(
			confjwt.WithConfContext(conf),
			types.WithMgrDBExecutorContext(d),
		)(context.Background())
		caa := &middleware.ContextAccountAuth{}

		var cases = []*struct {
			name      string
			authValue interface{}
			expectErr status.Error
		}{
			{
				name:      "#Failed##Bytes",
				authValue: []byte("any"),
				expectErr: status.InvalidAuthAccountID,
			},
			{
				name:      "#Failed##String",
				authValue: "any",
				expectErr: status.InvalidAuthAccountID,
			},
			{
				name:      "#Failed##Stringer",
				authValue: time.Now(),
				expectErr: status.InvalidAuthAccountID,
			},
			{
				name:      "#Failed#InvalidAuthValue",
				authValue: 1,
				expectErr: status.InvalidAuthValue,
			},
			{
				name:      "#Failed#QueryFailed",
				authValue: types.SFID(100),
				expectErr: status.DatabaseError,
			},
			{
				name:      "#Success",
				authValue: types.SFID(100),
				expectErr: 0,
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				confjwt.SetBuiltInTokenFn(
					func(ctx context.Context, s string) (interface{}, error, bool) {
						return c.authValue, nil, true
					},
				)

				ctx := contextx.WithValue(ctx, (&confjwt.Auth{}).ContextKey(), c.authValue)

				_, err := caa.Output(ctx)
				if c.expectErr != 0 {
					mock_sqlx.ExpectError(t, err, c.expectErr)
				} else {
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
		}
	})
}
