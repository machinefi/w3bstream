package patch_models

import (
	"reflect"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/models"
)

func ConfigList(patch *gomonkey.Patches, data []models.Config, err error) *gomonkey.Patches {
	return patch.ApplyMethod(
		reflect.TypeOf(&models.Config{}),
		"List",
		func(_ *models.Config, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.Config, error) {
			return data, err
		},
	)
}
