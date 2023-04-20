package resource

import (
	"context"
	"mime/multipart"

	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func FetchOrCreateResource(ctx context.Context, f *multipart.FileHeader) (*models.Resource, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "FetchOrCreateResource")
	defer l.End()

	_, fullName, md5, err := Upload(ctx, f, idg.MustGenSFID().String())
	if err != nil {
		l.Error(err)
		return nil, status.UploadFileFailed.StatusErr().WithDesc(err.Error())
	}

	m := &models.Resource{ResourceInfo: models.ResourceInfo{Md5: md5}}

	if err = m.FetchByMd5(d); err != nil && sqlx.DBErr(err).IsNotFound() {
		m.ResourceID = idg.MustGenSFID()
		m.ResourceInfo.Path = fullName
		m.ResourceInfo.Md5 = md5
		if err = m.Create(d); err != nil {
			l.Error(errors.Wrap(err, "create wasm resource db failed"))
			return nil, status.CheckDatabaseError(err, "CreateResource")
		}
		l.Info("wasm resource created")
	}
	l.Info("get wasm resource from db")
	return m, err
}

func CheckResourceExist(ctx context.Context, path string) bool {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "CheckResourceExist")
	defer l.End()

	return IsPathExists(path)
}

func ListResource(ctx context.Context) ([]models.Resource, error) {
	res, err := (&models.Resource{}).List(types.MustMgrDBExecutorFromContext(ctx), nil)
	if err != nil {
		return nil, status.CheckDatabaseError(err)
	}
	return res, err
}

func DeleteResource(ctx context.Context, resID types.SFID) error {
	return status.CheckDatabaseError((&models.Resource{
		RelResource: models.RelResource{ResourceID: resID},
	}).DeleteByResourceID(types.MustMgrDBExecutorFromContext(ctx)))
}

func GetBySFID(ctx context.Context, id types.SFID) (*models.Resource, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Resource{RelResource: models.RelResource{ResourceID: id}}

	if err := m.FetchByResourceID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.ResourceNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}
