package default_setter

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
)

func Set(dft, tar interface{}) error {
	// TODO should traverse all filed if struct
	// make sure tar can set
	rvTar := reflectx.Indirect(reflect.ValueOf(tar))
	rtTar := reflect.TypeOf(tar)
	if !rvTar.CanSet() {
		return errors.Errorf("invalid tar parameter: %v", rtTar)
	}
	rvDft := reflectx.Indirect(reflect.ValueOf(dft))
	rtDft := reflect.TypeOf(dft)

	if rtDft.AssignableTo(rtTar) {
		return errors.Errorf("unassignable from %v to %v", rtDft, rtTar)
	}
	rvTar.Set(rvDft)
	return nil
}
