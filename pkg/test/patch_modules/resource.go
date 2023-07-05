package patch_modules

import (
	"context"
	"mime/multipart"

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

func ResourceGetBySFID(patch *gomonkey.Patches, m *models.Resource, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		resource.GetBySFID,
		func(_ context.Context, _ types.SFID) (*models.Resource, error) { return m, err },
	)
}

func ResourceCreate(patch *gomonkey.Patches, m *models.Resource, data []byte, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		resource.Create,
		func(_ context.Context, _ types.SFID, _ *multipart.FileHeader, _, _ string) (*models.Resource, []byte, error) {
			return m, data, err
		},
	)
}
