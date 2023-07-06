package patch_models

import (
	"reflect"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/models"
)

var _targetModelResource = reflect.TypeOf(&models.Resource{})

func ResourceFetchByResourceID(patch *gomonkey.Patches, overwrite *models.Resource, err error) *gomonkey.Patches {
	return patch.ApplyMethod(
		_targetModelResource,
		"FetchByResourceID",
		func(receiver *models.Resource, _ sqlx.DBExecutor) error {
			if overwrite != nil {
				*receiver = *overwrite
			}
			return err
		},
	)
}
