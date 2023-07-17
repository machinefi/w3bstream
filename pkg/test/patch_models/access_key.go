package patch_models

import (
	"reflect"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/models"
)

var _tarModelAccessKey = reflect.TypeOf(&models.AccessKey{})

func AccessKeyFetchByRand(patch *gomonkey.Patches, overwrite *models.AccessKey, err error) *gomonkey.Patches {
	return patch.ApplyMethod(
		_tarModelAccessKey,
		"FetchByRand",
		func(receiver *models.AccessKey, d sqlx.DBExecutor) error {
			if overwrite != nil {
				*receiver = *overwrite
			}
			return nil
		},
	)
}

func AccessKeyList(patch *gomonkey.Patches, v []models.AccessKey, err error) *gomonkey.Patches {
	return patch.ApplyMethod(
		_tarModelAccessKey,
		"List",
		func(_ *models.AccessKey, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.AccessKey, error) {
			return v, err
		},
	)
}

func AccessKeyCount(patch *gomonkey.Patches, n int64, err error) *gomonkey.Patches {
	return patch.ApplyMethod(
		_tarModelAccessKey,
		"Count",
		func(_ *models.AccessKey, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) (int64, error) {
			return n, err
		},
	)
}

func AccessKeyFetchByAccountIDAndName(patch *gomonkey.Patches, overwrite *models.AccessKey, err error) *gomonkey.Patches {
	return patch.ApplyMethod(
		_tarModelAccessKey,
		"FetchByAccountIDAndName",
		func(receiver *models.AccessKey, _ sqlx.DBExecutor) error {
			if overwrite != nil {
				*receiver = *overwrite
			}
			return err
		},
	)
}
