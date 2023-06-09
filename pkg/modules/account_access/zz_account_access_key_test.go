package account_access_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/account_access"
	"github.com/machinefi/w3bstream/pkg/test/mock/mock_sqlx"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestGenAndParseAccessKey(t *testing.T) {
	id := confid.MustNewSFIDGenerator().MustGenSFID()

	rand, key, ts := account_access.GenAccessKey(id)

	_id, _rand, _ts, err := account_access.ParseAccessKey(key)

	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(rand).To(Equal(_rand))
	NewWithT(t).Expect(id).To(Equal(_id))
	NewWithT(t).Expect(ts.Equal(_ts)).To(BeTrue())

	t.Logf(key)
}

func TestTimeParseAndFormat(t *testing.T) {
	formatAndParse := func(layout string) (error, bool) {
		ts := time.Now().UTC()
		_ts, err := time.ParseInLocation(layout, ts.Format(layout), time.UTC)
		return err, ts.Equal(_ts)
	}

	for _, layout := range []string{time.RFC3339, time.RFC3339Nano} {
		err, equal := formatAndParse(layout)
		t.Logf("layout: %s err: %v, equal: %v", layout, err, equal)
	}
}

func TestAccountAccessKey(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	d := mock_sqlx.NewMockDBExecutor(ctl)
	idg := confid.MustNewSFIDGenerator()
	acc := &models.Account{
		RelAccount: models.RelAccount{AccountID: idg.MustGenSFID()},
	}

	ctx := contextx.WithContextCompose(
		types.WithMgrDBExecutorContext(d),
		confid.WithSFIDGeneratorContext(idg),
		types.WithAccountContext(acc),
	)(context.Background())

	t.Run("Create", func(t *testing.T) {
		t.Run("#Success", func(t *testing.T) {
			d.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			d.EXPECT().Exec(gomock.Any()).Return(nil, nil)

			_, err := account_access.Create(ctx, &account_access.CreateReq{})
			NewWithT(t).Expect(err).To(BeNil())
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#AccountKeyNameConflict", func(t *testing.T) {
				d.EXPECT().T(gomock.Any()).Return(&builder.Table{})
				d.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrConflict)

				_, err := account_access.Create(ctx, &account_access.CreateReq{})
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.AccountKeyNameConflict.Key()))
			})
			t.Run("#DatabaseError", func(t *testing.T) {
				d.EXPECT().T(gomock.Any()).Return(&builder.Table{})
				d.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrDatabase)

				_, err := account_access.Create(ctx, &account_access.CreateReq{})
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.DatabaseError.Key()))
			})
		})
	})

	t.Run("DeleteByName", func(t *testing.T) {
		t.Run("#Success", func(t *testing.T) {
			d.EXPECT().T(gomock.Any()).Return(&builder.Table{})
			d.EXPECT().Exec(gomock.Any()).Return(nil, nil)

			err := account_access.DeleteByName(ctx, "any_name")
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#AccountKeyNotFound", func(t *testing.T) {
				d.EXPECT().T(gomock.Any()).Return(&builder.Table{})
				d.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrNotFound)

				err := account_access.DeleteByName(ctx, "any")
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.AccountKeyNotFound.Key()))
			})
			t.Run("#DatabaseError", func(t *testing.T) {
				d.EXPECT().T(gomock.Any()).Return(&builder.Table{})
				d.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrDatabase)

				err := account_access.DeleteByName(ctx, "any")
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.DatabaseError.Key()))
			})
		})
	})
}
