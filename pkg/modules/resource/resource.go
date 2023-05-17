package resource

import (
	"context"
	"mime/multipart"
	"time"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func Create(ctx context.Context, acc types.SFID, fh *multipart.FileHeader, filename, md5 string) (*models.Resource, []byte, error) {
	data, sum, err := CheckFileMd5SumAndGetData(ctx, fh, md5)
	if err != nil {
		return nil, nil, err
	}

	id := confid.MustNewSFIDGenerator().MustGenSFID()
	res := &models.Resource{}
	found := false

	err = sqlx.NewTasks(types.MustMgrDBExecutorFromContext(ctx)).With(
		func(d sqlx.DBExecutor) error {
			res.Md5 = sum
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
			path, err := UploadFile(ctx, data, id)
			if err != nil {
				return err
			}
			res = &models.Resource{
				RelResource:  models.RelResource{ResourceID: id},
				ResourceInfo: models.ResourceInfo{Path: path, Md5: sum},
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
			}
			err = own.FetchByResourceIDAndAccountID(d)
			if err != nil {
				if sqlx.DBErr(err).IsNotFound() {
					own.UploadedAt = types.Timestamp{Time: time.Now()}
					own.Filename = filename
					if err := own.Create(d); err != nil {
						return status.DatabaseError.StatusErr().WithDesc(err.Error())
					}
					return nil
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			} else {
				own.Filename = filename
				if err = own.UpdateByResourceIDAndAccountID(d); err != nil {
					return status.DatabaseError.StatusErr().WithDesc(err.Error())
				}
				return nil
			}
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
	data, err := ReadContent(ctx, res)
	if err != nil {
		return nil, nil, err
	}
	return res, data, nil
}

func GetContentByMd5(ctx context.Context, md5 string) (*models.Resource, []byte, error) {
	res, err := GetByMd5(ctx, md5)
	if err != nil {
		return nil, nil, err
	}
	data, err := ReadContent(ctx, res)
	if err != nil {
		return nil, nil, err
	}
	return res, data, nil
}

func ReadContent(ctx context.Context, m *models.Resource) ([]byte, error) {
	fs := types.MustFileSystemOpFromContext(ctx)
	data, err := fs.Read(m.Path)
	if err != nil {
		return nil, status.FetchResourceFailed.StatusErr().WithDesc(err.Error())
	}
	return data, nil
}

func GetOwnerByAccountAndSFID(ctx context.Context, acc, res types.SFID) (*models.ResourceOwnership, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.ResourceOwnership{
		RelAccount:  models.RelAccount{AccountID: acc},
		RelResource: models.RelResource{ResourceID: res},
	}

	if err := m.FetchByResourceIDAndAccountID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.ResourcePermNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func List(ctx context.Context, r *ListReq) (*ListRsp, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	res := &models.Resource{}
	own := &models.ResourceOwnership{}
	rsp := &ListRsp{}

	err := d.QueryAndScan(
		builder.Select(
			builder.MultiWith(",",
				builder.Alias(res.ColResourceID(), "f_resource_id"),
				builder.Alias(res.ColMd5(), "f_md5"),
				builder.Alias(own.ColUploadedAt(), "f_uploaded_at"),
				builder.Alias(own.ColExpireAt(), "f_expire_at"),
				builder.Alias(own.ColFilename(), "f_filename"),
				builder.Alias(own.ColComment(), "f_comment"),
				builder.Alias(own.ColCreatedAt(), "f_created_at"),
				builder.Alias(own.ColUpdatedAt(), "f_updated_at"),
			),
		).From(
			d.T(res),
			builder.LeftJoin(d.T(own)).On(res.ColResourceID().Eq(own.ColResourceID())),
			builder.Where(r.Condition()),
		), &rsp.Data)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}

	err = d.QueryAndScan(builder.Select(builder.Count()).From(
		d.T(res),
		builder.LeftJoin(d.T(own)).On(res.ColResourceID().Eq(own.ResourceID)),
		builder.Where(r.Condition()),
	), &rsp.Total)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return rsp, nil
}

func RemoveOwnershipBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	acc := types.MustAccountFromContext(ctx)

	m := &models.ResourceOwnership{
		RelResource: models.RelResource{ResourceID: id},
		RelAccount:  acc.RelAccount,
	}
	if err := m.DeleteByResourceIDAndAccountID(d); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}
