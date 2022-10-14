package applet

import (
	"context"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/middleware"

	"github.com/iotexproject/w3bstream/pkg/modules/applet"
)

type CreateApplet struct {
	httpx.MethodPost
	ProjectID              types.SFID `in:"path" name:"projectID"`
	applet.CreateAppletReq `in:"body" mime:"multipart"`
}

func (r *CreateApplet) Path() string { return "/:projectID" }

func (r *CreateApplet) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	if _, err := ca.ValidateProjectPerm(ctx, r.ProjectID); err != nil {
		return nil, err
	}

	return applet.CreateApplet(ctx, r.ProjectID, &r.CreateAppletReq)
}

type CreateAppletAndDeploy struct{}
