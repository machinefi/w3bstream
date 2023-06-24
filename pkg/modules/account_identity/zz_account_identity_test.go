package account_identity_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/account_identity"
	mock_sqlx "github.com/machinefi/w3bstream/pkg/test/mock_depends_kit_sqlx"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestAccountIdentity(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	d := mock_sqlx.NewMockDBExecutor(c)
	idg := confid.MustNewSFIDGenerator()
	ctx := contextx.WithContextCompose(
		types.WithMgrDBExecutorContext(d),
		confid.WithSFIDGeneratorContext(idg),
	)(context.Background())

	t.Run("GetBySFIDAndType", func(t *testing.T) {
		d.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()

		id, typ := idg.MustGenSFID(), enums.ACCOUNT_IDENTITY_TYPE__USERNAME

		f := account_identity.GetBySFIDAndType

		t.Run("#Success", func(t *testing.T) {
			d.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil)

			_, err := f(ctx, id, typ)
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#AccountIdentityNotFound", func(t *testing.T) {
				d.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound)

				_, err := f(ctx, id, typ)
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.AccountIdentityNotFound.Key()))
			})
			t.Run("#DatabaseError", func(t *testing.T) {
				d.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrDatabase)

				_, err := f(ctx, id, typ)
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.DatabaseError.Key()))
			})
		})
	})
}
