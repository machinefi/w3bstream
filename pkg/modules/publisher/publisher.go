package publisher

import (
	"context"

	confid "github.com/iotexproject/Bumblebee/conf/id"
	"github.com/iotexproject/Bumblebee/conf/jwt"

	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type CreatePublisherReq = models.PublisherInfo

func CreatePublisher(ctx context.Context, projectID types.SFID, r *CreatePublisherReq) (m *models.Publisher, err error) {
	d := types.MustDBExecutorFromContext(ctx)
	j := jwt.MustConfFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	// TODO generate token, maybe use public key
	r.Token, err = j.GenerateTokenByPayload(projectID)
	if err != nil {
		return nil, status.InternalServerError.StatusErr().WithDesc(err.Error())
	}

	m = &models.Publisher{
		RelProject:    models.RelProject{ProjectID: projectID},
		RelPublisher:  models.RelPublisher{PublisherID: idg.MustGenSFID()},
		PublisherInfo: *r,
	}
	if err = m.Create(d); err != nil {
		return nil, err
	}

	return m, nil
}
