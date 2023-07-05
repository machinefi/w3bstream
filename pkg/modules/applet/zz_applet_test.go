package applet_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"reflect"
	"runtime"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
	mock_sqlx "github.com/machinefi/w3bstream/pkg/test/mock_depends_kit_sqlx"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

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
	)(context.Background())

	d.MockDBExecutor.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
	d.MockDBExecutor.EXPECT().Context().Return(ctx).AnyTimes()
	d.MockTxExecutor.EXPECT().IsTx().Return(true).AnyTimes()

	t.Run("GetBySFID", func(t *testing.T) {
		t.Run("#Success", func(t *testing.T) {
			d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)

			_, err := applet.GetBySFID(ctx, idg.MustGenSFID())
			NewWithT(t).Expect(err).To(BeNil())
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#AppletNotFound", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrNotFound).MaxTimes(1)

				_, err := applet.GetBySFID(ctx, idg.MustGenSFID())
				mock_sqlx.ExpectError(t, err, status.AppletNotFound)
			})
			t.Run("#DatabaseError", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrDatabase).MaxTimes(1)

				_, err := applet.GetBySFID(ctx, idg.MustGenSFID())
				mock_sqlx.ExpectError(t, err, status.DatabaseError)
			})
		})
	})

	t.Run("RemoveBySFID", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		t.Run("#Success", func(t *testing.T) {
			patch := gomonkey.ApplyMethod(
				reflect.TypeOf(&models.Applet{}),
				"DeleteByAppletID",
				func(_ *models.Applet, _ sqlx.DBExecutor) error {
					return nil
				},
			).ApplyFunc(
				strategy.Remove,
				func(_ context.Context, _ *strategy.CondArgs) error {
					return nil
				},
			).ApplyFunc(
				deploy.RemoveByAppletSFID,
				func(_ context.Context, _ types.SFID) error {
					return nil
				},
			)
			defer patch.Reset()

			err := applet.RemoveBySFID(ctx, idg.MustGenSFID())
			NewWithT(t).Expect(err).To(BeNil())
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#DeleteByAppletIDFailed", func(t *testing.T) {
				patch := gomonkey.NewPatches().ApplyMethod(
					reflect.TypeOf(&models.Applet{}),
					"DeleteByAppletID",
					func(_ *models.Applet, _ sqlx.DBExecutor) error {
						return errors.New("any")
					},
				)
				defer patch.Reset()

				err := applet.RemoveBySFID(ctx, idg.MustGenSFID())
				mock_sqlx.ExpectError(t, err, status.DatabaseError)
			})
			t.Run("#RemoveStrategyFailed", func(t *testing.T) {
				patch := gomonkey.ApplyMethod(
					reflect.TypeOf(&models.Applet{}),
					"DeleteByAppletID",
					func(_ *models.Applet, _ sqlx.DBExecutor) error {
						return nil
					},
				).ApplyFunc(
					strategy.Remove,
					func(_ context.Context, _ *strategy.CondArgs) error {
						return errors.New("any")
					},
				)
				defer patch.Reset()
				err := applet.RemoveBySFID(ctx, idg.MustGenSFID())
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
			t.Run("#RemoveInstanceFailed", func(t *testing.T) {
				patch := gomonkey.ApplyMethod(
					reflect.TypeOf(&models.Applet{}),
					"DeleteByAppletID",
					func(_ *models.Applet, _ sqlx.DBExecutor) error {
						return nil
					},
				).ApplyFunc(
					strategy.Remove,
					func(_ context.Context, _ *strategy.CondArgs) error {
						return nil
					},
				).ApplyFunc(
					deploy.RemoveByAppletSFID,
					func(_ context.Context, _ types.SFID) error {
						return errors.New("any")
					},
				)
				defer patch.Reset()

				err := applet.RemoveBySFID(ctx, idg.MustGenSFID())
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
		})
	})

	t.Run("Remove", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		t.Run("#Success", func(t *testing.T) {
			patch := gomonkey.ApplyMethod(
				reflect.TypeOf(&models.Applet{}),
				"List",
				func(_ *models.Applet, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.Applet, error) {
					return []models.Applet{{}}, nil
				},
			).ApplyFunc(
				applet.RemoveBySFID,
				func(_ context.Context, _ types.SFID) error {
					return nil
				},
			)
			defer patch.Reset()

			err := applet.Remove(ctx, &applet.CondArgs{})
			NewWithT(t).Expect(err).To(BeNil())
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#ListAppletFailed", func(t *testing.T) {
				patch := gomonkey.ApplyMethod(
					reflect.TypeOf(&models.Applet{}),
					"List",
					func(_ *models.Applet, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.Applet, error) {
						return nil, errors.New("any")
					},
				)
				defer patch.Reset()

				err := applet.Remove(ctx, &applet.CondArgs{})
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
			t.Run("#BatchRemoveAppletFailed", func(t *testing.T) {
				patch := gomonkey.ApplyMethod(
					reflect.TypeOf(&models.Applet{}),
					"List",
					func(_ *models.Applet, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.Applet, error) {
						return []models.Applet{{}}, nil
					},
				).ApplyFunc(
					applet.RemoveBySFID,
					func(_ context.Context, _ types.SFID) error {
						return errors.New("any")
					},
				)
				defer patch.Reset()

				err := applet.Remove(ctx, &applet.CondArgs{})
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
		})
	})

	t.Run("List", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		t.Run("#Success", func(t *testing.T) {
			patch := gomonkey.ApplyMethod(
				reflect.TypeOf(&models.Applet{}),
				"List",
				func(_ *models.Applet, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.Applet, error) {
					return []models.Applet{{}}, nil
				},
			).ApplyMethod(
				reflect.TypeOf(&models.Applet{}),
				"Count",
				func(_ *models.Applet, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) (int64, error) {
					return 1, nil
				},
			)
			defer patch.Reset()

			_, err := applet.List(ctx, &applet.ListReq{})
			NewWithT(t).Expect(err).To(BeNil())
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#ListFailed", func(t *testing.T) {
				patch := gomonkey.ApplyMethod(
					reflect.TypeOf(&models.Applet{}),
					"List",
					func(_ *models.Applet, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.Applet, error) {
						return nil, errors.New("any")
					},
				)
				defer patch.Reset()
				_, err := applet.List(ctx, &applet.ListReq{})
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
			t.Run("#CountFailed", func(t *testing.T) {
				patch := gomonkey.ApplyMethod(
					reflect.TypeOf(&models.Applet{}),
					"List",
					func(_ *models.Applet, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.Applet, error) {
						return []models.Applet{{}}, nil
					},
				).ApplyMethod(
					reflect.TypeOf(&models.Applet{}),
					"Count",
					func(_ *models.Applet, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) (int64, error) {
						return 0, errors.New("any")
					},
				)
				defer patch.Reset()

				_, err := applet.List(ctx, &applet.ListReq{})
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
		})
	})

	t.Run("ListDetail", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		t.Run("#Success", func(t *testing.T) {
			patch := gomonkey.ApplyFunc(
				applet.List,
				func(_ context.Context, _ *applet.ListReq) (*applet.ListRsp, error) {
					return &applet.ListRsp{
						Data:  []models.Applet{{}},
						Total: 1,
					}, nil
				},
			).ApplyFunc(
				deploy.GetByAppletSFID,
				func(_ context.Context, _ types.SFID) (*models.Instance, error) {
					return &models.Instance{}, nil
				},
			).ApplyFunc(
				resource.GetBySFID,
				func(_ context.Context, _ types.SFID) (*models.Resource, error) {
					return &models.Resource{}, nil
				},
			)
			defer patch.Reset()

			_, err := applet.ListDetail(ctx, &applet.ListReq{})
			NewWithT(t).Expect(err).To(BeNil())
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#ListFailed", func(t *testing.T) {
				patch := gomonkey.ApplyFunc(
					applet.List,
					func(_ context.Context, _ *applet.ListReq) (*applet.ListRsp, error) {
						return nil, errors.New("any")
					},
				)
				defer patch.Reset()

				_, err := applet.ListDetail(ctx, &applet.ListReq{})
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
			t.Run("#InstanceNotFound", func(t *testing.T) {
				patch := gomonkey.ApplyFunc(
					applet.List,
					func(_ context.Context, _ *applet.ListReq) (*applet.ListRsp, error) {
						return &applet.ListRsp{
							Data:  []models.Applet{{}},
							Total: 1,
						}, nil
					},
				).ApplyFunc(
					deploy.GetByAppletSFID,
					func(_ context.Context, _ types.SFID) (*models.Instance, error) {
						return nil, errors.New("any")
					},
				)
				defer patch.Reset()

				rsp, err := applet.ListDetail(ctx, &applet.ListReq{})
				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(len(rsp.Data)).To(Equal(0))
			})
			t.Run("#FetchResourceFailed", func(t *testing.T) {
				patch := gomonkey.ApplyFunc(
					applet.List,
					func(_ context.Context, _ *applet.ListReq) (*applet.ListRsp, error) {
						return &applet.ListRsp{
							Data:  []models.Applet{{}},
							Total: 1,
						}, nil
					},
				).ApplyFunc(
					deploy.GetByAppletSFID,
					func(_ context.Context, _ types.SFID) (*models.Instance, error) {
						return &models.Instance{}, nil
					},
				).ApplyFunc(
					resource.GetBySFID,
					func(_ context.Context, _ types.SFID) (*models.Resource, error) {
						return nil, errors.New("any")
					},
				)
				defer patch.Reset()

				_, err := applet.ListDetail(ctx, &applet.ListReq{})
				NewWithT(t).Expect(err).NotTo(BeNil())

			})
		})
	})

	t.Run("Create", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		ctx := contextx.WithContextCompose(
			types.WithProjectContext(&models.Project{}),
			types.WithAccountContext(&models.Account{}),
		)(ctx)
		t.Run("#Success", func(t *testing.T) {
			patch := gomonkey.ApplyFunc(
				resource.Create,
				func(_ context.Context, _ types.SFID, _ *multipart.FileHeader, _, _ string) (*models.Resource, []byte, error) {
					return &models.Resource{}, []byte("any"), nil
				},
			).ApplyMethod(
				reflect.TypeOf(&models.Applet{}),
				"Create",
				func(_ *models.Applet, _ sqlx.DBExecutor) error {
					return nil
				},
			).ApplyFunc(
				strategy.BatchCreate,
				func(_ context.Context, _ []models.Strategy) error {
					return nil
				},
			).ApplyFunc(
				deploy.UpsertByCode,
				func(_ context.Context, _ *deploy.CreateReq, _ []byte, _ enums.InstanceState, _ ...types.SFID) (*models.Instance, error) {
					return &models.Instance{}, nil
				},
			)
			defer patch.Reset()

			_, err := applet.Create(ctx, &applet.CreateReq{})
			NewWithT(t).Expect(err).To(BeNil())
		})
		t.Run("#Failed", func(t *testing.T) {
			patch := gomonkey.ApplyFunc(
				resource.Create,
				func(_ context.Context, _ types.SFID, _ *multipart.FileHeader, _, _ string) (*models.Resource, []byte, error) {
					return nil, nil, errors.New("any")
				},
			)
			defer patch.Reset()

			t.Run("#CreateResourceFailed", func(t *testing.T) {
				_, err := applet.Create(ctx, &applet.CreateReq{})
				NewWithT(t).Expect(err).NotTo(BeNil())
			})

			patch = patch.ApplyFunc(
				resource.Create,
				func(_ context.Context, _ types.SFID, _ *multipart.FileHeader, _, _ string) (*models.Resource, []byte, error) {
					return &models.Resource{}, []byte("any"), nil
				},
			)
			t.Run("#CreateAppletFailed", func(t *testing.T) {
				t.Run("#AppletNameConflict", func(t *testing.T) {
					patch = patch.ApplyMethod(
						reflect.TypeOf(&models.Applet{}),
						"Create",
						func(_ *models.Applet, _ sqlx.DBExecutor) error {
							return mock_sqlx.ErrConflict
						},
					)

					_, err := applet.Create(ctx, &applet.CreateReq{})
					mock_sqlx.ExpectError(t, err, status.AppletNameConflict)
				})
				t.Run("#DatabaseError", func(t *testing.T) {
					patch = patch.ApplyMethod(
						reflect.TypeOf(&models.Applet{}),
						"Create",
						func(_ *models.Applet, _ sqlx.DBExecutor) error {
							return mock_sqlx.ErrDatabase
						},
					)

					_, err := applet.Create(ctx, &applet.CreateReq{})
					mock_sqlx.ExpectError(t, err, status.DatabaseError)
				})
			})

			patch = patch.ApplyMethod(
				reflect.TypeOf(&models.Applet{}),
				"Create",
				func(_ *models.Applet, _ sqlx.DBExecutor) error {
					return nil
				},
			)
			t.Run("#BatchCreateStrategyFailed", func(t *testing.T) {
				patch = patch.ApplyFunc(
					strategy.BatchCreate,
					func(_ context.Context, _ []models.Strategy) error {
						return errors.New("any")
					},
				)

				_, err := applet.Create(ctx, &applet.CreateReq{})
				NewWithT(t).Expect(err).NotTo(BeNil())
			})

			patch = patch.ApplyFunc(
				strategy.BatchCreate,
				func(_ context.Context, _ []models.Strategy) error {
					return nil
				},
			)
			t.Run("#CreateDeployFailed", func(t *testing.T) {
				patch = patch.ApplyFunc(
					deploy.UpsertByCode,
					func(_ context.Context, _ *deploy.CreateReq, _ []byte, _ enums.InstanceState, _ ...types.SFID) (*models.Instance, error) {
						return nil, errors.New("any")
					},
				)
				_, err := applet.Create(ctx, &applet.CreateReq{})
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
		})
	})

	t.Run("Update", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}

		header, err := transformer.NewFileHeader("file", "any_file", bytes.NewBuffer([]byte("any_content")))
		NewWithT(t).Expect(err).To(BeNil())

		cases := []*struct {
			name string
			req  *applet.UpdateReq
		}{
			{
				name: "$NoNeedAnyFields",
				req:  &applet.UpdateReq{},
			},
			{
				name: "$NeedUpdateAllFields",
				req: &applet.UpdateReq{
					File: header,
					Info: applet.Info{
						Strategies: []models.StrategyInfo{},
						AppletName: "any_applet_name",
						WasmCache:  &wasm.Cache{},
					},
				},
			},
		}

		ctx := contextx.WithContextCompose(
			types.WithMgrDBExecutorContext(d),
			conflog.WithLoggerContext(conflog.Std()),
			confid.WithSFIDGeneratorContext(idg),
			types.WithProjectContext(&models.Project{}),
			types.WithAccountContext(&models.Account{}),
			types.WithAppletContext(&models.Applet{}),
		)(ctx)

		patch := gomonkey.ApplyFunc(
			resource.Create,
			func(_ context.Context, _ types.SFID, _ *multipart.FileHeader, _, _ string) (*models.Resource, []byte, error) {
				return &models.Resource{}, []byte("any_content"), nil
			},
		).ApplyFunc(
			strategy.Remove,
			func(_ context.Context, _ *strategy.CondArgs) error {
				return nil
			},
		).ApplyFunc(
			strategy.BatchCreate,
			func(_ context.Context, _ []models.Strategy) error {
				return nil
			},
		).ApplyMethod(
			reflect.TypeOf(&models.Applet{}),
			"UpdateByAppletID",
			func(_ *models.Applet, _ sqlx.DBExecutor) error {
				return nil
			},
		).ApplyFunc(
			deploy.GetByAppletSFID,
			func(_ context.Context, _ types.SFID) (*models.Instance, error) {
				return &models.Instance{}, nil
			},
		).ApplyFunc(
			deploy.UpsertByCode,
			func(_ context.Context, _ *deploy.CreateReq, _ []byte, _ enums.InstanceState, _ ...types.SFID) (*models.Instance, error) {
				return &models.Instance{}, nil
			},
		)
		defer patch.Reset()

		t.Run("#Success", func(t *testing.T) {

			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					_, err = applet.Update(ctx, c.req)
					NewWithT(t).Expect(err).To(BeNil())
				})
			}
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#CreateResFailed", func(t *testing.T) {
				patch = patch.ApplyFunc(
					resource.Create,
					func(_ context.Context, _ types.SFID, _ *multipart.FileHeader, _, _ string) (*models.Resource, []byte, error) {
						return nil, nil, errors.New("any")
					},
				)
				_, err = applet.Update(ctx, &applet.UpdateReq{File: header})
				NewWithT(t).Expect(err).NotTo(BeNil())
			})

			req := &applet.UpdateReq{
				Info: applet.Info{
					Strategies: []models.StrategyInfo{{}},
				},
			}
			t.Run("#RemoveStrategyFailed", func(t *testing.T) {
				patch = patch.ApplyFunc(
					resource.Create,
					func(_ context.Context, _ types.SFID, _ *multipart.FileHeader, _, _ string) (*models.Resource, []byte, error) {
						return &models.Resource{}, []byte("any_content"), nil
					},
				).ApplyFunc(
					strategy.Remove,
					func(_ context.Context, _ *strategy.CondArgs) error {
						return errors.New("any")
					},
				)
				_, err = applet.Update(ctx, req)
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
			t.Run("#BatchCreateStrategyFailed", func(t *testing.T) {
				patch = patch.ApplyFunc(
					strategy.Remove,
					func(_ context.Context, _ *strategy.CondArgs) error {
						return nil
					},
				).ApplyFunc(
					strategy.BatchCreate,
					func(_ context.Context, _ []models.Strategy) error {
						return errors.New("any")
					},
				)
				_, err = applet.Update(ctx, req)
				NewWithT(t).Expect(err).NotTo(BeNil())
			})

			req = &applet.UpdateReq{
				Info: applet.Info{AppletName: "any_applet_name"},
			}
			t.Run("#UpdateAppletFailed", func(t *testing.T) {
				patch = patch.ApplyFunc(
					strategy.BatchCreate,
					func(_ context.Context, _ []models.Strategy) error {
						return nil
					},
				)
				cases := []*struct {
					name      string
					errReturn error
					errExpect status.Error
				}{
					{
						name:      "#AppletNameConflict",
						errReturn: mock_sqlx.ErrConflict,
						errExpect: status.AppletNameConflict,
					},
					{
						name:      "#DatabaseError",
						errReturn: mock_sqlx.ErrDatabase,
						errExpect: status.DatabaseError,
					},
				}

				for _, c := range cases {
					t.Run(c.name, func(t *testing.T) {
						patch = patch.ApplyMethod(
							reflect.TypeOf(&models.Applet{}),
							"UpdateByAppletID",
							func(_ *models.Applet, _ sqlx.DBExecutor) error {
								return c.errReturn
							},
						)
						_, err = applet.Update(ctx, req)
						NewWithT(t).Expect(err).NotTo(BeNil())
						mock_sqlx.ExpectError(t, err, c.errExpect)
					})
				}
			})
			req = &applet.UpdateReq{File: header}
			t.Run("#GetInstanceFailed", func(t *testing.T) {
				patch = patch.ApplyMethod(
					reflect.TypeOf(&models.Applet{}),
					"UpdateByAppletID",
					func(_ *models.Applet, _ sqlx.DBExecutor) error {
						return nil
					},
				).ApplyFunc(
					deploy.GetByAppletSFID,
					func(_ context.Context, _ types.SFID) (*models.Instance, error) {
						return nil, errors.New("any")
					},
				)
				_, err = applet.Update(ctx, req)
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
			t.Run("#UpsertInstanceFailed", func(t *testing.T) {
				patch = patch.ApplyFunc(
					deploy.GetByAppletSFID,
					func(_ context.Context, _ types.SFID) (*models.Instance, error) {
						return &models.Instance{}, nil
					},
				).ApplyFunc(
					deploy.UpsertByCode,
					func(_ context.Context, _ *deploy.CreateReq, _ []byte, _ enums.InstanceState, _ ...types.SFID) (*models.Instance, error) {
						return nil, errors.New("any")
					},
				)
				_, err = applet.Update(ctx, req)
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
		})
	})
}
