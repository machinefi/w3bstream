package publisher

import (
	"context"
	"time"

	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
	"github.com/machinefi/w3bstream/pkg/types"
)

func CreateAnonymousPublisher(ctx context.Context) (*models.Publisher, error) {
	ctx, l := logr.Start(ctx, "modules.AnonymousPublisher.Create")
	defer l.End()

	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		prj = types.MustProjectFromContext(ctx)
		idg = confid.MustSFIDGeneratorFromContext(ctx)

		pub *models.Publisher
		tok *access_key.CreateRsp
		err error
	)

	id := idg.MustGenSFID()
	tok, err = access_key.Create(ctx, &access_key.CreateReq{
		IdentityID:   id,
		IdentityType: enums.ACCESS_KEY_IDENTITY_TYPE__PUBLISHER,
		CreateReqBase: access_key.CreateReqBase{
			Name: "pub_" + id.String(),
			Desc: "pub_" + id.String(),
			Privileges: access_key.GroupAccessPrivileges{{
				Name: enums.ApiGroupEvent,
				Perm: enums.ACCESS_PERMISSION__READ_WRITE,
			}},
		},
	})
	if err != nil {
		return nil, err
	}

	pub = &models.Publisher{
		PrimaryID:    datatypes.PrimaryID{ID: 0},
		RelProject:   models.RelProject{ProjectID: prj.ProjectID},
		RelPublisher: models.RelPublisher{PublisherID: tok.IdentityID},
		PublisherInfo: models.PublisherInfo{
			Name:  "anonymous",
			Key:   "anonymous",
			Token: tok.AccessKey,
		},
	}
	pub.CreatedAt.Set(time.Now())
	pub.UpdatedAt.Set(time.Now())

	found := false
	err = sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) (err error) {
			prj.Public = datatypes.TRUE
			if err := prj.UpdateByName(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.ProjectNameConflict
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) (err error) {
			if err := pub.FetchByID(d); err != nil {
				if sqlx.DBErr(err).IsNotFound() {
					found = false
					return nil
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			found = true
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if found {
				return nil
			}

			if _, err = d.Exec(
				builder.Insert().Into(d.T(pub)).Values(builder.Cols(
					pub.ColID().Name, pub.ColProjectID().Name, pub.ColPublisherID().Name, pub.ColName().Name,
					pub.ColKey().Name, pub.ColToken().Name, pub.ColCreatedAt().Name, pub.ColUpdatedAt().Name),
					pub.ID, pub.ProjectID, pub.PublisherID, pub.Name, pub.Key, pub.Token, pub.CreatedAt, pub.UpdatedAt),
			); err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()

	if err != nil {
		l.Warn(errors.Wrap(err, "insert anonymous publisher"))
		return nil, err
	}

	if jwt.WithAnonymousPublisherFn == nil {
		jwt.SetWithAnonymousPublisherFn(SetAnonymousPublisher)
	}

	return pub, nil
}

func SetAnonymousPublisher(ctx context.Context, tok string) (string, error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)

		pub = &models.Publisher{PrimaryID: datatypes.PrimaryID{ID: 0}}
	)

	opId := httptransport.OperationIDFromContext(ctx)

	if tok == "" && opId == "HandleEvent" {
		if err := pub.FetchByID(d); err != nil {
			return "", status.PublisherNotFound
		}
		tok = pub.Token
	}

	return tok, nil
}
