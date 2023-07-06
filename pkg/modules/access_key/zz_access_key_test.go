package access_key_test

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
	"github.com/pkg/errors"

	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
	mock_sqlx "github.com/machinefi/w3bstream/pkg/test/mock_depends_kit_sqlx"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestAccessKeyContext(t *testing.T) {
	kctx := access_key.NewDefaultAccessKeyContext()
	target := &access_key.AccessKeyContext{}

	key, err := kctx.MarshalText()
	NewWithT(t).Expect(err).To(BeNil())
	t.Logf(string(key))

	err = target.UnmarshalText(key)
	NewWithT(t).Expect(err).To(BeNil())

	NewWithT(t).Expect(kctx.Equal(target)).To(BeTrue())

	cases := []*struct {
		name string
		key  string
		err  error
	}{
		{
			name: "#Failed#InvalidPartCountOrPrefix",
			key:  "invalid_key",
			err:  access_key.ErrInvalidPrefixOrPartCount,
		},
		{
			name: "#Failed#Base64DecodeFailed",
			key:  "w3b_YWJjZA=====",
			err:  access_key.ErrBase64DecodeFailed,
		},
		{
			name: "#Failed#InvalidContentsPartCount",
			key:  "w3b_" + base64.RawURLEncoding.EncodeToString([]byte("1xxx")),
			err:  access_key.ErrInvalidContentsPartCount,
		},
		{
			name: "#Failed#ParseVersionFailed",
			key:  "w3b_" + base64.RawURLEncoding.EncodeToString([]byte("b_100_xxx")),
			err:  access_key.ErrParseVersionFailed,
		},
		{
			name: "#Failed#ParseGenTsFailed",
			key:  "w3b_" + base64.RawURLEncoding.EncodeToString([]byte("1_aaa_xxx")),
			err:  access_key.ErrParseGenTsFailed,
		},
		{
			name: "#Failed#InvalidVersion",
			key:  "w3b_" + base64.RawURLEncoding.EncodeToString([]byte("2_100_xxx")),
			err:  access_key.ErrInvalidVersion,
		},
		{
			name: "#Success",
			key:  "w3b_" + base64.StdEncoding.EncodeToString([]byte("1_100_xxx")),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			k := &access_key.AccessKeyContext{}
			err = k.UnmarshalText([]byte(c.key))
			if c.err != nil {
				NewWithT(t).Expect(strings.Contains(err.Error(), c.err.Error())).To(BeTrue())
			} else {
				NewWithT(t).Expect(err).To(BeNil())
			}
		})
	}
}

