package patch_modules

import (
	"context"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
)

func DeployGetBySFID(patch *gomonkey.Patches, m *models.Instance, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		deploy.GetBySFID,
		func(_ context.Context, _ types.SFID) (*models.Instance, error) { return m, err },
	)
}

func DeployUpsertByCode(patch *gomonkey.Patches, m *models.Instance, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		deploy.UpsertByCode,
		func(_ context.Context, _ *deploy.CreateReq, _ []byte, _ enums.InstanceState, _ ...types.SFID) (*models.Instance, error) {
			return m, err
		},
	)
}

func DeployRemoveBySFID(patch *gomonkey.Patches, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		deploy.RemoveBySFID,
		func(_ context.Context, _ types.SFID) error { return err },
	)
}

func DeployGetByAppletSFID(patch *gomonkey.Patches, v *models.Instance, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		deploy.GetByAppletSFID,
		func(_ context.Context, _ types.SFID) (*models.Instance, error) { return v, err },
	)
}

func DeployListByCond(patch *gomonkey.Patches, v []models.Instance, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		deploy.ListByCond,
		func(_ context.Context, _ *deploy.CondArgs) ([]models.Instance, error) { return v, err },
	)
}

func DeployWithInstanceRuntimeContext(patch *gomonkey.Patches, ctx context.Context, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		deploy.WithInstanceRuntimeContext,
		func(_ context.Context) (context.Context, error) { return ctx, err },
	)
}
