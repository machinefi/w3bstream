package wasm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/types"
)

type TraceInfo struct {
	ProjectName string
	AppletName  string
	InstanceID  types.SFID
}

func (v *TraceInfo) GlobalConfigType() ConfigType { return ConfigTraceInfo }

func (v *TraceInfo) Init(parent context.Context) error {
	prj := types.MustProjectFromContext(parent)
	app := types.MustAppletFromContext(parent)
	ins := types.MustInstanceFromContext(parent)

	v.ProjectName = prj.Name
	v.AppletName = app.Name
	v.InstanceID = ins.InstanceID
	return nil
}

func (v *TraceInfo) WithContext(ctx context.Context) context.Context {
	return WithTraceInfo(ctx, v)
}