func TestCondArgs(t *testing.T) {
	args := &access_key.CondArgs{}
	t.Log(builder.ResolveExpr(args.Condition()).Query())
	args = &access_key.CondArgs{
		AccountID:      100,
		Names:          []string{"test"},
		ExpiredAtBegin: types.Timestamp{Time: time.Now().Add(-time.Hour)},
		ExpiredAtEnd:   types.Timestamp{Time: time.Now().Add(time.Hour)},
		IdentityIDs:    types.SFIDs{1, 2, 3},
		IdentityTypes:  []enums.AccessKeyIdentityType{1, 2},
	}
	t.Log(builder.ResolveExpr(args.Condition()).Query())
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

func TestAccessKey(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	d := &struct {
		*mock_sqlx.MockDBExecutor
		*mock_sqlx.MockTxExecutor
	}{
		MockDBExecutor: mock_sqlx.NewMockDBExecutor(ctl),
		MockTxExecutor: mock_sqlx.NewMockTxExecutor(ctl),
	}
	idg := confid.MustNewSFIDGenerator()
	acc := &models.Account{
		RelAccount: models.RelAccount{AccountID: idg.MustGenSFID()},
	}

	ctx := contextx.WithContextCompose(
		types.WithMgrDBExecutorContext(d),
		confid.WithSFIDGeneratorContext(idg),
		types.WithAccountContext(acc),
	)(context.Background())

	d.MockDBExecutor.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
	d.MockTxExecutor.EXPECT().IsTx().Return(true).AnyTimes()
	d.MockDBExecutor.EXPECT().Context().Return(ctx).AnyTimes()

	anyErr := errors.New("any error")

	t.Run("Create", func(t *testing.T) {
		pub := &models.Publisher{
			RelPublisher: models.RelPublisher{PublisherID: idg.MustGenSFID()},
		}

		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Success", func(t *testing.T) {
			cases := []*struct {
				name string
				req  *access_key.CreateReq
			}{
				{
					name: "#NoExpiration",
					req: &access_key.CreateReq{
						CreateReqBase: access_key.CreateReqBase{
							Name:           "test",
							ExpirationDays: 0,
						},
					},
				},
				{
					name: "#HasExpiration",
					req: &access_key.CreateReq{
						CreateReqBase: access_key.CreateReqBase{
							Name:           "test",
							ExpirationDays: 30,
						},
					},
				},
			}
			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).Times(1)
					d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)
					_, err := access_key.Create(ctx, c.req)
					NewWithT(t).Expect(err).To(BeNil())
				})
			}
			t.Run("CreateForPublisher", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).Times(1)
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)
				_, err := access_key.Create(ctx, &access_key.CreateReq{
					IdentityID:   pub.PublisherID,
					IdentityType: enums.ACCESS_KEY_IDENTITY_TYPE__PUBLISHER,
				})
				NewWithT(t).Expect(err).To(BeNil())
			})
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#InvalidAccessKeyIdentityType", func(t *testing.T) {
				_, err := access_key.Create(ctx, &access_key.CreateReq{
					IdentityID:   0,
					IdentityType: enums.AccessKeyIdentityType(100),
				})
				mock_sqlx.ExpectError(t, err, status.InvalidAccessKeyIdentityType)
			})
			t.Run("#FetchByRandDatbaseError", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrDatabase).Times(1)
				_, err := access_key.Create(ctx, &access_key.CreateReq{})
				mock_sqlx.ExpectError(t, err, status.DatabaseError)
			})
			t.Run("#CreateAccessKeyConflict", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).Times(1)
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrConflict).Times(1)
				_, err := access_key.Create(ctx, &access_key.CreateReq{})
				mock_sqlx.ExpectError(t, err, status.AccessKeyNameConflict)
			})
			t.Run("#CreateAccessKeyNameDatabaseError", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).Times(1)
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrDatabase).Times(1)
				_, err := access_key.Create(ctx, &access_key.CreateReq{})
				mock_sqlx.ExpectError(t, err, status.DatabaseError)
			})
		})
	})

	t.Run("DeleteByName", func(t *testing.T) {
		t.Run("#Success", func(t *testing.T) {
			d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).MaxTimes(1)

			err := access_key.DeleteByName(ctx, "any_name")
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#AccountKeyNotFound", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrNotFound).MaxTimes(1)

				err := access_key.DeleteByName(ctx, "any")
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.AccessKeyNotFound.Key()))
			})
			t.Run("#DatabaseError", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrDatabase).MaxTimes(1)

				err := access_key.DeleteByName(ctx, "any")
				NewWithT(t).Expect(err).NotTo(BeNil())

				se, ok := statusx.IsStatusErr(err)
				NewWithT(t).Expect(ok).To(BeTrue())
				NewWithT(t).Expect(se.Key).To(Equal(status.DatabaseError.Key()))
			})
		})
	})

	t.Run("Validate", func(t *testing.T) {
		kctx := access_key.NewAccessKeyContext(1)

		id := idg.MustGenSFID()
		m := &models.AccessKey{
			RelAccount: models.RelAccount{AccountID: id},
			AccessKeyInfo: models.AccessKeyInfo{
				IdentityID:   id,
				IdentityType: enums.ACCESS_KEY_IDENTITY_TYPE__ACCOUNT,
				Name:         "test",
				Rand:         kctx.Rand,
				ExpiredAt:    base.Timestamp{Time: kctx.GenTS.UTC().Add(2 * time.Second)},
			},
			OperationTimesWithDeleted: datatypes.OperationTimesWithDeleted{
				OperationTimes: datatypes.OperationTimes{
					CreatedAt: base.Timestamp{Time: kctx.GenTS.UTC()},
				},
			},
		}
		key, _ := kctx.MarshalText()

		t.Run("#Success", func(t *testing.T) {
			if runtime.GOOS == `darwin` {
				return
			}
			patch := gomonkey.ApplyMethod(
				reflect.TypeOf(&models.AccessKey{}),
				"FetchByRand",
				func(receiver *models.AccessKey, d sqlx.DBExecutor) error {
					*receiver = *m
					return nil
				},
			)
			defer patch.Reset()
			d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)

			idAny, err, canBeValidated := access_key.Validate(ctx, string(key))

			NewWithT(t).Expect(canBeValidated).To(BeTrue())
			NewWithT(t).Expect(err).To(BeNil())

			idVal, ok := idAny.(types.SFID)
			NewWithT(t).Expect(ok).To(BeTrue())
			NewWithT(t).Expect(idVal).To(Equal(id))
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#AccessKeyContextUnmarshalFailed", func(t *testing.T) {
				t.Run("#CanbeValidated", func(t *testing.T) {
					_, err, canbe := access_key.Validate(ctx, "invalid_key")
					NewWithT(t).Expect(err).NotTo(BeNil())
					NewWithT(t).Expect(canbe).To(BeFalse())
				})
				t.Run("#CannotBeValidated", func(t *testing.T) {
					key := "w3b_xxxx"
					_, err, canbe := access_key.Validate(ctx, key)
					NewWithT(t).Expect(err).NotTo(BeNil())
					NewWithT(t).Expect(canbe).To(BeTrue())
				})
			})
			t.Run("#FetchByRandFailed", func(t *testing.T) {
				t.Run("#AccessKeyNotFound", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).Times(1)
					_, err, _ := access_key.Validate(ctx, string(key))
					mock_sqlx.ExpectError(t, err, status.AccessKeyNotFound)
				})
				t.Run("#DatabaseError", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrDatabase).Times(1)
					_, err, _ := access_key.Validate(ctx, string(key))
					mock_sqlx.ExpectError(t, err, status.DatabaseError)
				})
			})
			t.Run("#UpdateLastUsedFailed", func(t *testing.T) {
				if runtime.GOOS == `darwin` {
					return
				}
				patch := gomonkey.ApplyMethod(
					reflect.TypeOf(&models.AccessKey{}),
					"FetchByRand",
					func(receiver *models.AccessKey, d sqlx.DBExecutor) error {
						*receiver = *m
						return nil
					},
				)
				defer patch.Reset()
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, anyErr).Times(1)
				_, err, _ := access_key.Validate(ctx, string(key))
				NewWithT(t).Expect(err).To(BeNil())
			})
			t.Run("#GenTsNotMatch", func(t *testing.T) {
				if runtime.GOOS == `darwin` {
					return
				}
				patch := gomonkey.ApplyMethod(
					reflect.TypeOf(&models.AccessKey{}),
					"FetchByRand",
					func(receiver *models.AccessKey, d sqlx.DBExecutor) error {
						*receiver = *m
						receiver.CreatedAt = base.Timestamp{
							Time: receiver.CreatedAt.Add(time.Second),
						}
						return nil
					},
				)
				defer patch.Reset()
				_, err, _ := access_key.Validate(ctx, string(key))
				mock_sqlx.ExpectError(t, err, status.InvalidAccessKey)
			})
			t.Run("#AccessKeyExpired", func(t *testing.T) {
				if runtime.GOOS == `darwin` {
					return
				}
				time.Sleep(3 * time.Second)
				patch := gomonkey.ApplyMethod(
					reflect.TypeOf(&models.AccessKey{}),
					"FetchByRand",
					func(receiver *models.AccessKey, d sqlx.DBExecutor) error {
						*receiver = *m
						return nil
					},
				)
				defer patch.Reset()
				_, err, _ := access_key.Validate(ctx, string(key))
				mock_sqlx.ExpectError(t, err, status.AccessKeyExpired)
			})
		})
	})

	t.Run("List", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.
			ApplyMethod(
				reflect.TypeOf(&models.AccessKey{}),
				"List",
				func(r *models.AccessKey, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.AccessKey, error) {
					return []models.AccessKey{
						{},
						{
							AccessKeyInfo: models.AccessKeyInfo{
								ExpiredAt: base.Timestamp{Time: time.Now()},
								LastUsed:  base.Timestamp{Time: time.Now()},
							},
						},
					}, nil
				},
			).
			ApplyMethod(
				reflect.TypeOf(&models.AccessKey{}),
				"Count",
				func(r *models.AccessKey, _ sqlx.DBExecutor, _ builder.SqlCondition) (int64, error) {
					return 2, nil
				},
			)
		defer patch.Reset()
		t.Run("#Success", func(t *testing.T) {
			rsp, err := access_key.List(ctx, &access_key.ListReq{})
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(rsp.Total).To(Equal(int64(2)))
			NewWithT(t).Expect(len(rsp.Data)).To(Equal(2))
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#ListFailed", func(t *testing.T) {
				patch = patch.
					ApplyMethod(
						reflect.TypeOf(&models.AccessKey{}),
						"List",
						func(r *models.AccessKey, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.AccessKey, error) {
							return nil, anyErr
						},
					)
				_, err := access_key.List(ctx, &access_key.ListReq{})
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
			t.Run("#CountFailed", func(t *testing.T) {
				patch = patch.
					ApplyMethod(
						reflect.TypeOf(&models.AccessKey{}),
						"List",
						func(r *models.AccessKey, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.AccessKey, error) {
							return []models.AccessKey{
								{},
								{
									AccessKeyInfo: models.AccessKeyInfo{
										ExpiredAt: base.Timestamp{Time: time.Now()},
										LastUsed:  base.Timestamp{Time: time.Now()},
									},
								},
							}, nil
						},
					).
					ApplyMethod(
						reflect.TypeOf(&models.AccessKey{}),
						"Count",
						func(r *models.AccessKey, _ sqlx.DBExecutor, _ builder.SqlCondition) (int64, error) {
							return 0, anyErr
						},
					)
				_, err := access_key.List(ctx, &access_key.ListReq{})
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
		})
	})
}
