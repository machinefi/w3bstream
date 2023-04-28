package version

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
)

type VersionRouter struct {
	httpx.MethodGet
}

func (v *VersionRouter) Path() string { return "/version" }

func (v *VersionRouter) Output(ctx context.Context) (interface{}, error) {
	return types.BuildVersion, nil
}
