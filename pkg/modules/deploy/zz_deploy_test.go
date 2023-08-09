package deploy_test

import (
	"context"
	"runtime"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	confredis "github.com/machinefi/w3bstream/pkg/depends/conf/redis"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	wasmapi "github.com/machinefi/w3bstream/pkg/modules/vm/wasmapi/types/mock"
	mock_sqlx "github.com/machinefi/w3bstream/pkg/test/mock_depends_kit_sqlx"
	"github.com/machinefi/w3bstream/pkg/test/patch_models"
	"github.com/machinefi/w3bstream/pkg/test/patch_modules"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

var (
	anyError    = errors.New("any")
	anySFID     = types.SFID(124)
	anyWasmCode = []byte("any")
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

	wasmApiServer := wasmapi.NewMockServer(ctl)

	ctx := contextx.WithContextCompose(
		types.WithMgrDBExecutorContext(d),
		conflog.WithLoggerContext(conflog.Std()),
		types.WithLoggerContext(conflog.Std()),
		confid.WithSFIDGeneratorContext(idg),
		types.WithAppletContext(&models.Applet{}),
		types.WithResourceContext(&models.Resource{}),
		types.WithInstanceContext(&models.Instance{}),
		types.WithWasmDBConfigContext(&types.WasmDBConfig{}),
		types.WithRedisEndpointContext(&confredis.Redis{}),
		types.WithTaskWorkerContext(&mq.TaskWorker{}),
		types.WithTaskBoardContext(&mq.TaskBoard{}),
		types.WithMqttBrokerContext(mqttBroker),
		types.WithETHClientConfigContext(&types.ETHClientConfig{}),
		types.WithChainConfigContext(&types.ChainConfig{}),
		wasm.WithMQTTClientContext(mqttClient),
		types.WithWasmApiServerContext(wasmApiServer),
	)(context.Background())

	d.MockDBExecutor.EXPECT().T(gomock.Any()).Return(&builder.Table{}).AnyTimes()
	d.MockTxExecutor.EXPECT().IsTx().Return(true).AnyTimes()
	d.MockDBExecutor.EXPECT().Context().Return(ctx).AnyTimes()

	errFrom := func(from string) error {
		return errors.New(from)
	}

	t.Run("#Init", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}

		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#ListInstanceFailed", func(t *testing.T) {
				from := "models.Instance.List"
				patch = patch_models.InstanceList(patch, nil, errFrom(from))
				NewWithT(t).Expect(deploy.Init(ctx).Error()).To(Equal(from))
			})

			patch = patch_models.InstanceList(patch, []models.Instance{{
				InstanceInfo: models.InstanceInfo{State: enums.INSTANCE_STATE__STARTED},
			}}, nil)

			t.Run("#FetchAppletFailed", func(t *testing.T) {
				from := "models.Applet.FetchByAppletID"
				patch = patch_models.AppletFetchByAppletID(patch, nil, errFrom(from))
				NewWithT(t).Expect(deploy.Init(ctx)).To(BeNil())
			})

			patch = patch_models.AppletFetchByAppletID(patch, &models.Applet{}, nil)

			t.Run("#FetchResourceFailed", func(t *testing.T) {
				from := "resource.GetContentBySFID"
				patch = patch_modules.ResourceGetContentBySFID(patch, nil, nil, errFrom(from))
				NewWithT(t).Expect(deploy.Init(ctx)).To(BeNil())
			})

			patch = patch_modules.ResourceGetContentBySFID(patch, &models.Resource{}, anyWasmCode, nil)

			t.Run("#UpsertInstanceFailed", func(t *testing.T) {
				from := "deploy.UpsertByCode"
				patch = patch_modules.DeployUpsertByCode(patch, nil, errFrom(from))
				NewWithT(t).Expect(deploy.Init(ctx)).To(BeNil())
			})

			t.Run("#UnexpectedDeployedVMState", func(t *testing.T) {
				ins := &models.Instance{}
				patch = patch_modules.DeployUpsertByCode(patch, ins, nil)
				NewWithT(t).Expect(deploy.Init(ctx)).To(BeNil())
			})
		})
		t.Run("#Success", func(t *testing.T) {
			ins := &models.Instance{
				InstanceInfo: models.InstanceInfo{State: enums.INSTANCE_STATE__STARTED},
			}
			patch = patch_modules.DeployUpsertByCode(patch, ins, nil)
			NewWithT(t).Expect(deploy.Init(ctx)).To(BeNil())
		})
	})

	t.Run("#GetBySFID", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#InstanceNotFound", func(t *testing.T) {
				patch = patch_models.InstanceFetchByInstanceID(patch, nil, mock_sqlx.ErrNotFound)
				_, err := deploy.GetBySFID(ctx, anySFID)
				mock_sqlx.ExpectError(t, err, status.InstanceNotFound)
			})
			t.Run("#DatabaseError", func(t *testing.T) {
				from := "models.Instance.FetchByInstanceID"
				patch = patch_models.InstanceFetchByInstanceID(patch, nil, errFrom(from))
				_, err := deploy.GetBySFID(ctx, anySFID)
				mock_sqlx.ExpectError(t, err, status.DatabaseError, from)
			})
		})

		t.Run("#Success", func(t *testing.T) {
			patch = patch_models.InstanceFetchByInstanceID(patch, &models.Instance{}, nil)
			_, err := deploy.GetBySFID(ctx, anySFID)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("#GetByAppletSFID", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#InstanceNotFound", func(t *testing.T) {
				patch = patch_models.InstanceFetchByAppletID(patch, nil, mock_sqlx.ErrNotFound)
				_, err := deploy.GetByAppletSFID(ctx, anySFID)
				mock_sqlx.ExpectError(t, err, status.InstanceNotFound)
			})
			t.Run("#DatabaseError", func(t *testing.T) {
				from := "models.Instance.FetchByAppletID"
				patch = patch_models.InstanceFetchByAppletID(patch, nil, errFrom(from))
				_, err := deploy.GetByAppletSFID(ctx, anySFID)
				mock_sqlx.ExpectError(t, err, status.DatabaseError, from)
			})
		})

		t.Run("#Success", func(t *testing.T) {
			patch = patch_models.InstanceFetchByAppletID(patch, &models.Instance{}, nil)
			_, err := deploy.GetByAppletSFID(ctx, anySFID)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("#List", func(t *testing.T) {
		r := &deploy.ListReq{}
		t.Run("#Failed", func(t *testing.T) {
			from := "ListError"
			d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(errFrom(from)).Times(1)
			_, err := deploy.List(ctx, r)
			mock_sqlx.ExpectError(t, err, status.DatabaseError, from)

			from = "CountError"
			d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(errFrom(from)).Times(1)
			r.ProjectID = 1000
			_, err = deploy.List(ctx, r)
			mock_sqlx.ExpectError(t, err, status.DatabaseError, from)
		})
		t.Run("#Success", func(t *testing.T) {
			d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(2)
			_, err := deploy.List(ctx, &deploy.ListReq{})
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("#ListByCond", func(t *testing.T) {
		r := &deploy.CondArgs{ProjectID: 1000}
		t.Run("#Failed", func(t *testing.T) {
			d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(mock_sqlx.ErrDatabase).Times(1)
			_, err := deploy.ListByCond(ctx, r)
			mock_sqlx.ExpectError(t, err, status.DatabaseError)
		})
		t.Run("#Success", func(t *testing.T) {
			d.MockDBExecutor.EXPECT().QueryAndScan(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			_, err := deploy.ListByCond(ctx, r)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("#RemoveBySFID", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#TxDeleteByInstanceIDFailed", func(t *testing.T) {
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrDatabase).Times(1)
				err := deploy.RemoveBySFID(ctx, anySFID)
				mock_sqlx.ExpectError(t, err, status.DatabaseError)
			})

			t.Run("#TxConfigRemoveFailed", func(t *testing.T) {
				from := "config.Remove"
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)
				patch = patch_modules.ConfigRemove(patch, errFrom(from))
				err := deploy.RemoveBySFID(ctx, anySFID)
				NewWithT(t).Expect(err.Error()).To(Equal(from))
			})

			patch = patch_modules.ConfigRemove(patch, nil)

			t.Run("#TxWasmLogRemoveFailed", func(t *testing.T) {
				from := "wasmlog.Remove"
				d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)
				patch = patch_modules.WasmLogRemove(patch, errFrom(from))
				err := deploy.RemoveBySFID(ctx, anySFID)
				NewWithT(t).Expect(err.Error()).To(Equal(from))
			})
		})
		t.Run("#Success", func(t *testing.T) {
			d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, nil).Times(1)
			patch = patch_modules.WasmLogRemove(patch, nil)
			err := deploy.RemoveBySFID(ctx, anySFID)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("#RemoveByAppletSFID", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#TxDeployGetByAppletSFIDFailed", func(t *testing.T) {
				from := "deploy.GetByAppletSFID"
				patch = patch_modules.DeployGetByAppletSFID(patch, nil, errFrom(from))
				NewWithT(t).Expect(deploy.RemoveByAppletSFID(ctx, anySFID).Error()).To(Equal(from))
			})

			patch = patch_modules.DeployGetByAppletSFID(patch, &models.Instance{}, nil)

			t.Run("#TxDeployRemoveBySFIDFailed", func(t *testing.T) {
				from := "deploy.RemoveBySFID"
				patch = patch_modules.DeployRemoveBySFID(patch, errFrom(from))
				NewWithT(t).Expect(deploy.RemoveByAppletSFID(ctx, anySFID).Error()).To(Equal(from))
			})
		})

		t.Run("#Success", func(t *testing.T) {
			patch = patch_modules.DeployRemoveBySFID(patch, nil)
			NewWithT(t).Expect(deploy.RemoveByAppletSFID(ctx, anySFID)).To(BeNil())
		})
	})

	t.Run("#Remove", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		r := &deploy.CondArgs{}
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#TxDeployListByCondFailed", func(t *testing.T) {
				from := "deploy.ListByCond"
				patch = patch_modules.DeployListByCond(patch, nil, errFrom(from))
				NewWithT(t).Expect(deploy.Remove(ctx, r).Error()).To(Equal(from))
			})

			patch = patch_modules.DeployListByCond(patch, []models.Instance{{}}, nil)

			t.Run("#TxDeployRemoveBySFIDFailed", func(t *testing.T) {
				from := "deploy.RemoveBySFID"
				patch = patch_modules.DeployRemoveBySFID(patch, errFrom(from))
				NewWithT(t).Expect(deploy.Remove(ctx, r).Error()).To(Equal(from))
			})
		})
		t.Run("#Success", func(t *testing.T) {
			patch = patch_modules.DeployRemoveBySFID(patch, nil)
			NewWithT(t).Expect(deploy.Remove(ctx, r)).To(BeNil())
		})
	})

	t.Run("#UpsertByCode", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		updateCasePatch := func(overwrite *models.Instance) {
			patch = patch_models.InstanceFetchByAppletID(patch, overwrite, nil)
		}
		createCasePatch := func() {
			patch = patch_models.InstanceFetchByAppletID(patch, nil, mock_sqlx.ErrNotFound)
		}

		req := &deploy.CreateReq{}
		code := []byte("any")
		state := enums.INSTANCE_STATE__STARTED

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#InvalidInstanceState", func(t *testing.T) {
				_, err := deploy.UpsertByCode(ctx, req, code, enums.InstanceState(0))
				mock_sqlx.ExpectError(t, err, status.InvalidVMState)
			})
			t.Run("#TxFetchByAppletIDFailed", func(t *testing.T) {
				t.Run("#DatabaseError", func(t *testing.T) {
					from := "models.Instance.FetchByAppletID"
					patch = patch_models.InstanceFetchByAppletID(patch, nil, errFrom(from))
					_, err := deploy.UpsertByCode(ctx, req, code, state)
					mock_sqlx.ExpectError(t, err, status.DatabaseError, from)
				})
				t.Run("#InvalidAppletContext", func(t *testing.T) {
					updateCasePatch(&models.Instance{
						RelInstance: models.RelInstance{InstanceID: anySFID + 1},
					})
					_, err := deploy.UpsertByCode(ctx, req, code, state, anySFID)
					mock_sqlx.ExpectError(t, err, status.InvalidAppletContext)
				})
			})
			t.Run("#TxUpdateOrCreateInstanceFailed", func(t *testing.T) {
				t.Run("#CaseUpdateFailed", func(t *testing.T) {
					updateCasePatch(&models.Instance{
						RelInstance: models.RelInstance{InstanceID: anySFID},
					})
					t.Run("#MultiInstanceDeployed", func(t *testing.T) {
						patch = patch_models.InstanceUpdateByInstanceID(patch, nil, mock_sqlx.ErrConflict)
						_, err := deploy.UpsertByCode(ctx, req, code, state, anySFID)
						mock_sqlx.ExpectError(t, err, status.MultiInstanceDeployed)
					})
					t.Run("#DatabaseError", func(t *testing.T) {
						from := "models.Instance.UpdateByInstanceID"
						patch = patch_models.InstanceUpdateByInstanceID(patch, nil, errFrom(from))
						_, err := deploy.UpsertByCode(ctx, req, code, state, anySFID)
						mock_sqlx.ExpectError(t, err, status.DatabaseError, from)
					})
				})
				t.Run("#CaseCreateFailed", func(t *testing.T) {
					createCasePatch()
					t.Run("#MultiInstanceDeployed", func(t *testing.T) {
						patch = patch_models.InstanceCreate(patch, nil, mock_sqlx.ErrConflict)
						_, err := deploy.UpsertByCode(ctx, req, code, state, anySFID)
						mock_sqlx.ExpectError(t, err, status.MultiInstanceDeployed)
					})
					t.Run("#DatabaseError", func(t *testing.T) {
						patch = patch_models.InstanceCreate(patch, nil, mock_sqlx.ErrDatabase)
						_, err := deploy.UpsertByCode(ctx, req, code, state, anySFID)
						mock_sqlx.ExpectError(t, err, status.DatabaseError)
					})
				})

				ins := &models.Instance{RelInstance: models.RelInstance{InstanceID: anySFID}}
				patch = patch_models.InstanceUpdateByInstanceID(patch, ins, nil)
				patch = patch_models.InstanceCreate(patch, ins, nil)
			})

			t.Run("#TxUpdateConfigFailed", func(t *testing.T) {
				req.Cache = &wasm.Cache{}

				t.Run("#ConfigRemoveFailed", func(t *testing.T) {
					from := "config.Remove"
					patch = patch_modules.ConfigRemove(patch, errFrom(from))
					_, err := deploy.UpsertByCode(ctx, req, code, state)
					NewWithT(t).Expect(err.Error()).To(Equal(from))
				})

				patch = patch_modules.ConfigRemove(patch, nil)

				t.Run("#ConfigCreateFailed", func(t *testing.T) {
					from := "config.Create"
					patch = patch_modules.ConfigCreate(patch, nil, errFrom(from))
					_, err := deploy.UpsertByCode(ctx, req, code, state)
					NewWithT(t).Expect(err.Error()).To(Equal(from))
				})

				patch = patch_modules.ConfigCreate(patch, nil, nil)
			})

			t.Run("#TxCreateInstanceVM", func(t *testing.T) {
				updateCasePatch(&models.Instance{})
				t.Run("#WihtInstanceRuntimeContextFailed", func(t *testing.T) {
					from := "deploy.WithInstanceRuntimeContext"
					patch = patch_modules.DeployWithInstanceRuntimeContext(patch, nil, errFrom(from))
					_, err := deploy.UpsertByCode(ctx, req, code, state)
					NewWithT(t).Expect(err.Error()).To(Equal(from))
				})
				patch = patch_modules.DeployWithInstanceRuntimeContext(patch, nil, nil)
				t.Run("#VmNewInstanceFailed", func(t *testing.T) {
					from := "vm.NewInstance"
					patch = patch_modules.VmNewInstance(patch, errFrom(from))
					_, err := deploy.UpsertByCode(ctx, req, code, state)
					mock_sqlx.ExpectError(t, err, status.CreateInstanceFailed, from)
				})
			})
		})

		t.Run("#Success", func(t *testing.T) {
			patch = patch_modules.VmNewInstance(patch, nil)
			_, err := deploy.UpsertByCode(ctx, req, code, state)
			NewWithT(t).Expect(err).To(BeNil())
		})
	})

	t.Run("#Upsert", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		req := &deploy.CreateReq{}
		state := enums.INSTANCE_STATE__STARTED

		t.Run("#Failed", func(t *testing.T) {
			from := "resource.GetContentBySFID"
			patch = patch_modules.ResourceGetContentBySFID(patch, nil, nil, errFrom(from))
			_, err := deploy.Upsert(ctx, req, state)
			NewWithT(t).Expect(err.Error()).To(Equal(from))
		})
		from := "deploy.UpsertByCode"
		patch = patch_modules.ResourceGetContentBySFID(patch, nil, nil, nil)
		patch = patch_modules.DeployUpsertByCode(patch, nil, errFrom(from))
		_, err := deploy.Upsert(ctx, req, state)
		NewWithT(t).Expect(err.Error()).To(Equal(from))
	})

	t.Run("#Deploy", func(t *testing.T) {
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#UnknownDeployCmd", func(t *testing.T) {
				err := deploy.Deploy(ctx, enums.DeployCmd(100))
				mock_sqlx.ExpectError(t, err, status.UnknownDeployCommand, "100")
			})

			t.Run("#TxUpdateInstanceFailed", func(t *testing.T) {
				t.Run("#MultiInstanceDeployed", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrConflict).Times(1)
					err := deploy.Deploy(ctx, enums.DEPLOY_CMD__HUNGUP)
					mock_sqlx.ExpectError(t, err, status.MultiInstanceDeployed)
				})
				t.Run("#InstanceNotFound", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrNotFound).Times(1)
					err := deploy.Deploy(ctx, enums.DEPLOY_CMD__HUNGUP)
					mock_sqlx.ExpectError(t, err, status.InstanceNotFound)
				})
				t.Run("#DatabaseError", func(t *testing.T) {
					d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(nil, mock_sqlx.ErrDatabase).Times(1)
					err := deploy.Deploy(ctx, enums.DEPLOY_CMD__HUNGUP)
					mock_sqlx.ExpectError(t, err, status.DatabaseError)
				})
			})
		})
		d.MockDBExecutor.EXPECT().Exec(gomock.Any()).Return(sqlmock.NewResult(1, 1), nil).Times(2)
		err := deploy.Deploy(ctx, enums.DEPLOY_CMD__HUNGUP)
		NewWithT(t).Expect(err).To(BeNil())

		err = deploy.Deploy(ctx, enums.DEPLOY_CMD__START)
		NewWithT(t).Expect(err).To(BeNil())
	})

	t.Run("#WithInstanceRuntimeContext", func(t *testing.T) {
		if runtime.GOOS == `darwin` {
			return
		}
		patch := gomonkey.NewPatches()
		defer patch.Reset()

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#FetchProjectFailed", func(t *testing.T) {
				patch = patch_models.ProjectFetchByProjectID(patch, nil, errFrom(t.Name()))
				_, err := deploy.WithInstanceRuntimeContext(ctx)
				NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
			})

			patch = patch_models.ProjectFetchByProjectID(patch, &models.Project{}, nil)

			t.Run("#GetProjectOperatorFailed", func(t *testing.T) {
				patch = patch_modules.ProjectOperatorGetByProject(patch, nil, errFrom(t.Name()))
				_, err := deploy.WithInstanceRuntimeContext(ctx)
				NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
			})

			patch = patch_modules.ProjectOperatorGetByProject(patch, &models.ProjectOperator{}, nil)

			t.Run("#ListAccountOperatorsFailed", func(t *testing.T) {
				patch = patch_modules.OperatorListByCond(patch, nil, errFrom(t.Name()))
				_, err := deploy.WithInstanceRuntimeContext(ctx)
				NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
			})

			patch = patch_modules.OperatorListByCond(patch, []models.Operator{{}}, nil)

			t.Run("#ListConfigurationsFailed", func(t *testing.T) {
				patch = patch_modules.ConfigList(patch, nil, errFrom(t.Name()))
				_, err := deploy.WithInstanceRuntimeContext(ctx)
				NewWithT(t).Expect(err.Error()).To(Equal(t.Name()))
			})

			patch = patch_modules.ConfigList(patch, []*config.Detail{{Configuration: &wasm.Cache{}}}, nil)

			t.Run("#ConfigurationsInitFailed", func(t *testing.T) {
				patch_modules.TypesWasmInitConfiguration(patch, errFrom(t.Name()))
				_, err := deploy.WithInstanceRuntimeContext(ctx)
				mock_sqlx.ExpectError(t, err, status.ConfigInitFailed)
			})

			patch = patch_modules.TypesWasmInitConfiguration(patch, nil)

			t.Run("#GlobalConfigurationInitFailed", func(t *testing.T) {
				patch_modules.TypesWasmInitGlobalConfiguration(patch, errFrom(t.Name()))
				_, err := deploy.WithInstanceRuntimeContext(ctx)
				mock_sqlx.ExpectError(t, err, status.ConfigInitFailed)
			})
		})

		t.Run("#Success", func(t *testing.T) {
			patch = patch_modules.TypesWasmInitGlobalConfiguration(patch, nil)
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
