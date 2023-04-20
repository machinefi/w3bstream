package resource

import (
	"context"
	"mime/multipart"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func FetchOrCreateResource(ctx context.Context, owner string, f *multipart.FileHeader) (*models.Resource, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "FetchOrCreateResource")
	defer l.End()

	fullName, err := UploadWithS3(ctx, f, owner)
	if err != nil {
		l.Error(err)
		return nil, status.UploadFileFailed.StatusErr().WithDesc(err.Error())
	}

	m := &models.Resource{ResourceInfo: models.ResourceInfo{Path: fullName}}

	var exists bool
	err = sqlx.NewTasks(d).With(
		// fetch Resource
		func(db sqlx.DBExecutor) error {
			err := m.FetchByPath(d)
			if err != nil {
				if sqlx.DBErr(err).IsNotFound() {
					exists = false
					return nil
				} else {
					return status.CheckDatabaseError(err, "FetchResource")
				}
			} else {
				exists = true
				return nil
			}
		},
		// create or update Resource
		func(db sqlx.DBExecutor) error {
			if exists {
				m.ResourceInfo.RefCnt += 1
				if err := m.UpdateByPath(d); err != nil {
					return status.CheckDatabaseError(err, "UpdateResource")
				}
				return nil
			} else {
				m.ResourceID = idg.MustGenSFID()
				m.ResourceInfo.Path = fullName
				m.ResourceInfo.RefCnt = 1
				if err := m.Create(db); err != nil {
					return status.CheckDatabaseError(err, "CreateResource")
				}
				return nil
			}
		},
	).Do()

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
