package version

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
)

type BuildVersion struct {
	Branch   string `json:"branch"`
	CommitID string `json:"commitID"`
	Tag      string `json:"tag"`
}

func (b *BuildVersion) String() string {
	if b.Branch+b.Tag+b.CommitID == "" {
		return "unknown"
	}
	return b.Branch + "@" + b.Tag + "-" + b.CommitID
}

var buildVersion BuildVersion

func init() {
	buildVersion.CommitID = types.Commit
	buildVersion.Branch = types.Branch
	buildVersion.Tag = types.Version
}

type VersionRouter struct {
	httpx.MethodGet
}

func (v *VersionRouter) Path() string { return "/version" }

func (v *VersionRouter) Output(ctx context.Context) (interface{}, error) {
	return buildVersion.String(), nil
}
