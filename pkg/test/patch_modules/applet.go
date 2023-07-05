package patch_modules

import (
	"context"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
)

func AppletRemoveBySFID(patch *gomonkey.Patches, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		applet.RemoveBySFID,
		func(_ context.Context, _ types.SFID) error { return err },
	)
}

func AppletList(patch *gomonkey.Patches, rsp *applet.ListRsp, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		applet.List,
		func(_ context.Context, _ *applet.ListReq) (*applet.ListRsp, error) {
			return rsp, err
		},
	)
}
