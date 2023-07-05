package patch_models

import (
	"reflect"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/models"
)

var _targetModelApplet = reflect.TypeOf(&models.Applet{})

func AppletFetchByAppletID(patch *gomonkey.Patches, overwrite *models.Applet, err error) *gomonkey.Patches {
	return patch.ApplyMethod(
		_targetModelApplet,
		"FetchByAppletID",
		func(receiver *models.Applet, _ sqlx.DBExecutor) error {
			if overwrite != nil {
				*receiver = *overwrite
			}
			return err
		},
	)
}
