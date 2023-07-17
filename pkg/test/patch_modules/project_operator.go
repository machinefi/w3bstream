package patch_modules

import (
	"context"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/projectoperator"
)

func ProjectOperatorGetByProject(patch *gomonkey.Patches, v *models.ProjectOperator, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		projectoperator.GetByProject,
		func(_ context.Context, _ types.SFID) (*models.ProjectOperator, error) { return v, err },
	)
}
