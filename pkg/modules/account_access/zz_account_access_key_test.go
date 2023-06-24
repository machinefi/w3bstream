package account_access_test

import (
	"context"
	"encoding/base64"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/account_access"
	mock_sqlx "github.com/machinefi/w3bstream/pkg/test/mock_depends_kit_sqlx"
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

	cases := []*struct {
		name string
		key  string
		err  error
	}{
		{
			name: "#InvalidPartedCountOrPrefix",
			key:  "w3b_key_fmt_error",
			err:  account_access.ErrMsgAccessKeyInvalidPartCountOrPrefix,
		},
		{
			name: "#Base64DecodeFailed",
			key:  "w3b_YWJjZA=====",
			err:  account_access.ErrMsgAccessKeyBase64Decode,
		},
		{
			name: "#InvalidPartedCount",
			key:  "w3b_" + base64.StdEncoding.EncodeToString([]byte("1_xxx")),
			err:  account_access.ErrMsgAccessKeyInvalidPartCount,
		},
		{
			name: "#InvalidAccountID",
			key:  "w3b_" + base64.StdEncoding.EncodeToString([]byte("a_xxx_xxx")),
			err:  account_access.ErrMsgAccessKeyInvalidAccountID,
		},
		{
			name: "#InvalidTimestamp",
			key:  "w3b_" + base64.StdEncoding.EncodeToString([]byte("100_xxx_xxx")),
			err:  account_access.ErrMsgAccessKeyInvalidTimestamp,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, _, _, err := account_access.ParseAccessKey(c.key)
			NewWithT(t).Expect(err).NotTo(BeNil())

			se, ok := statusx.IsStatusErr(err)
			NewWithT(t).Expect(ok).To(BeTrue())
			NewWithT(t).Expect(se.Key).To(Equal(status.InvalidAccountAccessKey.Key()))
			NewWithT(t).Expect(strings.Contains(se.Desc, c.err.Error())).To(BeTrue())
		})
	}

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

	d.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()

	t.Run("Create", func(t *testing.T) {
		t.Run("#Success", func(t *testing.T) {
			t.Run("#NoExpiration", func(t *testing.T) {
				d.EXPECT().Exec(gomock.Any()).Return(nil, nil).MaxTimes(1)

				_, err := account_access.Create(ctx, &account_access.CreateReq{})
				NewWithT(t).Expect(err).To(BeNil())
			})
			t.Run("#HasExpiration", func(t *testing.T) {
				d.EXPECT().Exec(gomock.Any()).Return(nil, nil).MaxTimes(1)

				_, err := account_access.Create(ctx, &account_access.CreateReq{
					ExpirationDays: 1,
				})
				NewWithT(t).Expect(err).To(BeNil())
			})
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#AccountKeyNameConflict", func(t *testing.T) {
				d.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrConflict).MaxTimes(1)

				_, err := account_access.Create(ctx, &account_access.CreateReq{})
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.AccountKeyNameConflict.Key()))
			})
			t.Run("#DatabaseError", func(t *testing.T) {
				d.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrDatabase).MaxTimes(1)

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
			d.EXPECT().Exec(gomock.Any()).Return(nil, nil).MaxTimes(1)

			err := account_access.DeleteByName(ctx, "any_name")
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#AccountKeyNotFound", func(t *testing.T) {
				d.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrNotFound).MaxTimes(1)

				err := account_access.DeleteByName(ctx, "any")
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.AccountKeyNotFound.Key()))
			})
			t.Run("#DatabaseError", func(t *testing.T) {
				d.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrDatabase).MaxTimes(1)

				err := account_access.DeleteByName(ctx, "any")
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.DatabaseError.Key()))
			})
		})
	})

	t.Run("Validate", func(t *testing.T) {
		id := idg.MustGenSFID()
		rand, key, ts := account_access.GenAccessKey(id)
		t.Logf("rand: %s", rand)
		t.Logf("key: %s", key)
		t.Logf("ts: %v", ts)

		m := &models.AccountAccessKey{
			RelAccount: models.RelAccount{AccountID: id},
			AccountAccessKeyInfo: models.AccountAccessKeyInfo{
				Name:      "test_gen_key",
				AccessKey: rand,
				ExpiredAt: base.Timestamp{Time: ts.Add(time.Hour).UTC()},
			},
			OperationTimesWithDeleted: datatypes.OperationTimesWithDeleted{
				OperationTimes: datatypes.OperationTimes{
					CreatedAt: base.Timestamp{Time: ts},
				},
			},
		}

		t.Run("#Success", func(t *testing.T) {
			if runtime.GOOS == `darwin` {
				return
			}
			patches := gomonkey.ApplyMethod(
				reflect.TypeOf(&models.AccountAccessKey{}),
				"FetchByAccessKey",
				func(receiver *models.AccountAccessKey, d sqlx.DBExecutor) error {
					*receiver = *m
					return nil
				},
			)
			defer patches.Reset()

			idAny, err, canBeValidated := account_access.Validate(ctx, key)

			NewWithT(t).Expect(canBeValidated).To(BeTrue())
			t.Log(err)
			NewWithT(t).Expect(err).To(BeNil())

			idVal, ok := idAny.(types.SFID)
			NewWithT(t).Expect(ok).To(BeTrue())
			NewWithT(t).Expect(idVal).To(Equal(id))
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#NotW3bValidateKey", func(t *testing.T) {
				_, _, canBeValidated := account_access.Validate(ctx, "not_w3b_key")
				NewWithT(t).Expect(canBeValidated).To(BeFalse())
			})
			t.Run("#ParseAccessKeyFailed", func(t *testing.T) {
				_, err, _ := account_access.Validate(ctx, "w3b_key_fmt_error")
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
			t.Run("#AccountKeyNotFound", func(t *testing.T) {
				d.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).MaxTimes(1)

				_, err, _ := account_access.Validate(ctx, key)
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.AccountKeyNotFound.Key()))
			})
			t.Run("#DatabaseError", func(t *testing.T) {
				d.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrDatabase).MaxTimes(1)

				_, err, _ := account_access.Validate(ctx, key)
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.DatabaseError.Key()))
			})
			t.Run("#AccountIDNotMatch", func(t *testing.T) {
				if runtime.GOOS == `darwin` {
					return
				}

				m2 := *m
				m2.AccountID = 0
				patches := gomonkey.ApplyMethod(
					reflect.TypeOf(&models.AccountAccessKey{}),
					"FetchByAccessKey",
					func(receiver *models.AccountAccessKey, d sqlx.DBExecutor) error {
						*receiver = m2
						return nil
					},
				)
				defer patches.Reset()

				_, err, _ := account_access.Validate(ctx, key)
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.InvalidAccountAccessKey.Key()))
			})
			t.Run("#AccountAccessKeyExpired", func(t *testing.T) {
				if runtime.GOOS == `darwin` {
					return
				}

				m2 := *m
				m2.ExpiredAt = types.Timestamp{Time: time.Now().UTC().Add(-time.Hour)}
				patches := gomonkey.ApplyMethod(
					reflect.TypeOf(&models.AccountAccessKey{}),
					"FetchByAccessKey",
					func(receiver *models.AccountAccessKey, d sqlx.DBExecutor) error {
						*receiver = m2
						return nil
					},
				)
				defer patches.Reset()

				_, err, _ := account_access.Validate(ctx, key)
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.AccountAccessKeyExpired.Key()))
			})
		})
	})
}
