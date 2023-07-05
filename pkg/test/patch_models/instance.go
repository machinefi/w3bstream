package patch_models

import (
	"reflect"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/models"
)

var _targetModelInstance = reflect.TypeOf(&models.Instance{})

func InstanceList(patch *gomonkey.Patches, data []models.Instance, err error) *gomonkey.Patches {
	return patch.ApplyMethod(
		_targetModelInstance,
		"List",
		func(_ *models.Instance, _ sqlx.DBExecutor, _ builder.SqlCondition, _ ...builder.Addition) ([]models.Instance, error) {
			return data, err
		},
	)
}

func InstanceFetchByInstanceID(patch *gomonkey.Patches, overwrite *models.Instance, err error) *gomonkey.Patches {
	return patch.ApplyMethod(
		_targetModelInstance,
		"FetchByInstanceID",
		func(receiver *models.Instance, _ sqlx.DBExecutor) error {
			if overwrite != nil {
				*receiver = *overwrite
			}
			return err
		},
	)
}

func InstanceFetchByAppletID(patch *gomonkey.Patches, overwrite *models.Instance, err error) *gomonkey.Patches {
	return patch.ApplyMethod(
		_targetModelInstance,
		"FetchByAppletID",
		func(receiver *models.Instance, _ sqlx.DBExecutor) error {
			if overwrite != nil {
				*receiver = *overwrite
			}
			return err
		},
	)
}

func InstanceUpdateByInstanceID(patch *gomonkey.Patches, overwrite *models.Instance, err error) *gomonkey.Patches {
	return patch.ApplyMethod(
		_targetModelInstance,
		"UpdateByInstanceID",
		func(receiver *models.Instance, _ sqlx.DBExecutor, _ ...string) error {
			if overwrite != nil {
				*receiver = *overwrite
			}
			return err
		},
	)
}

func InstanceCreate(patch *gomonkey.Patches, overwrite *models.Instance, err error) *gomonkey.Patches {
	return patch.ApplyMethod(
		_targetModelInstance,
		"Create",
		func(receiver *models.Instance, _ sqlx.DBExecutor) error {
			if overwrite != nil {
				*receiver = *overwrite
			}
			return err
		},
	)
}
