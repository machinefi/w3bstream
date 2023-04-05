package tag

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/tag"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateTag struct {
	httpx.MethodPost
	ProjectName      string `in:"path" name:"projectName"`
	tag.CreateTagReq `in:"body" mime:"multipart"`
}

func (r *CreateTag) Path() string { return "/:projectName" }

func (r *CreateTag) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	return tag.CreateTag(ctx, types.MustProjectFromContext(ctx), &r.CreateTagReq)
}
