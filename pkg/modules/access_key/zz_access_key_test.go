package access_key_test

import (
	"context"
	"database/sql"
	"encoding/base64"
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
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
	mock_sqlx "github.com/machinefi/w3bstream/pkg/test/mock_depends_kit_sqlx"
	"github.com/machinefi/w3bstream/pkg/test/patch_models"
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

var (
	RootRouter1 *kit.Router
	RootRouter2 *kit.Router
	RootRouter3 *kit.Router
)

type MockGet struct{ httpx.MethodGet }

func (*MockGet) Output(_ context.Context) (interface{}, error) { return nil, nil }

type MockPos struct{ httpx.MethodPost }

func (*MockPos) Output(_ context.Context) (interface{}, error) { return nil, nil }

func (*MockPos) OperatorAttr() enums.ApiOperatorAttr { return enums.API_OPERATOR_ATTR__COMMON }

func init() {
	RootRouter1 = kit.NewRouter(httptransport.Group("mock1"))
	RootRouter2 = kit.NewRouter(httptransport.Group("mock2"))
	RootRouter3 = kit.NewRouter(httptransport.Group("mock3"))

	RootRouter1.Register(kit.NewRouter(&MockGet{}))
	RootRouter1.Register(kit.NewRouter(&MockPos{}))
	RootRouter2.Register(kit.NewRouter(&MockGet{}))
	RootRouter2.Register(kit.NewRouter(&MockPos{}))
	RootRouter3.Register(kit.NewRouter(&MockGet{}))
	RootRouter3.Register(kit.NewRouter(&MockPos{}))

	access_key.RouterRegister(RootRouter1, "mock1", "mock group operator")
	access_key.RouterRegister(RootRouter2, "mock2", "mock group operator")
	access_key.RouterRegister(RootRouter3, "mock3", "mock group operator")
}

func TestRouterRegister(t *testing.T) {
	defer func() {
		err := recover()
		NewWithT(t).Expect(strings.Contains(err.(error).Error(), "already registered"))
	}()
	access_key.RouterRegister(RootRouter1, "mock1", "mock group operator")
}

func TestOperatorGroupMetaList(t *testing.T) {
	metas := access_key.OperatorGroupMetaList()
	NewWithT(t).Expect(len(metas)).To(Equal(3))
	NewWithT(t).Expect(metas[0].Name).To(Equal("mock1"))
	NewWithT(t).Expect(metas[1].Name).To(Equal("mock2"))
	NewWithT(t).Expect(metas[2].Name).To(Equal("mock3"))
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

	errFrom := func(from string) error { return errors.New(from) }

	t.Run("Create", func(t *testing.T) {
		pub := &models.Publisher{
			RelPublisher: models.RelPublisher{PublisherID: idg.MustGenSFID()},
		}

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
							Privileges: access_key.GroupAccessPrivileges{
								{Name: "mock1", Perm: enums.ACCESS_PERMISSION__NO_ACCESS},
								{Name: "mock2", Perm: enums.ACCESS_PERMISSION__READONLY},
								{Name: "mock3", Perm: enums.ACCESS_PERMISSION__READ_WRITE},
								{Name: "not_exists"},
							},
						},
					},
				},
				{
					name: "#HasExpiration",
					req: &access_key.CreateReq{
						CreateReqBase: access_key.CreateReqBase{
							Name:           "test",
							ExpirationDays: 30,
							Privileges: access_key.GroupAccessPrivileges{
								{Name: "mock1", Perm: enums.ACCESS_PERMISSION__NO_ACCESS},
								{Name: "mock2", Perm: enums.ACCESS_PERMISSION__READONLY},
								{Name: "mock3", Perm: enums.ACCESS_PERMISSION__READ_WRITE},
							},
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

	t.Run("#Update", func(t *testing.T) {
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#FetchByAccountIDAndNameFailed", func(t *testing.T) {
				t.Run("#AccessKeyNotFound", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).Times(1)
					err := access_key.UpdateByName(ctx, "any_name", &access_key.UpdateReq{})
					mock_sqlx.ExpectError(t, err, status.AccessKeyNotFound)
				})
				t.Run("#DatabaseError", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(errFrom(t.Name())).Times(1)
					err := access_key.UpdateByName(ctx, "any_name", &access_key.UpdateReq{})
					mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
				})
			})
			req := &access_key.UpdateReq{
				ExpirationDays: 10,
				Desc:           "desc",
			}
			if runtime.GOOS == `darwin` {
				return
			}
			patch = patch_models.AccessKeyFetchByAccountIDAndName(patch, &models.AccessKey{}, nil)
			t.Run("#UpdateFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(sql.Result(nil), errFrom(t.Name())).Times(1)
				err := access_key.UpdateByName(ctx, "any_name", req)
				mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
			})
		})
		t.Run("#Success", func(t *testing.T) {
			if runtime.GOOS == `darwin` {
				return
			}
			d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(sql.Result(nil), nil).Times(1)
			err := access_key.UpdateByName(ctx, "any_name", &access_key.UpdateReq{})
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("#DeleteByName", func(t *testing.T) {
		t.Run("#Success", func(t *testing.T) {
			d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).MaxTimes(1)

			err := access_key.DeleteByName(ctx, "any_name")
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#AccountKeyNotFound", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrNotFound).MaxTimes(1)

				err := access_key.DeleteByName(ctx, "any")
				mock_sqlx.ExpectError(t, err, status.AccessKeyNotFound)
			})
			t.Run("#DatabaseError", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, errFrom(t.Name())).MaxTimes(1)

				err := access_key.DeleteByName(ctx, "any")
				mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
			})
		})
	})

	t.Run("#Validate", func(t *testing.T) {
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
				Privileges:   models.AccessPrivileges{"MockOperatorID": {}},
			},
			OperationTimesWithDeleted: datatypes.OperationTimesWithDeleted{
				OperationTimes: datatypes.OperationTimes{
					CreatedAt: base.Timestamp{Time: kctx.GenTS.UTC()},
				},
			},
		}
		key, _ := kctx.MarshalText()
		ctx := httptransport.ContextWithRouteMetaID(ctx, "MockOperatorID")

		t.Run("#Success", func(t *testing.T) {
			if runtime.GOOS == `darwin` {
				return
			}
			patch := patch_models.AccessKeyFetchByRand(gomonkey.NewPatches(), m, nil)
			defer patch.Reset()
			d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)

			idAny, err, canBeValidated := access_key.Validate(ctx, string(key))

			NewWithT(t).Expect(canBeValidated).To(BeTrue())
			NewWithT(t).Expect(err).To(BeNil())

			idVal, ok := idAny.(*models.AccessKey)
			NewWithT(t).Expect(ok).To(BeTrue())
			NewWithT(t).Expect(idVal.IdentityID).To(Equal(id))
		})

		patch := gomonkey.NewPatches()
		defer patch.Reset()

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

			if runtime.GOOS == `darwin` {
				return
			}

			t.Run("#GenTsNotMatch", func(t *testing.T) {
				overwrite := *m
				overwrite.CreatedAt.Time = overwrite.CreatedAt.Add(time.Second)
				patch = patch_models.AccessKeyFetchByRand(patch, &overwrite, nil)
				_, err, _ := access_key.Validate(ctx, string(key))
				mock_sqlx.ExpectError(t, err, status.InvalidAccessKey)
			})

			patch = patch_models.AccessKeyFetchByRand(patch, m, nil)
			t.Run("#AccessKeyPermissionDenied", func(t *testing.T) {
				ctx := httptransport.ContextWithRouteMetaID(ctx, "NotEqualMockOperatorID")
				_, err, _ := access_key.Validate(ctx, string(key))
				mock_sqlx.ExpectError(t, err, status.AccessKeyPermissionDenied)
			})

			t.Run("#UpdateLastUsedFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, errFrom(t.Name())).Times(1)
				_, err, _ := access_key.Validate(ctx, string(key))
				NewWithT(t).Expect(err).To(BeNil())
			})

			t.Run("#AccessKeyExpired", func(t *testing.T) {
				time.Sleep(3 * time.Second)
				_, err, _ := access_key.Validate(ctx, string(key))
				mock_sqlx.ExpectError(t, err, status.AccessKeyExpired)
			})
		})

	})

	t.Run("#List", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#ListFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(errFrom(t.Name())).Times(1)
				_, err := access_key.List(ctx, &access_key.ListReq{})
				mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
			})
			t.Run("#CountFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(errFrom(t.Name())).Times(1)
				_, err := access_key.List(ctx, &access_key.ListReq{})
				mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
			})
		})

		result := []models.AccessKey{
			{},
			{
				AccessKeyInfo: models.AccessKeyInfo{
					ExpiredAt: base.Timestamp{Time: time.Now()},
					LastUsed:  base.Timestamp{Time: time.Now()},
				},
			},
		}

		patch = patch_models.AccessKeyList(patch, result, nil)
		patch = patch_models.AccessKeyCount(patch, 2, nil)
		t.Run("#Success", func(t *testing.T) {
			rsp, err := access_key.List(ctx, &access_key.ListReq{})
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(rsp.Total).To(Equal(int64(2)))
			NewWithT(t).Expect(len(rsp.Data)).To(Equal(2))
		})
	})
}
