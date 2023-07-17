package patch_modules

import (
	"context"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
)

func StrategyRemove(patch *gomonkey.Patches, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		strategy.Remove,
		func(_ context.Context, _ *strategy.CondArgs) error { return err },
	)
}

func StrategyBatchCreate(patch *gomonkey.Patches, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		strategy.BatchCreate,
		func(_ context.Context, _ []models.Strategy) error { return err },
	)
}
