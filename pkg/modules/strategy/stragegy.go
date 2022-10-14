package strategy

import (
	"context"

	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"

	"github.com/iotexproject/w3bstream/pkg/enums"
	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type InstanceHandler struct {
	AppletID   types.SFID
	InstanceID types.SFID
	Handler    string
}

func FindStrategyInstances(ctx context.Context, prjName string, eventType enums.EventType) ([]*InstanceHandler, error) {
	d := types.MustDBExecutorFromContext(ctx)

	mProject := &models.Project{ProjectInfo: models.ProjectInfo{Name: prjName}}

	if err := mProject.FetchByName(d); err != nil {
		return nil, status.CheckDatabaseError(err, "FetchProjectByName")
	}

	mStrategy := &models.Strategy{}

	strategies, err := mStrategy.List(d,
		builder.And(
			mStrategy.ColProjectID().Eq(mProject.ProjectID),
			builder.Or(
				mStrategy.ColEventType().Eq(eventType),
				mStrategy.ColEventType().Eq(enums.EVENT_TYPE__ANY),
			),
		),
	)
	if err != nil {
		return nil, status.CheckDatabaseError(err, "ListStrategy")
	}

	if len(strategies) == 0 {
		return nil, status.NotFound.StatusErr().WithDesc("not found strategy")
	}
	strategiesMap := make(map[types.SFID]*models.Strategy)
	for i := range strategies {
		strategiesMap[strategies[i].AppletID] = &strategies[i]
	}

	appletIDs := make(types.SFIDs, 0, len(strategies))

	for i := range strategies {
		appletIDs = append(appletIDs, strategies[i].AppletID)
	}

	mInstance := &models.Instance{}

	instances, err := mInstance.List(d,
		builder.And(
			mInstance.ColAppletID().In(appletIDs),
			mInstance.ColState().Eq(enums.INSTANCE_STATE__STARTED),
		),
	)
	if err != nil {
		return nil, status.CheckDatabaseError(err, "ListInstances")
	}

	if len(instances) == 0 {
		return nil, status.NotFound.StatusErr().WithDesc("not found instance")
	}

	handlers := make([]*InstanceHandler, 0)

	for _, instance := range instances {
		handlers = append(handlers, &InstanceHandler{
			AppletID:   instance.AppletID,
			InstanceID: instance.InstanceID,
			Handler:    strategiesMap[instance.AppletID].Handler,
		})
	}
	return handlers, nil
}
