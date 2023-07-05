package patch_modules

import (
	"context"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
)

func OperatorListByCond(patch *gomonkey.Patches, v []models.Operator, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		operator.ListByCond,
		func(_ context.Context, _ *operator.CondArgs) ([]models.Operator, error) { return v, err },
	)
}
