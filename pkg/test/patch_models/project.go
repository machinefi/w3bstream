package patch_models

import (
	"reflect"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/models"
)

var _targetModelProject = reflect.TypeOf(&models.Project{})

func ProjectFetchByProjectID(patch *gomonkey.Patches, overwrite *models.Project, err error) *gomonkey.Patches {
	return patch.ApplyMethod(
		_targetModelProject,
		"FetchByProjectID",
		func(receiver *models.Project, _ sqlx.DBExecutor) error {
			if overwrite != nil {
				*receiver = *overwrite
			}
			return err
		},
	)
}
