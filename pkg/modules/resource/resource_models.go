package resource

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CondArgs struct {
	AccountID      types.SFID      `name:"-"`
	ResourceIDs    []types.SFID    `in:"query" name:"resourceID,omitempty"`
	UploadedBefore types.Timestamp `in:"query" name:"uploadedBefore,omitempty"`
	UploadedAfter  types.Timestamp `in:"query" name:"uploadedAfter,omitempty"`
	ExpireBefore   types.Timestamp `in:"query" name:"expireBefore,omitempty"`
	ExpireAfter    types.Timestamp `in:"query" name:"expireAfter,omitempty"`
	FilenameLike   string          `in:"query" name:"filenameLike,omitempty"`
	Filenames      []string        `in:"query" name:"filename,omitempty"`
	Md5            string          `in:"query" name:"md5,omitempty"`
}

func (r *CondArgs) Condition() builder.SqlCondition {
	var (
		own = &models.ResourceOwnership{}
		res = &models.Resource{}
		c   []builder.SqlCondition
	)
	if r.AccountID != 0 {
		c = append(c, own.ColAccountID().Eq(r.AccountID))
	}
	if len(r.ResourceIDs) > 0 {
		c = append(c, res.ColResourceID().In(r.ResourceIDs))
	}
	if !r.UploadedBefore.IsZero() {
		c = append(c, own.ColUploadedAt().Lte(r.UploadedBefore))
	}
	if !r.UploadedAfter.IsZero() {
		c = append(c, own.ColUploadedAt().Gte(r.UploadedAfter))
	}
	if !r.ExpireBefore.IsZero() {
		c = append(c, own.ColExpireAt().Lte(r.ExpireBefore))
	}
	if !r.UploadedAfter.IsZero() {
		c = append(c, own.ColExpireAt().Gte(r.ExpireAfter))
	}
	if len(r.Filenames) > 0 {
		c = append(c, own.ColFilename().In(r.Filenames))
	}
	if len(r.FilenameLike) > 0 {
		c = append(c, own.ColFilename().Like(r.FilenameLike))
	}
	if r.Md5 != "" {
		c = append(c, res.ColMd5().Eq(r.Md5))
	}
	return builder.And(c...)
}

type ListReq struct {
	CondArgs
	datatypes.Pager
}

type ResourceInfo struct {
	models.RelResource
	models.ResourceInfo
	models.ResourceOwnerInfo
	datatypes.OperationTimes
}

type ListRsp struct {
	Data  []*ResourceInfo `json:"data"`
	Total int64           `json:"total"`
}

type DownLoadResourceRsp struct {
	FileName string `json:"fileName"`
	Url      string `json:"url"`
}
