package event

import (
	"context"
	"encoding/json"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/types"
)

type HandleEvent struct {
	httpx.MethodPost
	event.EventReq
}

func (r *HandleEvent) Path() string {
	return "/:channel"
}

func (r *HandleEvent) Output(ctx context.Context) (interface{}, error) {
	ctx, l := logr.Start(ctx, "api.HandleEvent")
	defer l.End()

	r.EventReq.SetDefault()

	prj, err := project.GetByName(ctx, r.Channel)
	if err != nil {
		return nil, err
	}
	ctx = types.WithProject(ctx, prj)

	if r.IsDataPush() {
		l.Info("event push")
		// require account auth
		ca, ok := middleware.CurrentAccountFromContext(ctx)
		if !ok {
			return nil, status.InvalidDataPushShouldAccount
		}
		ctx = ca.WithAccount(ctx)
		if ca.Role != enums.ACCOUNT_ROLE__ADMIN && ca.AccountID != prj.AccountID {
			return nil, status.NoProjectPermission.StatusErr().
				WithDesc("data push request should have account permission")
		}
		reqs := event.DataPushReqs{}
		err = json.Unmarshal(r.Payload.Bytes(), &reqs)
		if err != nil {
			return nil, status.InvalidDataPushPayload.StatusErr().WithDesc(err.Error())
		}
		reqs.SetDefault()
		return event.BatchCreate(ctx, reqs)
	}

	// require publisher auth
	l.Info("event request")
	pub, ok := middleware.MaybePublisher(ctx)
	if !ok {
		return nil, status.InvalidDataPushShouldPublisher
	}
	if pub.ProjectID != prj.ProjectID {
		return nil, status.NoProjectPermission.StatusErr().
			WithDesc("publisher request should have project permission")
	}
	ctx = types.WithPublisher(ctx, pub.Publisher)
	return event.Create(ctx, &r.EventReq)
}
