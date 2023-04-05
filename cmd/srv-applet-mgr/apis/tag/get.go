package tag

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/tag"
	"github.com/machinefi/w3bstream/pkg/types"
)

type ListTag struct {
	httpx.MethodGet
	ProjectName string `in:"path" name:"projectName"`
	tag.ListTagReq
}

func (r *ListTag) Path() string { return "/:projectName" }

func (r *ListTag) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}

	return tag.ListTags(ctx, types.MustProjectFromContext(ctx), &r.ListTagReq)
}
