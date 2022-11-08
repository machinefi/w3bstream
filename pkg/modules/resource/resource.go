package resource

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Info struct {
	models.Resource
	ResourceName string
}

func FetchOrCreateResource(ctx context.Context, f *multipart.FileHeader) (*Info, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "FetchOrCreateResource")
	defer l.End()

	_, fullName, fileName, md5, err := Upload(ctx, f, idg.MustGenSFID().String())
	if err != nil {
		l.Error(err)
		return nil, status.UploadFileFailed.StatusErr().WithDesc(err.Error())
	}

	m := &models.Resource{ResourceInfo: models.ResourceInfo{Md5: md5}}

	if err = m.FetchByMd5(d); err != nil && sqlx.DBErr(err).IsNotFound() {
		l.Error(errors.Wrap(err, fmt.Sprintf("fetch wasm resource by md5 - %s, maybe it doesnt exist.", md5)))
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
	return &Info{Resource: *m, ResourceName: fileName}, err
}
