package resource

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/util"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func FetchOrCreateResource(ctx context.Context, accountID, appletID types.SFID, fileName string, f *multipart.FileHeader) (*models.Resource, error) {
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
	mMeta := &models.ResourceMeta{}

	var resExists, metaExists bool
	err = sqlx.NewTasks(d).With(
		// fetch Resource
		func(db sqlx.DBExecutor) error {
			err := m.FetchByPath(db)
			if err != nil {
				if sqlx.DBErr(err).IsNotFound() {
					resExists = false
					return nil
				} else {
					return status.DatabaseError.StatusErr().
						WithDesc(errors.Wrap(err, "FetchResource").Error())
				}
			} else {
				resExists = true
				return nil
			}
		},
		// create or update Resource
		func(db sqlx.DBExecutor) error {
			if !resExists {
				if err := fileSystemOp.Upload(md5, data); err != nil {
					return status.UploadFileFailed.StatusErr().WithDesc(err.Error())
				}

				m.ResourceID = idg.MustGenSFID()
				m.ResourceInfo.Path = md5
				if err := m.Create(db); err != nil {
					l.WithValues("stg", "CreateResource").Error(err)
					if sqlx.DBErr(err).IsConflict() {
						return status.ResourcePathConflict
					}
					return status.DatabaseError.StatusErr().
						WithDesc(errors.Wrap(err, "CreateResource").Error())
				}
			}
			return nil
		},

		// fetch resource meta info
		func(db sqlx.DBExecutor) error {
			mMeta.ResourceID = m.ResourceID
			mMeta.AccountID = accountID
			mMeta.AppletID = appletID
			err := mMeta.FetchByResourceIDAndAccountIDAndAppletID(db)
			if err != nil {
				if sqlx.DBErr(err).IsNotFound() {
					metaExists = false
					return nil
				} else {
					return status.DatabaseError.StatusErr().
						WithDesc(errors.Wrap(err, "FetchResourceMeta").Error())
				}
			} else {
				metaExists = true
				return nil
			}
		},
		// create or update resource meta info
		func(db sqlx.DBExecutor) error {
			if metaExists {
				mMeta.MetaInfo.FileName = fileName
				if err := mMeta.UpdateByMetaID(db); err != nil {
					return status.DatabaseError.StatusErr().
						WithDesc(errors.Wrap(err, "UpdateResourceMeta").Error())
				}
				return nil
			} else {
				mMeta.MetaID = idg.MustGenSFID()
				mMeta.MetaInfo.FileName = fileName
				if err := mMeta.Create(db); err != nil {
					l.WithValues("stg", "CreateResourceMeta").Error(err)
					if sqlx.DBErr(err).IsConflict() {
						return status.ResourceAccountConflict
					}
					return status.DatabaseError.StatusErr().
						WithDesc(errors.Wrap(err, "CreateResourceMeta").Error())
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

func ListResourceMeta(ctx context.Context, r *ListReq) (*ListRsp, error) {
	var (
		d    = types.MustMgrDBExecutorFromContext(ctx)
		m    = &models.ResourceMeta{}
		ret  = &ListRsp{}
		cond = r.Condition()
		adds = r.Additions()

		err error
	)

	ret.Data, err = m.List(d, cond, adds...)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	ret.Total, err = m.Count(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}

func ListResourceMetaDetail(ctx context.Context, r *ListReq) (*ListDetailRsp, error) {
	var (
		d     = types.MustMgrDBExecutorFromContext(ctx)
		mMeta = &models.ResourceMeta{}
		mPrj  = &models.Project{}
		mApp  = &models.Applet{}
		ret   = &ListDetailRsp{}
		cond  = r.Condition()
	)

	expr := builder.Select(builder.MultiWith(",",
		builder.Alias(mMeta.ColAccountID(), "f_acc_id"),
		builder.Alias(mMeta.ColAccountID(), "f_res_id"),
		builder.Alias(mPrj.ColProjectID(), "f_prj_id"),
		builder.Alias(mPrj.ColName(), "f_prj_name"),
		builder.Alias(mApp.ColAppletID(), "f_app_id"),
		builder.Alias(mApp.ColName(), "f_app_name"),
		builder.Alias(mMeta.ColFileName(), "f_file_name"),
		builder.Alias(mMeta.ColUpdatedAt(), "f_updated_at"),
		builder.Alias(mMeta.ColCreatedAt(), "f_created_at"),
	)).From(
		d.T(mMeta),
		append([]builder.Addition{
			builder.LeftJoin(d.T(mApp)).On(mMeta.ColResourceID().Eq(mApp.ColResourceID())),
			builder.LeftJoin(d.T(mPrj)).On(mApp.ColProjectID().Eq(mPrj.ColProjectID())),
			builder.Where(cond),
		}, r.Addition())...,
	)
	err := d.QueryAndScan(expr, &ret.Data)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	if ret.Total, err = mMeta.Count(d, cond); err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
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

func RemoveBySFID(ctx context.Context, id types.SFID) error {
	m := &models.ResourceMeta{RelMeta: models.RelMeta{MetaID: id}}

	if err := m.DeleteByMetaID(types.MustMgrDBExecutorFromContext(ctx)); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func RemoveByResIDAndAccIDAndAppID(ctx context.Context, resID, accID, appID types.SFID) error {
	m := &models.ResourceMeta{
		RelResource: models.RelResource{ResourceID: resID},
		RelAccount:  models.RelAccount{AccountID: accID},
		RelApplet:   models.RelApplet{AppletID: appID},
	}

	if err := m.DeleteByResourceIDAndAccountIDAndAppletID(types.MustMgrDBExecutorFromContext(ctx)); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}
func Remove(ctx context.Context, r *CondArgs) error {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.ResourceMeta{}
	)

	_, err := d.Exec(builder.Delete().From(
		d.T(m),
		builder.Where(r.Condition()),
	))
	if err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}
