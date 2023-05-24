package projectoperator

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func IsOperatorOccupied(ctx context.Context, operatorID types.SFID) (bool, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.ProjectOperator{}
	pos, err := m.List(d, m.ColOperatorID().Eq(operatorID))
	if err != nil {
		return false, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return len(pos) > 0, nil
}

func GetByProject(ctx context.Context, projectID types.SFID) (*models.ProjectOperator, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.ProjectOperator{
		RelProject: models.RelProject{ProjectID: projectID},
	}

	if err := m.FetchByProjectID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.ProjectOperatorNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func RemoveByProject(ctx context.Context, projectID types.SFID) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.ProjectOperator{
		RelProject: models.RelProject{ProjectID: projectID},
	}
	if err := m.DeleteByProjectID(d); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func Create(ctx context.Context, projectID, operatorID types.SFID) (*models.ProjectOperator, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)

	m := &models.ProjectOperator{
		RelProject:  models.RelProject{ProjectID: projectID},
		RelOperator: models.RelOperator{OperatorID: operatorID},
	}

	if err := m.Create(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return nil, status.ProjectOperatorConflict
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}
