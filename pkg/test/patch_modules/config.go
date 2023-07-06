package patch_modules

import (
	"context"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func ConfigRemove(patch *gomonkey.Patches, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		config.Remove,
		func(_ context.Context, _ *config.CondArgs) error { return err },
	)
}
func ConfigCreate(patch *gomonkey.Patches, v *models.Config, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		config.Create,
		func(_ context.Context, _ types.SFID, _ wasm.Configuration) (*models.Config, error) { return v, err },
	)
}

func ConfigList(patch *gomonkey.Patches, v []*config.Detail, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		config.List,
		func(_ context.Context, _ *config.CondArgs) ([]*config.Detail, error) { return v, err },
	)
}
