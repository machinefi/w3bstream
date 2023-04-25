package resource

import (
	"context"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func Create(ctx context.Context, acc types.SFID, fh *multipart.FileHeader, md5 string) (*models.Resource, []byte, error) {
	f, err := fh.Open()
	if err != nil {
		err = status.UploadFileFailed.StatusErr().WithDesc(err.Error())
		return nil, nil, err
	}

	path, data, err := UploadFile(ctx, f, md5)
	if err != nil {
		err = status.UploadFileFailed.StatusErr().WithDesc(err.Error())
		return nil, nil, err
	}

	id := confid.MustNewSFIDGenerator().MustGenSFID()
	res := &models.Resource{}
	found := false

	err = sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(d sqlx.DBExecutor) error {
			res.Md5 = md5
			if err = res.FetchByMd5(d); err != nil {
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
			res = &models.Resource{
				RelResource:  models.RelResource{ResourceID: id},
				ResourceInfo: models.ResourceInfo{Path: path, Md5: md5},
			}
			if err = res.Create(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.ResourceConflict
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			own := &models.ResourceOwnership{
				RelResource: models.RelResource{ResourceID: res.ResourceID},
				RelAccount:  models.RelAccount{AccountID: acc},
				ResourceOwnerInfo: models.ResourceOwnerInfo{
					UploadedAt: types.Timestamp{Time: time.Now()},
				},
			}
			if err = own.Create(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.ResourceOwnerConflict
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()

	if err != nil {
		return nil, nil, err
	}
	return res, data, nil
}

func GetBySFID(ctx context.Context, id types.SFID) (*models.Resource, error) {
	res := &models.Resource{}
	res.ResourceID = id
	if err := res.FetchByResourceID(types.MustMgrDBExecutorFromContext(ctx)); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.ResourceNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return res, nil
}

func GetByMd5(ctx context.Context, md5 string) (*models.Resource, error) {
	res := &models.Resource{}
	res.Md5 = md5
	if err := res.FetchByMd5(types.MustMgrDBExecutorFromContext(ctx)); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.ResourceNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return res, nil
}

func GetContentBySFID(ctx context.Context, id types.SFID) (*models.Resource, []byte, error) {
	res, err := GetBySFID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	if fs, _ := types.FileSystemFromContext(ctx); fs != nil {
		data, err := fs.FileCli.Read(res.Md5)
		if err != nil {
			return nil, nil, status.LocalResReadFailed.StatusErr().WithDesc(err.Error())
		}
		return res, data, nil
	}

	// try local filesystem
	c := types.MustUploadConfigFromContext(ctx)
	data, err := os.ReadFile(filepath.Join(c.Root, res.Path))
	if err != nil {
		return nil, nil, status.LocalResReadFailed.StatusErr().WithDesc(err.Error())
	}
	return res, data, nil
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
