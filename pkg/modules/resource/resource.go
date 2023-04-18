package resource

import (
	"context"
	"mime/multipart"
	"os"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

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

func RemoveBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Resource{RelResource: models.RelResource{ResourceID: id}}

	if err := m.DeleteByResourceID(d); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil

}

func Create(ctx context.Context, acc types.SFID, md5 string, f *multipart.FileHeader) (*models.Resource, error) {
	ff, err := f.Open()
	if err != nil {
		return nil, status.UploadFileFailed.StatusErr().WithDesc(err.Error())
	}
	defer ff.Close()

	fid := confid.MustSFIDGeneratorFromContext(ctx).MustGenSFID()

	filename, sum, err := Upload(ctx, fid.String(), ff)
	if err != nil {
		return nil, status.UploadFileFailed.StatusErr().WithDesc(err.Error())
	}

	if md5 != "" && sum != md5 {
		_ = os.Remove(filename)
		return nil, status.MD5ChecksumFailed
	}

	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Resource{ResourceInfo: models.ResourceInfo{Md5: md5}}
	exists := false

	err = sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			if err = m.FetchByMd5(d); err != nil {
				if sqlx.DBErr(err).IsNotFound() {
					exists = false
					return nil
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			exists = true
			_ = os.Remove(filename)
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if exists {
				return nil
			}
			// TODO m.Owner = acc
			m.ResourceID = fid
			m.ResourceInfo.Path = filename
			m.ResourceInfo.Md5 = md5
			if err = m.Create(d); err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()
	if err != nil {
		return nil, err
	}
	return m, nil
}
