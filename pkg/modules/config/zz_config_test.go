package config_test

import (
	"context"
	"runtime"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	mock_sqlx "github.com/machinefi/w3bstream/pkg/test/mock_depends_kit_sqlx"
	"github.com/machinefi/w3bstream/pkg/test/patch_models"
	"github.com/machinefi/w3bstream/pkg/test/patch_modules"
	"github.com/machinefi/w3bstream/pkg/test/patch_std"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func TestMarshal(t *testing.T) {
	if runtime.GOOS == `darwin` {
		return
	}

	patch := gomonkey.NewPatches()
	defer patch.Reset()

	patch = patch_std.JsonMarshal(patch, nil, errors.New("any"))
	_, err := config.Marshal(&wasm.Cache{})
	mock_sqlx.ExpectError(t, err, status.ConfigParseFailed, "any")

	patch = patch_std.JsonMarshal(patch, []byte("any"), nil)
	_, err = config.Marshal(&wasm.Cache{})
	NewWithT(t).Expect(err).To(BeNil())
}

func TestUnmarshal(t *testing.T) {
	patch := gomonkey.NewPatches()
	defer patch.Reset()

	t.Run("#Failed", func(t *testing.T) {
		t.Run("#InvalidConfigType", func(t *testing.T) {
			_, err := config.Unmarshal([]byte("any"), enums.ConfigType(100))
			mock_sqlx.ExpectError(t, err, status.InvalidConfigType)
		})
		t.Run("#ConfigParseFailed", func(t *testing.T) {
			if runtime.GOOS == `darwin` {
				return
			}

			patch = patch_std.JsonUnmarshal(patch, errors.New("any"))
			_, err := config.Unmarshal([]byte("any"), enums.CONFIG_TYPE__INSTANCE_CACHE)
			mock_sqlx.ExpectError(t, err, status.ConfigParseFailed, "any")
		})
	})
}

func TestConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := mock_sqlx.NewTx(ctrl)
	ctx := contextx.WithContextCompose(
		types.WithMgrDBExecutorContext(d),
		confid.WithSFIDGeneratorContext(confid.MustNewSFIDGenerator()),
	)(context.Background())

	d.MockDBExecutor.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
	d.MockTxExecutor.EXPECT().IsTx().Return(true).AnyTimes()
	d.MockDBExecutor.EXPECT().Context().Return(ctx).AnyTimes()

	cause := func(msg string) error { return errors.New(msg) }

	t.Run("#GetValueBySFID", func(t *testing.T) {
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#NotFoundError", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).Times(1)
				_, err := config.GetValueBySFID(ctx, 100)
				mock_sqlx.ExpectError(t, err, status.ConfigNotFound)
			})
			t.Run("#DatabaseError", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrDatabase).Times(1)
				_, err := config.GetValueBySFID(ctx, 100)
				mock_sqlx.ExpectError(t, err, status.DatabaseError)
			})
		})
		t.Run("#Success", func(t *testing.T) {
			if runtime.GOOS == `darwin` {
				return
			}
			d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(1)

			patch = patch_std.JsonUnmarshal(patch, nil)
			patch = patch_modules.TypesWasmNewConfigurationByType(patch, &wasm.Cache{}, nil)
			defer patch.Reset()

			_, err := config.GetValueBySFID(ctx, 100)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("#GetValueByRelAndType", func(t *testing.T) {
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#NotFoundError", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).Times(1)
				_, err := config.GetValueByRelAndType(ctx, 100, enums.CONFIG_TYPE__INSTANCE_CACHE)
				mock_sqlx.ExpectError(t, err, status.ConfigNotFound)
			})
			t.Run("#DatabaseError", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrDatabase).Times(1)
				_, err := config.GetValueByRelAndType(ctx, 100, enums.CONFIG_TYPE__INSTANCE_CACHE)
				mock_sqlx.ExpectError(t, err, status.DatabaseError)
			})
		})
		t.Run("#Success", func(t *testing.T) {
			if runtime.GOOS == `darwin` {
				return
			}
			d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(1)

			patch = patch_std.JsonUnmarshal(patch, nil)
			patch = patch_modules.TypesWasmNewConfigurationByType(patch, &wasm.Cache{}, nil)
			defer patch.Reset()

			_, err := config.GetValueByRelAndType(ctx, 100, enums.CONFIG_TYPE__INSTANCE_CACHE)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("#Upsert", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#TxTryFetchAndUninitOldConfig", func(t *testing.T) {
				t.Run("#FetchByRelIDAndTypeFailed", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(cause(t.Name())).Times(1)
					_, err := config.Upsert(ctx, 100, &wasm.Cache{})
					t.Log(err)
					mock_sqlx.ExpectError(t, err, status.DatabaseError)
				})
				t.Run("#UnmarshalConfigFailed", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(1)
					patch = patch_modules.ConfigUnmarshal(patch, nil, cause(t.Name()))
					_, err := config.Upsert(ctx, 100, &wasm.Cache{})
					NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
				})

				patch = patch_modules.ConfigUnmarshal(patch, &wasm.Cache{}, nil)

				t.Run("#UninitConfigFailed", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(1)
					patch = patch_modules.TypesWasmUninitConfiguration(patch, cause(t.Name()))
					_, err := config.Upsert(ctx, 100, &wasm.Cache{})
					t.Log(err)
					mock_sqlx.ExpectError(t, err, status.ConfigUninitFailed)
				})
			})

			patch = patch_modules.TypesWasmUninitConfiguration(patch, nil)

			t.Run("#TxCreateOrUpdateConfig", func(t *testing.T) {
				t.Run("#ConfigMarshalFailed", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).Times(1)
					patch = patch_modules.ConfigMarshal(patch, nil, cause(t.Name()))
					_, err := config.Upsert(ctx, 100, &wasm.Cache{})
					NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
				})

				patch = patch_modules.ConfigMarshal(patch, []byte("any"), nil)

				t.Run("#CreateConfigFailed", func(t *testing.T) {
					t.Run("#ConflictError", func(t *testing.T) {
						d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).Times(1)
						d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrConflict).Times(1)
						_, err := config.Upsert(ctx, 100, &wasm.Cache{})
						mock_sqlx.ExpectError(t, err, status.ConfigConflict)
					})
					t.Run("#DatabaseError", func(t *testing.T) {
						d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).Times(1)
						d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrDatabase).Times(1)
						_, err := config.Upsert(ctx, 100, &wasm.Cache{})
						mock_sqlx.ExpectError(t, err, status.DatabaseError)
					})
				})

				t.Run("#UpdateConfigFailed", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(1)
					d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, cause(t.Name())).Times(1)
					_, err := config.Upsert(ctx, 100, &wasm.Cache{})
					mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
				})
			})

			t.Run("#TxInitConfigFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).Times(1)
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)

				patch = patch_modules.TypesWasmInitConfiguration(patch, cause(t.Name()))
				_, err := config.Upsert(ctx, 100, &wasm.Cache{})
				t.Log(err)
				mock_sqlx.ExpectError(t, err, status.ConfigInitFailed)
			})
		})

		t.Run("#Success", func(t *testing.T) {
			patch = patch_modules.TypesWasmInitConfiguration(patch, nil)
			d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).Times(1)
			d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)

			_, err := config.Upsert(ctx, 100, &wasm.Cache{})
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("#List", func(t *testing.T) {
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#DatabaseListFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(cause(t.Name())).Times(1)
				_, err := config.List(ctx, &config.CondArgs{})
				mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
			})

			t.Run("#ConfigUnmarshalFailed", func(t *testing.T) {
				if runtime.GOOS == `darwin` {
					return
				}

				patch = patch_models.ConfigList(patch, []models.Config{{}}, nil)
				patch = patch_modules.ConfigUnmarshal(patch, nil, cause(t.Name()))
				_, err := config.List(ctx, &config.CondArgs{})
				NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
			})
		})

		t.Run("#Success", func(t *testing.T) {
			if runtime.GOOS == `darwin` {
				return
			}

			patch = patch_modules.ConfigUnmarshal(patch, &wasm.Cache{}, nil)
			_, err := config.List(ctx, &config.CondArgs{})
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("#Remove", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}

		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#ConfigListFailed", func(t *testing.T) {
				patch = patch_modules.ConfigList(patch, nil, cause(t.Name()))

				err := config.Remove(ctx, &config.CondArgs{})
				NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
			})
			t.Run("#RemoveConfigModelFailed", func(t *testing.T) {
				patch = patch_modules.ConfigList(patch, []*config.Detail{{
					Configuration: &wasm.Cache{},
				}}, nil)
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, cause(t.Name())).Times(1)

				err := config.Remove(ctx, &config.CondArgs{})
				mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
			})
			t.Run("#WasmUninitConfigurationFailed", func(t *testing.T) {
				patch = patch_modules.TypesWasmUninitConfiguration(patch, cause(t.Name()))
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)

				err := config.Remove(ctx, &config.CondArgs{})
				t.Log(err)
				mock_sqlx.ExpectError(t, err, status.ConfigUninitFailed)
			})
		})
		t.Run("#Success", func(t *testing.T) {
			d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)
			patch = patch_modules.TypesWasmUninitConfiguration(patch, nil)

			err := config.Remove(ctx, &config.CondArgs{})
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("#Create", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#ConfigMarshalFailed", func(t *testing.T) {
				patch = patch_std.JsonMarshal(patch, nil, cause(t.Name()))
				_, err := config.Create(ctx, 100, &wasm.Cache{})
				mock_sqlx.ExpectError(t, err, status.ConfigParseFailed, t.Name())
			})

			patch = patch_modules.ConfigMarshal(patch, []byte("any"), nil)

			t.Run("#CreateConfigDatabaseFailed", func(t *testing.T) {
				t.Run("#ConflictError", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrConflict).Times(1)
					_, err := config.Create(ctx, 100, &wasm.Cache{})
					mock_sqlx.ExpectError(t, err, status.ConfigConflict)
				})
				t.Run("#DatabaseError", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrDatabase).Times(1)
					_, err := config.Create(ctx, 100, &wasm.Cache{})
					mock_sqlx.ExpectError(t, err, status.DatabaseError)
				})
			})

			t.Run("#WasmInitConfigurationFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)
				patch = patch_modules.TypesWasmInitConfiguration(patch, cause(t.Name()))
				_, err := config.Create(ctx, 100, &wasm.Cache{})
				t.Log(err)
				mock_sqlx.ExpectError(t, err, status.ConfigInitFailed)
			})
		})
		t.Run("#Success", func(t *testing.T) {
			d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)
			patch = patch_modules.TypesWasmInitConfiguration(patch, nil)
			_, err := config.Create(ctx, 100, &wasm.Cache{})
			NewWithT(t).Expect(err).To(BeNil())
		})
	})
}

func TestCondArgs(t *testing.T) {
	for _, c := range []*config.CondArgs{
		{},
		{
			ConfigIDs: types.SFIDs{100, 101},
			RelIDs:    types.SFIDs{100, 101},
			Types: []enums.ConfigType{
				enums.CONFIG_TYPE__PROJECT_DATABASE,
				enums.CONFIG_TYPE__INSTANCE_CACHE,
			}},
	} {
		t.Log(builder.ResolveExpr(c.Condition()).Query())
	}
}

func TestDetail(t *testing.T) {
	d := &config.Detail{RelID: 100, Configuration: &wasm.Cache{}}
	t.Run("#String", func(t *testing.T) {
		t.Log(d.String())
	})
	t.Run("#Log", func(t *testing.T) {
		NewWithT(t).Expect(d.Log(nil)).To(Equal(d.String()))
		t.Log(d.Log(errors.New("any")))
	})
}
