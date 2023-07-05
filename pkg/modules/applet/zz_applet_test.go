package applet_test

import (
	"bytes"
	"context"
	"errors"
	"runtime"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	mock_sqlx "github.com/machinefi/w3bstream/pkg/test/mock_depends_kit_sqlx"
	"github.com/machinefi/w3bstream/pkg/test/patch_models"
	"github.com/machinefi/w3bstream/pkg/test/patch_modules"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func TestCondArgs_Condition(t *testing.T) {
	expr := (&applet.CondArgs{}).Condition()
	NewWithT(t).Expect(expr.IsNil()).To(BeTrue())

	ex := builder.ResolveExpr(
		(&applet.CondArgs{
			ProjectID: 100,
			AppletIDs: types.SFIDs{100},
			Names:     []string{"test_cond"},
			NameLike:  "test_cond",
			LNameLike: "test_cond",
		}).Condition(),
	)
	NewWithT(t).Expect(len(ex.Args())).To(Equal(5))
	t.Log(ex.Query(), ex.Args())
}

func TestCreateReq_BuildStrategies(t *testing.T) {
	ctx := contextx.WithContextCompose(
		conflog.WithLoggerContext(conflog.Std()),
		confid.WithSFIDGeneratorContext(confid.MustNewSFIDGenerator()),
		types.WithAppletContext(&models.Applet{}),
		types.WithProjectContext(&models.Project{}),
	)(context.Background())

	r := &applet.CreateReq{
		File: nil,
		Info: applet.Info{},
	}

	// len(r.Strategies) == 0 build with default strategy info
	strategies := r.BuildStrategies(ctx)
	NewWithT(t).Expect(len(strategies)).To(Equal(1))

	// len(r.Strategies) == 0 build with input strategy info
	r.Info.Strategies = append(r.Info.Strategies, models.StrategyInfo{})
	strategies = r.BuildStrategies(ctx)
	NewWithT(t).Expect(len(strategies)).To(Equal(len(r.Info.Strategies)))
}

func TestApplet(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	d := &struct {
		*mock_sqlx.MockDBExecutor
		*mock_sqlx.MockTxExecutor
	}{
		MockDBExecutor: mock_sqlx.NewMockDBExecutor(c),
		MockTxExecutor: mock_sqlx.NewMockTxExecutor(c),
	}
	idg := confid.MustNewSFIDGenerator()
	ctx := contextx.WithContextCompose(
		types.WithMgrDBExecutorContext(d),
		conflog.WithLoggerContext(conflog.Std()),
		confid.WithSFIDGeneratorContext(idg),
		types.WithProjectContext(&models.Project{}),
		types.WithAccountContext(&models.Account{}),
		types.WithAppletContext(&models.Applet{}),
	)(context.Background())

	d.MockDBExecutor.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
	d.MockDBExecutor.EXPECT().Context().Return(ctx).AnyTimes()
	d.MockTxExecutor.EXPECT().IsTx().Return(true).AnyTimes()

	t.Run("GetBySFID", func(t *testing.T) {
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#AppletNotFound", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).Times(1)

				_, err := applet.GetBySFID(ctx, idg.MustGenSFID())
				mock_sqlx.ExpectError(t, err, status.AppletNotFound)
			})
			t.Run("#DatabaseError", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrDatabase).Times(1)

				_, err := applet.GetBySFID(ctx, idg.MustGenSFID())
				mock_sqlx.ExpectError(t, err, status.DatabaseError)
			})
		})

		t.Run("#Success", func(t *testing.T) {
			d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(1)

			_, err := applet.GetBySFID(ctx, idg.MustGenSFID())
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	cause := func(msg string) error { return errors.New(msg) }

	t.Run("RemoveBySFID", func(t *testing.T) {
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#DeleteByAppletIDFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, cause(t.Name())).Times(1)
				err := applet.RemoveBySFID(ctx, 100)
				mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
			})
			t.Run("RemoveStrategyFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, cause(t.Name())).Times(1)
				err := applet.RemoveBySFID(ctx, 100)
				mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
			})
			t.Run("#RemoveInstanceFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(2)
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(cause(t.Name())).Times(1)
				err := applet.RemoveBySFID(ctx, idg.MustGenSFID())
				mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
			})
		})

		t.Run("#Success", func(t *testing.T) {
			if runtime.GOOS == `darwin` {
				return
			}
			patch := gomonkey.NewPatches()
			defer patch.Reset()

			d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(2)
			patch = patch_modules.DeployRemoveByAppletSFID(patch, nil)

			err := applet.RemoveBySFID(ctx, idg.MustGenSFID())
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("Remove", func(t *testing.T) {
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		req := &applet.CondArgs{}
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#ListAppletFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(cause(t.Name())).Times(1)
				err := applet.Remove(ctx, req)
				mock_sqlx.ExpectError(t, err, status.DatabaseError)
			})
			t.Run("#BatchRemoveAppletFailed", func(t *testing.T) {
				if runtime.GOOS == `darwin` {
					return
				}

				patch = patch_models.AppletList(patch, []models.Applet{{}}, nil)
				patch = patch_modules.AppletRemoveBySFID(patch, cause(t.Name()))

				err := applet.Remove(ctx, req)
				mock_sqlx.ExpectError(t, err, status.BatchRemoveAppletFailed)
			})
		})
		t.Run("#Success", func(t *testing.T) {
			if runtime.GOOS == `darwin` {
				return
			}
			patch = patch_models.AppletList(patch, nil, nil)

			err := applet.Remove(ctx, req)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("List", func(t *testing.T) {
		req := &applet.ListReq{}
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#ListFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(cause(t.Name())).Times(1)
				_, err := applet.List(ctx, req)
				mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
			})
			t.Run("#CountFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(cause(t.Name())).Times(1)
				_, err := applet.List(ctx, req)
				mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
			})
		})
		t.Run("#Success", func(t *testing.T) {
			d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(2)
			_, err := applet.List(ctx, req)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("ListDetail", func(t *testing.T) {
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		req := &applet.ListReq{}
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#AppletListFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(cause(t.Name())).Times(1)
				_, err := applet.ListDetail(ctx, req)
				mock_sqlx.ExpectError(t, err, status.DatabaseError, t.Name())
			})

			t.Run("#AppendDetailFailed", func(t *testing.T) {
				if runtime.GOOS == `darwin` {
					return
				}
				patch = patch_modules.AppletList(patch, &applet.ListRsp{
					Data:  []models.Applet{{}},
					Total: 1,
				}, nil)
				patch = patch_modules.DeployGetByAppletSFID(patch, &models.Instance{}, nil)
				patch = patch_modules.ResourceGetBySFID(patch, nil, cause(t.Name()+"#1"))

				_, err := applet.ListDetail(ctx, req)
				NewWithT(t).Expect(err.Error()).To(Equal(t.Name() + "#1"))

				patch = patch_modules.DeployGetByAppletSFID(patch, nil, cause(t.Name()))
				_, err = applet.ListDetail(ctx, req)
				NewWithT(t).Expect(err).To(BeNil())
			})
		})

		t.Run("#Success", func(t *testing.T) {
			if runtime.GOOS == `darwin` {
				return
			}
			patch = patch_modules.DeployGetByAppletSFID(patch, &models.Instance{}, nil)
			patch = patch_modules.ResourceGetBySFID(patch, &models.Resource{}, nil)
			_, err := applet.ListDetail(ctx, req)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("Create", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}

		patch := gomonkey.NewPatches()
		defer patch.Reset()

		var (
			req  = &applet.CreateReq{}
			code = []byte("any")
		)
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#CreateResourceFailed", func(t *testing.T) {
				patch = patch_modules.ResourceCreate(patch, nil, nil, cause(t.Name()))
				_, err := applet.Create(ctx, req)
				NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
			})

			patch = patch_modules.ResourceCreate(patch, &models.Resource{}, code, nil)

			t.Run("#CreateAppletFailed", func(t *testing.T) {
				t.Run("#AppletNameConflict", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrConflict).Times(1)
					_, err := applet.Create(ctx, req)
					mock_sqlx.ExpectError(t, err, status.AppletNameConflict)
				})
				t.Run("#DatabaseError", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrDatabase).Times(1)
					_, err := applet.Create(ctx, req)
					mock_sqlx.ExpectError(t, err, status.DatabaseError)
				})
			})

			patch = patch_modules.StrategyBatchCreate(patch, nil)

			t.Run("#DeployUpsertByCodeFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)
				patch = patch_modules.DeployUpsertByCode(patch, nil, cause(t.Name()))
				_, err := applet.Create(ctx, req)
				NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
			})
		})

		t.Run("#Success", func(t *testing.T) {
			d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)
			patch = patch_modules.DeployUpsertByCode(patch, &models.Instance{}, nil)

			r := *req
			r.WasmCache = &wasm.Cache{}
			_, err := applet.Create(ctx, &r)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("Update", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}

		patch := gomonkey.NewPatches()
		defer patch.Reset()

		var (
			req  = &applet.UpdateReq{}
			code = []byte("any")

			header, err = transformer.NewFileHeader("file", "any_file", bytes.NewBuffer(code))
		)
		NewWithT(t).Expect(err).To(BeNil())

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#ResourceCreateFailed", func(t *testing.T) {
				r := *req
				r.File = header

				patch = patch_modules.ResourceCreate(patch, nil, nil, cause(t.Name()))
				_, err = applet.Update(ctx, &r)
				NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
			})

			patch = patch_modules.ResourceCreate(patch, &models.Resource{}, code, nil)

			t.Run("#UpdateStrategyFailed", func(t *testing.T) {
				r := *req
				r.Strategies = []models.StrategyInfo{{}}
				t.Run("#StrategyRemoveFailed", func(t *testing.T) {
					patch = patch_modules.StrategyRemove(patch, cause(t.Name()))
					_, err = applet.Update(ctx, &r)
					NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
				})

				patch = patch_modules.StrategyRemove(patch, nil)

				t.Run("#StrategyBatchCreateFailed", func(t *testing.T) {
					patch = patch_modules.StrategyBatchCreate(patch, cause(t.Name()))
					_, err = applet.Update(ctx, &r)
					NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
				})
			})

			patch = patch_modules.StrategyBatchCreate(patch, nil)

			t.Run("#UpdateAppletDatabaseFailed", func(t *testing.T) {
				r := *req
				r.AppletName = "any"
				r.File = header

				t.Run("#AppletNameConflict", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrConflict).Times(1)
					_, err = applet.Update(ctx, &r)
					mock_sqlx.ExpectError(t, err, status.AppletNameConflict)
				})
				t.Run("#DatabaseError", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrDatabase).Times(1)
					_, err = applet.Update(ctx, &r)
					mock_sqlx.ExpectError(t, err, status.DatabaseError)
				})
			})

			d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(sqlmock.NewResult(1, 1), nil).AnyTimes()

			t.Run("#UpdateAndDeployInstanceFailed", func(t *testing.T) {
				r := *req
				r.File = header
				r.WasmCache = &wasm.Cache{}
				t.Run("#DeployGetByAppletSFIDFailed", func(t *testing.T) {
					patch = patch_modules.DeployGetByAppletSFID(patch, &models.Instance{}, cause(t.Name()))
					_, err = applet.Update(ctx, &r)
					NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
				})

				patch = patch_modules.DeployGetByAppletSFID(patch, &models.Instance{}, nil)

				t.Run("#DeployGetByAppletSFIDFailed", func(t *testing.T) {
					patch = patch_modules.DeployUpsertByCode(patch, nil, cause(t.Name()))
					_, err = applet.Update(ctx, &r)
					NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
				})
			})
		})
		t.Run("#Success", func(t *testing.T) {
			patch = patch_modules.DeployUpsertByCode(patch, &models.Instance{}, nil)

			r := *req
			_, err = applet.Update(ctx, &r)
			NewWithT(t).Expect(err).To(BeNil())

			r.Strategies = []models.StrategyInfo{{}}
			_, err = applet.Update(ctx, &r)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})
}
