package patch_modules

import (
	"context"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/types"
)

func ResourceGetContentBySFID(patch *gomonkey.Patches, m *models.Resource, data []byte, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		resource.GetContentBySFID,
		func(_ context.Context, _ types.SFID) (*models.Resource, []byte, error) {
			return m, data, err
		},
	)
}
