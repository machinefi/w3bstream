package deploy_test

import (
	"context"
	"reflect"
	"runtime"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	confredis "github.com/machinefi/w3bstream/pkg/depends/conf/redis"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	"github.com/machinefi/w3bstream/pkg/modules/projectoperator"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/modules/wasmlog"
	mock_sqlx "github.com/machinefi/w3bstream/pkg/test/mock_depends_kit_sqlx"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func TestDeploy(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	var (
		d = &struct {
			*mock_sqlx.MockDBExecutor
			*mock_sqlx.MockTxExecutor
		}{
			MockDBExecutor: mock_sqlx.NewMockDBExecutor(ctl),
			MockTxExecutor: mock_sqlx.NewMockTxExecutor(ctl),
		}
		idg        = confid.MustNewSFIDGenerator()
		mqttBroker *confmqtt.Broker
		mqttClient *confmqtt.Client
	)

	// mock default mqtt context
	{
		mqttBroker = &confmqtt.Broker{}
		mqttBroker.SetDefault()

		mqttClient = &confmqtt.Client{}
	}

	ctx := contextx.WithContextCompose(
		types.WithMgrDBExecutorContext(d),
		conflog.WithLoggerContext(conflog.Std()),
		types.WithLoggerContext(conflog.Std()),
		confid.WithSFIDGeneratorContext(idg),
		types.WithAppletContext(&models.Applet{}),
		types.WithResourceContext(&models.Resource{}),
		types.WithInstanceContext(&models.Instance{}),
		types.WithWasmDBEndpointContext(&base.Endpoint{}),
		types.WithRedisEndpointContext(&confredis.Redis{}),
		types.WithTaskWorkerContext(&mq.TaskWorker{}),
		types.WithTaskBoardContext(&mq.TaskBoard{}),
		types.WithMqttBrokerContext(mqttBroker),
		types.WithETHClientConfigContext(&types.ETHClientConfig{}),
		wasm.WithMQTTClientContext(&wasm.MqttClient{Client: mqttClient}),
	)(context.Background())

	d.MockDBExecutor.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
	d.MockTxExecutor.EXPECT().IsTx().Return(true).AnyTimes()
	d.MockDBExecutor.EXPECT().Context().Return(ctx).AnyTimes()

	var (
		anySFID  = types.SFID(124)
		anyError = errors.New("any")
		anyCode  = []byte("any")
		anyState = enums.INSTANCE_STATE__STARTED
		anyReq   = &deploy.CreateReq{}
	)

	t.Run("#Init", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}

		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#ListFailed", func(t *testing.T) {
				patch = patch.
					ApplyMethod(
						reflect.TypeOf(&models.Instance{}),
						"List",
						func(_ *models.Instance, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.Instance, error) {
							return nil, anyError
						},
					)
				NewWithT(t).Expect(deploy.Init(ctx)).To(Equal(anyError))
			})
			t.Run("#FetchAppletFailed", func(t *testing.T) {
				patch = patch.
					ApplyMethod(
						reflect.TypeOf(&models.Instance{}),
						"List",
						func(_ *models.Instance, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.Instance, error) {
							return []models.Instance{{}}, nil
						},
					).
					ApplyMethod(
						reflect.TypeOf(&models.Applet{}),
						"FetchByAppletID",
						func(_ *models.Applet, _ sqlx.DBExecutor) error {
							return anyError
						},
					)
				NewWithT(t).Expect(deploy.Init(ctx)).To(BeNil())
			})
			t.Run("#FetchResourceFailed", func(t *testing.T) {
				patch = patch.
					ApplyMethod(
						reflect.TypeOf(&models.Applet{}),
						"FetchByAppletID",
						func(_ *models.Applet, _ sqlx.DBExecutor) error {
							return nil
						},
					).
					ApplyFunc(
						resource.GetContentBySFID,
						func(_ context.Context, _ types.SFID) (*models.Resource, []byte, error) {
							return nil, nil, anyError
						},
					)
				NewWithT(t).Expect(deploy.Init(ctx)).To(BeNil())
			})

			ins := &models.Instance{}
			ins.State = enums.INSTANCE_STATE__STARTED
			t.Run("#UpsertInstanceFailed", func(t *testing.T) {
				patch = patch.
					ApplyMethod(
						reflect.TypeOf(&models.Instance{}),
						"List",
						func(_ *models.Instance, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.Instance, error) {
							return []models.Instance{*ins}, nil
						},
					).
					ApplyFunc(
						resource.GetContentBySFID,
						func(_ context.Context, _ types.SFID) (*models.Resource, []byte, error) {
							return &models.Resource{}, anyCode, nil
						},
					).
					ApplyFunc(
						deploy.UpsertByCode,
						func(_ context.Context, _ *deploy.CreateReq, _ []byte, _ enums.InstanceState, _ ...types.SFID) (*models.Instance, error) {
							return nil, anyError
						},
					)
				NewWithT(t).Expect(deploy.Init(ctx)).To(BeNil())
			})
			t.Run("$CreateVMFailed", func(t *testing.T) {
				patch = patch.
					ApplyFunc(
						deploy.UpsertByCode,
						func(_ context.Context, _ *deploy.CreateReq, _ []byte, _ enums.InstanceState, _ ...types.SFID) (*models.Instance, error) {
							ret := *ins
							ret.State = enums.INSTANCE_STATE__STOPPED
							return &ret, nil
						},
					)
				NewWithT(t).Expect(deploy.Init(ctx)).To(BeNil())
			})
		})
		t.Run("#Success", func(t *testing.T) {
			ins := &models.Instance{}
			ins.State = enums.INSTANCE_STATE__STARTED
			patch = patch.
				ApplyMethod(
					reflect.TypeOf(&models.Instance{}),
					"List",
					func(_ *models.Instance, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.Instance, error) {
						return []models.Instance{*ins}, nil
					},
				).
				ApplyFunc(
					deploy.UpsertByCode,
					func(_ context.Context, _ *deploy.CreateReq, _ []byte, _ enums.InstanceState, _ ...types.SFID) (*models.Instance, error) {
						ret := *ins
						ret.State = enums.INSTANCE_STATE__STARTED
						return &ret, nil
					},
				)
			NewWithT(t).Expect(deploy.Init(ctx)).To(BeNil())
		})
	})

	t.Run("#GetBySFID", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		cases := []*struct {
			name        string
			mockError   error
			expectError status.Error
		}{
			{
				name: "#Success",
			},
			{
				name:        "#Failed#InstanceNotFound",
				mockError:   mock_sqlx.ErrNotFound,
				expectError: status.InstanceNotFound,
			},
			{
				name:        "#Failed#DatabaseError",
				mockError:   mock_sqlx.ErrDatabase,
				expectError: status.DatabaseError,
			},
		}

		patch := gomonkey.NewPatches()
		defer patch.Reset()

		for _, c := range cases {
			patch.ApplyMethod(
				reflect.TypeOf(&models.Instance{}),
				"FetchByInstanceID",
				func(_ *models.Instance, _ sqlx.DBExecutor) error {
					return c.mockError
				},
			)

			_, err := deploy.GetBySFID(ctx, anySFID)
			if c.expectError != 0 {
				mock_sqlx.ExpectError(t, err, c.expectError)
			} else {
				NewWithT(t).Expect(err).To(BeNil())
			}

		}
	})

	t.Run("#GetByAppletSFID", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		cases := []*struct {
			name        string
			mockError   error
			expectError status.Error
		}{
			{
				name: "#Success",
			},
			{
				name:        "#Failed#InstanceNotFound",
				mockError:   mock_sqlx.ErrNotFound,
				expectError: status.InstanceNotFound,
			},
			{
				name:        "#Failed#DatabaseError",
				mockError:   mock_sqlx.ErrDatabase,
				expectError: status.DatabaseError,
			},
		}

		patch := gomonkey.NewPatches()
		defer patch.Reset()

		for _, c := range cases {
			patch.ApplyMethod(
				reflect.TypeOf(&models.Instance{}),
				"FetchByAppletID",
				func(_ *models.Instance, _ sqlx.DBExecutor) error {
					return c.mockError
				},
			)

			_, err := deploy.GetByAppletSFID(ctx, anySFID)
			if c.expectError != 0 {
				mock_sqlx.ExpectError(t, err, c.expectError)
			} else {
				NewWithT(t).Expect(err).To(BeNil())
			}

		}
	})

	t.Run("#ListWithCond", func(t *testing.T) {
		t.Run("#WithoutProjectID", func(t *testing.T) {
			if runtime.GOOS == `darwin` {
				return
			}
			patch := gomonkey.NewPatches()
			defer patch.Reset()

			arg := &deploy.CondArgs{ProjectID: 0}

			patch = patch.ApplyMethod(
				reflect.TypeOf(&models.Instance{}),
				"List",
				func(_ *models.Instance, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.Instance, error) {
					return []models.Instance{{}}, nil
				},
			)

			_, err := deploy.ListWithCond(ctx, arg)
			NewWithT(t).Expect(err).To(BeNil())

			patch = patch.ApplyMethod(
				reflect.TypeOf(&models.Instance{}),
				"List",
				func(_ *models.Instance, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.Instance, error) {
					return nil, anyError
				},
			)

			_, err = deploy.ListWithCond(ctx, arg)
			NewWithT(t).Expect(err).NotTo(BeNil())
		})
		t.Run("#WithProjectID", func(t *testing.T) {
			arg := &deploy.CondArgs{ProjectID: 1}
			d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			_, err := deploy.ListWithCond(ctx, arg)
			NewWithT(t).Expect(err).To(BeNil())
			d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(anyError).MaxTimes(1)
			_, err = deploy.ListWithCond(ctx, arg)
			NewWithT(t).Expect(err).NotTo(BeNil())
		})
	})

	t.Run("#RemoveBySFID", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.ApplyMethod(
			reflect.TypeOf(&models.Instance{}),
			"DeleteByInstanceID",
			func(_ *models.Instance, _ sqlx.DBExecutor) error {
				return nil
			},
		).ApplyFunc(
			config.Remove,
			func(_ context.Context, _ *config.CondArgs) error {
				return nil
			},
		).ApplyFunc(
			wasmlog.Remove,
			func(_ context.Context, _ *wasmlog.CondArgs) error {
				return nil
			},
		)
		defer patch.Reset()

		t.Run("#Success", func(t *testing.T) {
			err := deploy.RemoveBySFID(ctx, anySFID)
			NewWithT(t).Expect(err).To(BeNil())
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#DeleteByInstanceIDFailed", func(t *testing.T) {
				patch = patch.ApplyMethod(
					reflect.TypeOf(&models.Instance{}),
					"DeleteByInstanceID",
					func(_ *models.Instance, _ sqlx.DBExecutor) error {
						return anyError
					},
				)
				err := deploy.RemoveBySFID(ctx, anySFID)
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
			t.Run("#RemoveConfigFailed", func(t *testing.T) {
				patch = patch.ApplyMethod(
					reflect.TypeOf(&models.Instance{}),
					"DeleteByInstanceID",
					func(_ *models.Instance, _ sqlx.DBExecutor) error {
						return nil
					},
				).ApplyFunc(
					config.Remove,
					func(_ context.Context, _ *config.CondArgs) error {
						return anyError
					},
				)
				err := deploy.RemoveBySFID(ctx, anySFID)
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
			t.Run("#RemoveWasmLogFailed", func(t *testing.T) {
				patch = patch.ApplyFunc(
					config.Remove,
					func(_ context.Context, _ *config.CondArgs) error {
						return nil
					},
				).ApplyFunc(
					wasmlog.Remove,
					func(_ context.Context, _ *wasmlog.CondArgs) error {
						return anyError
					},
				)
				err := deploy.RemoveBySFID(ctx, anySFID)
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
		})
	})

	t.Run("#RemoveByAppletSFID", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		patch = patch.ApplyFunc(
			deploy.GetByAppletSFID,
			func(_ context.Context, _ types.SFID) (*models.Instance, error) {
				return &models.Instance{}, nil
			},
		).ApplyFunc(
			deploy.RemoveBySFID,
			func(_ context.Context, _ types.SFID) error {
				return nil
			},
		)

		t.Run("#Success", func(t *testing.T) {
			err := deploy.RemoveBySFID(ctx, anySFID)
			NewWithT(t).Expect(err).To(BeNil())
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#GetByAppletIDFailed", func(t *testing.T) {
				patch = patch.ApplyFunc(
					deploy.GetByAppletSFID,
					func(_ context.Context, _ types.SFID) (*models.Instance, error) {
						return nil, anyError
					},
				)
				err := deploy.RemoveByAppletSFID(ctx, anySFID)
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
			t.Run("#RemoveBySFIDFailed", func(t *testing.T) {
				patch = patch.ApplyFunc(
					deploy.GetByAppletSFID,
					func(_ context.Context, _ types.SFID) (*models.Instance, error) {
						return &models.Instance{}, nil
					},
				).ApplyFunc(
					deploy.RemoveBySFID,
					func(_ context.Context, _ types.SFID) error {
						return anyError
					},
				)
				err := deploy.RemoveByAppletSFID(ctx, anySFID)
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
		})
	})

	t.Run("Remove", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		anyArg := &deploy.CondArgs{}

		patch := gomonkey.NewPatches()
		defer patch.Reset()

		patch = patch.
			ApplyFunc(
				deploy.ListWithCond,
				func(_ context.Context, _ *deploy.CondArgs) ([]models.Instance, error) {
					return []models.Instance{{}}, nil
				},
			).
			ApplyFunc(
				deploy.RemoveBySFID,
				func(_ context.Context, _ types.SFID) error {
					return nil
				},
			)

		t.Run("#Success", func(t *testing.T) {
			err := deploy.Remove(ctx, anyArg)
			NewWithT(t).Expect(err).To(BeNil())
		})
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#ListWithCondFailed", func(t *testing.T) {
				patch = patch.
					ApplyFunc(
						deploy.ListWithCond,
						func(_ context.Context, _ *deploy.CondArgs) ([]models.Instance, error) {
							return nil, anyError
						},
					)
				err := deploy.Remove(ctx, anyArg)
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
			t.Run("#RemoveBySFIDFailed", func(t *testing.T) {
				patch = patch.
					ApplyFunc(
						deploy.ListWithCond,
						func(_ context.Context, _ *deploy.CondArgs) ([]models.Instance, error) {
							return []models.Instance{{}}, nil
						},
					).
					ApplyFunc(
						deploy.RemoveBySFID,
						func(_ context.Context, _ types.SFID) error {
							return anyError
						},
					)
				err := deploy.Remove(ctx, anyArg)
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
		})
	})

	t.Run("#UpsertByCode", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("InvalidVMState", func(t *testing.T) {
				_, err := deploy.UpsertByCode(ctx, anyReq, anyCode, 0)
				mock_sqlx.ExpectError(t, err, status.InvalidVMState)
			})
			t.Run("TryFetchAppletFromDB", func(t *testing.T) {
				t.Run("#FetchByAppletIDFailed", func(t *testing.T) {
					patch = patch.
						ApplyMethod(
							reflect.TypeOf(&models.Instance{}),
							"FetchByAppletID",
							func(_ *models.Instance, _ sqlx.DBExecutor) error {
								return mock_sqlx.ErrDatabase
							},
						)
					_, err := deploy.UpsertByCode(ctx, anyReq, anyCode, anyState)
					mock_sqlx.ExpectError(t, err, status.DatabaseError)
				})
				t.Run("#UpdateCheckOldIDFailed", func(t *testing.T) {
					id := idg.MustGenSFID()
					patch = patch.
						ApplyMethod(
							reflect.TypeOf(&models.Instance{}),
							"FetchByAppletID",
							func(receiver *models.Instance, _ sqlx.DBExecutor) error {
								*receiver = models.Instance{}
								receiver.InstanceID = id + 1 // not equal to argument
								return nil
							},
						)
					_, err := deploy.UpsertByCode(ctx, anyReq, anyCode, anyState, id)
					mock_sqlx.ExpectError(t, err, status.InvalidAppletContext)
				})
			})
			t.Run("#UpdateOrCreateInstance", func(t *testing.T) {
				t.Run("#CaseUpdateExistedInstace", func(t *testing.T) {
					t.Run("#UpdateByInstanceIDFailed", func(t *testing.T) {
						id := idg.MustGenSFID()
						patch = patch.
							ApplyMethod(
								reflect.TypeOf(&models.Instance{}),
								"FetchByAppletID",
								func(receiver *models.Instance, _ sqlx.DBExecutor) error {
									*receiver = models.Instance{}
									receiver.InstanceID = id // equal to argument
									return nil
								},
							)
						t.Run("#ConflictError", func(t *testing.T) {
							patch = patch.
								ApplyMethod(
									reflect.TypeOf(&models.Instance{}),
									"UpdateByInstanceID",
									func(_ *models.Instance, _ sqlx.DBExecutor, _ ...string) error {
										return mock_sqlx.ErrConflict
									},
								)
							_, err := deploy.UpsertByCode(ctx, anyReq, anyCode, anyState, id)
							mock_sqlx.ExpectError(t, err, status.MultiInstanceDeployed)
						})
						t.Run("#DatabaseError", func(t *testing.T) {
							patch = patch.
								ApplyMethod(
									reflect.TypeOf(&models.Instance{}),
									"UpdateByInstanceID",
									func(_ *models.Instance, _ sqlx.DBExecutor, _ ...string) error {
										return mock_sqlx.ErrDatabase
									},
								)
							_, err := deploy.UpsertByCode(ctx, anyReq, anyCode, anyState, id)
							mock_sqlx.ExpectError(t, err, status.DatabaseError)
						})
					})
				})
				t.Run("#CreateNewInstance", func(t *testing.T) {
					patch = patch.
						ApplyMethod(
							reflect.TypeOf(&models.Instance{}),
							"FetchByAppletID",
							func(_ *models.Instance, _ sqlx.DBExecutor) error {
								return mock_sqlx.ErrNotFound
							},
						)
					t.Run("#ConflictError", func(t *testing.T) {
						patch = patch.
							ApplyMethod(
								reflect.TypeOf(&models.Instance{}),
								"Create",
								func(_ *models.Instance, _ sqlx.DBExecutor) error {
									return mock_sqlx.ErrConflict
								},
							)
						_, err := deploy.UpsertByCode(ctx, anyReq, anyCode, anyState)
						mock_sqlx.ExpectError(t, err, status.MultiInstanceDeployed)
					})
					t.Run("#DatabaseError", func(t *testing.T) {
						patch = patch.
							ApplyMethod(
								reflect.TypeOf(&models.Instance{}),
								"Create",
								func(_ *models.Instance, _ sqlx.DBExecutor) error {
									return mock_sqlx.ErrDatabase
								},
							)
						_, err := deploy.UpsertByCode(ctx, anyReq, anyCode, anyState)
						mock_sqlx.ExpectError(t, err, status.DatabaseError)
					})
				})
			})

			t.Run("#TryUpdateCacheConfig", func(t *testing.T) {
				req := *anyReq
				req.Cache = &wasm.Cache{}

				id := idg.MustGenSFID()
				patch = patch.
					ApplyMethod(
						reflect.TypeOf(&models.Instance{}),
						"Create",
						func(_ *models.Instance, _ sqlx.DBExecutor) error {
							return nil
						},
					)

				t.Run("#RemovePrevConfigFailed", func(t *testing.T) {
					patch = patch.
						ApplyFunc(
							config.Remove,
							func(_ context.Context, _ *config.CondArgs) error {
								return anyError
							},
						)
					_, err := deploy.UpsertByCode(ctx, &req, anyCode, anyState, id)
					NewWithT(t).Expect(err).To(Equal(anyError))
				})
				t.Run("#CreateNewConfigFailed", func(t *testing.T) {
					patch = patch.
						ApplyFunc(
							config.Remove,
							func(_ context.Context, _ *config.CondArgs) error {
								return nil
							},
						).
						ApplyFunc(
							config.Create,
							func(_ context.Context, _ types.SFID, _ wasm.Configuration) (*models.Config, error) {
								return nil, anyError
							},
						)
					_, err := deploy.UpsertByCode(ctx, &req, anyCode, anyState, id)
					NewWithT(t).Expect(err).To(Equal(anyError))
				})
			})
			t.Run("#RemovePrevVMAndCreateNewOne", func(t *testing.T) {
				id := idg.MustGenSFID()
				patch = patch.
					ApplyMethod(
						reflect.TypeOf(&models.Instance{}),
						"FetchByAppletID",
						func(receiver *models.Instance, _ sqlx.DBExecutor) error {
							*receiver = models.Instance{}
							receiver.InstanceID = id
							return nil
						},
					).
					ApplyMethod(
						reflect.TypeOf(&models.Instance{}),
						"UpdateByInstanceID",
						func(_ *models.Instance, _ sqlx.DBExecutor, _ ...string) error {
							return nil
						},
					).
					ApplyFunc(
						config.Create,
						func(_ context.Context, _ types.SFID, _ wasm.Configuration) (*models.Config, error) {
							return &models.Config{}, nil
						},
					).
					ApplyFunc(
						deploy.WithInstanceRuntimeContext,
						func(_ context.Context) (context.Context, error) {
							return ctx, nil
						},
					)
				t.Run("#ForUpdateToRemoveOldInstanceFailed", func(t *testing.T) {
					patch = patch.
						ApplyFunc(
							vm.DelInstance,
							func(_ context.Context, _ types.SFID) error {
								return anyError
							},
						).
						ApplyFunc(
							vm.NewInstance,
							func(_ context.Context, _ []byte, _ types.SFID, _ enums.InstanceState) error {
								return anyError
							},
						)
					_, err := deploy.UpsertByCode(ctx, anyReq, anyCode, anyState, id)
					mock_sqlx.ExpectError(t, err, status.CreateInstanceFailed)
				})
				t.Run("#WithRuntimeContextFailed", func(t *testing.T) {
					patch = patch.
						ApplyFunc(
							deploy.WithInstanceRuntimeContext,
							func(_ context.Context) (context.Context, error) {
								return nil, anyError
							},
						)
					_, err := deploy.UpsertByCode(ctx, anyReq, anyCode, anyState, id)
					NewWithT(t).Expect(err).To(Equal(anyError))
				})
				t.Run("#NewInstanceFailed", func(t *testing.T) {
					patch = patch.
						ApplyFunc(
							deploy.WithInstanceRuntimeContext,
							func(_ context.Context) (context.Context, error) {
								return context.Background(), nil
							},
						).
						ApplyFunc(
							vm.NewInstance,
							func(_ context.Context, _ []byte, _ types.SFID, _ enums.InstanceState) error {
								return anyError
							},
						)
					_, err := deploy.UpsertByCode(ctx, anyReq, anyCode, anyState, id)
					mock_sqlx.ExpectError(t, err, status.CreateInstanceFailed)
				})
			})
		})

		t.Run("#Success", func(t *testing.T) {
			patch = patch.
				ApplyFunc(
					vm.NewInstance,
					func(_ context.Context, _ []byte, _ types.SFID, _ enums.InstanceState) error {
						return nil
					},
				)
			_, err := deploy.UpsertByCode(ctx, anyReq, anyCode, anyState)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("#Upsert", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.ApplyFunc(
			deploy.UpsertByCode,
			func(_ context.Context, _ *deploy.CreateReq, _ []byte, _ enums.InstanceState, _ ...types.SFID) (*models.Instance, error) {
				return nil, nil
			},
		)
		defer patch.Reset()

		t.Run("#Success", func(t *testing.T) {
			patch.ApplyFunc(
				resource.GetContentBySFID,
				func(_ context.Context, _ types.SFID) (*models.Resource, []byte, error) {
					return nil, anyCode, nil
				},
			)
			_, err := deploy.Upsert(ctx, anyReq, anyState)
			NewWithT(t).Expect(err).To(BeNil())
		})
		t.Run("#Failed", func(t *testing.T) {
			patch.ApplyFunc(
				resource.GetContentBySFID,
				func(_ context.Context, _ types.SFID) (*models.Resource, []byte, error) {
					return nil, nil, anyError
				},
			)
			_, err := deploy.Upsert(ctx, anyReq, anyState)
			NewWithT(t).Expect(err).To(Equal(anyError))
		})
	})

	t.Run("#Deploy", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		ctx = contextx.WithContextCompose()(ctx)

		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#UnknownDeployCmd", func(t *testing.T) {
				err := deploy.Deploy(ctx, enums.DeployCmd(100))
				mock_sqlx.ExpectError(t, err, status.UnknownDeployCommand)
			})

			t.Run("#UpdateInstance", func(t *testing.T) {
				cases := []*struct {
					name        string
					mockError   error
					expectError status.Error
				}{
					{
						name:        "#MultiInstanceDeployed",
						mockError:   mock_sqlx.ErrConflict,
						expectError: status.MultiInstanceDeployed,
					},
					{
						name:        "#InstanceNotFound",
						mockError:   mock_sqlx.ErrNotFound,
						expectError: status.InstanceNotFound,
					},
					{
						name:        "#DatabaseError",
						mockError:   mock_sqlx.ErrDatabase,
						expectError: status.DatabaseError,
					},
				}

				for _, c := range cases {
					t.Run(c.name, func(t *testing.T) {
						patch = patch.ApplyMethod(
							reflect.TypeOf(&models.Instance{}),
							"UpdateByInstanceID",
							func(_ *models.Instance, _ sqlx.DBExecutor, _ ...string) error {
								return c.mockError
							},
						)
						err := deploy.Deploy(ctx, enums.DEPLOY_CMD__HUNGUP)
						mock_sqlx.ExpectError(t, err, c.expectError)
					})
				}
			})
			t.Run("#ExecUpdateVMState", func(t *testing.T) {
				patch = patch.ApplyMethod(
					reflect.TypeOf(&models.Instance{}),
					"UpdateByInstanceID",
					func(_ *models.Instance, _ sqlx.DBExecutor, _ ...string) error {
						return nil
					},
				)

				t.Run("#StopVMFailed", func(t *testing.T) {
					patch = patch.ApplyFunc(
						vm.StopInstance,
						func(_ context.Context, _ types.SFID) error {
							return anyError
						},
					)
					err := deploy.Deploy(ctx, enums.DEPLOY_CMD__HUNGUP)
					NewWithT(t).Expect(err).To(BeNil())
				})
				t.Run("#StartVMFailed", func(t *testing.T) {
					patch = patch.ApplyFunc(
						vm.StartInstance,
						func(_ context.Context, _ types.SFID) error {
							return anyError
						},
					)
					err := deploy.Deploy(ctx, enums.DEPLOY_CMD__START)
					NewWithT(t).Expect(err).To(BeNil())
				})
			})
		})
		t.Run("#Success", func(t *testing.T) {
			err := deploy.Deploy(ctx, enums.DEPLOY_CMD__START)
			t.Log(err)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("#WithInstanceRuntimeContext", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#FetchProjectFailed", func(t *testing.T) {
				patch = patch.
					ApplyMethod(
						reflect.TypeOf(&models.Project{}),
						"FetchByProjectID",
						func(_ *models.Project, _ sqlx.DBExecutor) error {
							return anyError
						},
					)
				_, err := deploy.WithInstanceRuntimeContext(ctx)
				NewWithT(t).Expect(err).To(Equal(anyError))
			})
			t.Run("#FetchResourceFailed", func(t *testing.T) {
				patch = patch.
					ApplyMethod(
						reflect.TypeOf(&models.Project{}),
						"FetchByProjectID",
						func(_ *models.Project, _ sqlx.DBExecutor) error {
							return nil
						},
					).
					ApplyMethod(
						reflect.TypeOf(&models.Resource{}),
						"FetchByResourceID",
						func(_ *models.Resource, _ sqlx.DBExecutor) error {
							return anyError
						},
					)
				_, err := deploy.WithInstanceRuntimeContext(ctx)
				NewWithT(t).Expect(err).To(Equal(anyError))
			})
			t.Run("#ListConfigurationsFailed", func(t *testing.T) {
				patch = patch.
					ApplyMethod(
						reflect.TypeOf(&models.Resource{}),
						"FetchByResourceID",
						func(_ *models.Resource, _ sqlx.DBExecutor) error {
							return nil
						},
					).
					ApplyFunc(
						config.List,
						func(_ context.Context, _ *config.CondArgs) ([]*config.Detail, error) {
							return nil, anyError
						},
					)
				_, err := deploy.WithInstanceRuntimeContext(ctx)
				NewWithT(t).Expect(err).To(Equal(anyError))
			})
			t.Run("#ConfigurationsInitFailed", func(t *testing.T) {
				patch = patch.
					ApplyFunc(
						config.List,
						func(_ context.Context, _ *config.CondArgs) ([]*config.Detail, error) {
							return []*config.Detail{
								{
									Configuration: &wasm.Database{},
								},
							}, nil
						},
					).
					ApplyMethod(
						reflect.TypeOf(&wasm.Database{}),
						"Init",
						func(_ *wasm.Database, _ context.Context) error {
							return anyError
						},
					)
				_, err := deploy.WithInstanceRuntimeContext(ctx)
				mock_sqlx.ExpectError(t, err, status.ConfigInitFailed)
			})
			t.Run("#GetProjectOperatorFailed", func(t *testing.T) {
				patch = patch.
					ApplyMethod(
						reflect.TypeOf(&wasm.Database{}),
						"Init",
						func(_ *wasm.Database, _ context.Context) error {
							return nil
						},
					).
					ApplyFunc(
						projectoperator.GetByProject,
						func(_ context.Context, _ types.SFID) (*models.ProjectOperator, error) {
							return nil, status.DatabaseError
						},
					)
				_, err := deploy.WithInstanceRuntimeContext(ctx)
				mock_sqlx.ExpectError(t, err, status.DatabaseError)
			})
			t.Run("#ListOperatorsFailed", func(t *testing.T) {
				patch = patch.
					ApplyFunc(
						projectoperator.GetByProject,
						func(_ context.Context, _ types.SFID) (*models.ProjectOperator, error) {
							return &models.ProjectOperator{}, nil
						},
					).
					ApplyFunc(
						operator.ListByCond,
						func(_ context.Context, _ *operator.CondArgs) ([]models.Operator, error) {
							return nil, anyError
						},
					)
				_, err := deploy.WithInstanceRuntimeContext(ctx)
				NewWithT(t).Expect(err).To(Equal(anyError))
			})
		})
		t.Run("#Success", func(t *testing.T) {
			patch = patch.
				ApplyFunc(
					operator.ListByCond,
					func(_ context.Context, _ *operator.CondArgs) ([]models.Operator, error) {
						return []models.Operator{}, nil
					},
				).
				ApplyFunc(
					wasm.NewChainClient,
					func(_ context.Context, _ []models.Operator, _ *models.ProjectOperator) *wasm.ChainClient {
						return &wasm.ChainClient{}
					},
				)
			_, err := deploy.WithInstanceRuntimeContext(ctx)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})
}

func TestCondArgs_Condition(t *testing.T) {
	cases := []*deploy.CondArgs{
		{ProjectID: 0},
		{InstanceIDs: []types.SFID{1}},
		{InstanceIDs: []types.SFID{1, 2}},
		{AppletIDs: []types.SFID{1}},
		{AppletIDs: []types.SFID{1, 2}},
		{States: []enums.InstanceState{1}},
		{States: []enums.InstanceState{1, 2}},
	}

	for _, c := range cases {
		t.Log(builder.ResolveExpr(c.Condition()).Query())
	}
}
