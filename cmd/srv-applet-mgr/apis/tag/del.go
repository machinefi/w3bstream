package tag

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/tag"
	"github.com/machinefi/w3bstream/pkg/types"
)

type RemoveTag struct {
	httpx.MethodDelete
	ProjectName string `in:"path" name:"projectName"`
	tag.RemoveTagReq
}

func (r *RemoveTag) Path() string { return "/:projectName" }

func (r *RemoveTag) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	return nil, tag.RemoveTag(ctx, types.MustProjectFromContext(ctx), &r.RemoveTagReq)
}
