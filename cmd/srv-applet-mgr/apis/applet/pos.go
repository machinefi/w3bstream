package applet

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateApplet struct {
	httpx.MethodPost
	ProjectName            string `in:"path" name:"projectName"`
	applet.CreateAppletReq `in:"body" mime:"multipart"`
}

func (r *CreateApplet) Path() string { return "/:projectName" }

func (r *CreateApplet) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)
	return applet.CreateApplet(ctx, prj.ProjectID, &r.CreateAppletReq)
}

type CreateAppletAndDeploy struct{}
