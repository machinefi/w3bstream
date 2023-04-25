package resource

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/util"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func FetchOrCreateResource(ctx context.Context, owner string, f *multipart.FileHeader) (*models.Resource, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)
	fileSystemOp := types.MustFileSystemOpFromContext(ctx)

	_, l = l.Start(ctx, "FetchOrCreateResource")
	defer l.End()

	data, md5, err := getDataFromFileHeader(ctx, f)
	if err != nil {
		return nil, err
	}
	m := &models.Resource{ResourceInfo: models.ResourceInfo{Path: md5}}

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
					return status.DatabaseError.StatusErr().
						WithDesc(errors.Wrap(err, "FetchResource").Error())
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
					return status.DatabaseError.StatusErr().
						WithDesc(errors.Wrap(err, "UpdateResource").Error())
				}
				return nil
			} else {
				if err := fileSystemOp.Upload(md5, data); err != nil {
					return status.UploadFileFailed.StatusErr().WithDesc(err.Error())
				}

				m.ResourceID = idg.MustGenSFID()
				m.ResourceInfo.Path = md5
				m.ResourceInfo.RefCnt = 1
				if err := m.Create(db); err != nil {
					l.WithValues("stg", "CreateResource").Error(err)
					if sqlx.DBErr(err).IsConflict() {
						return status.ResourcePathConflict
					}
					return status.DatabaseError.StatusErr().
						WithDesc(errors.Wrap(err, "CreateResource").Error())
				}
				return nil
			}
		},
	).Do()

	l.Info("get wasm resource from db")
	return m, err
}

func getDataFromFileHeader(ctx context.Context, f *multipart.FileHeader) (data []byte, sum string, err error) {
	l := types.MustLoggerFromContext(ctx)
	uploadConf := types.MustUploadConfigFromContext(ctx)

	var (
		fr       io.ReadSeekCloser
		filesize = int64(0)
	)

	_, l = l.Start(ctx, "getDataFromFileHeader")
	defer l.End()

	if fr, err = f.Open(); err != nil {
		return
	}
	defer fr.Close()

	if filesize, err = fr.Seek(0, io.SeekEnd); err != nil {
		l.Error(err)
		return
	}
	if filesize > uploadConf.FileSizeLimit {
		err = errors.Wrap(err, "filesize over limit")
		l.Error(err)
		return
	}

	_, err = fr.Seek(0, io.SeekStart)
	if err != nil {
		l.Error(err)
		return
	}

	data = make([]byte, filesize)
	_, err = fr.Read(data)
	if err != nil {
		l.Error(err)
		return
	}

	sum, err = util.ByteMD5(data)
	return
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
