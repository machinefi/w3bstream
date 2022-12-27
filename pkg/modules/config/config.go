package config

import (
	"context"
	"encoding/json"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Configuration interface {
	RelID() types.SFID
	Type() enums.ConfigType
	Value() interface{}
}

func GetConfigValue(ctx context.Context, rel types.SFID, typ enums.ConfigType, v interface{}) error {
	l := types.MustLoggerFromContext(ctx).WithValues("rel", rel)

	_, l = l.Start(ctx, "GetConfigValue")
	defer l.End()

	cfg, err := GetConfigByRelIdAndType(ctx, rel, typ)
	if err != nil {
		l.Error(err)
		return status.CheckDatabaseError(err)
	}
	if err = json.Unmarshal(cfg.Value, v); err != nil {
		l.Error(err)
		return status.InternalServerError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func GetConfigByRelIdAndType(ctx context.Context, rel types.SFID, typ enums.ConfigType) (*models.Config, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "GetConfigByRelIdAndType")
	defer l.End()

	cfg := &models.Config{
		ConfigBase: models.ConfigBase{
			RelID: rel,
			Type:  typ,
		},
	}

	if err := cfg.FetchByRelIDAndType(d); err != nil {
		return nil, err
	}

	return cfg, nil
}

func CreateConfig(ctx context.Context, rel types.SFID, typ enums.ConfigType, raw []byte) (*models.Config, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "CreateConfig")
	defer l.End()

	cfg := &models.Config{
		RelConfig: models.RelConfig{ConfigID: idg.MustGenSFID()},
		ConfigBase: models.ConfigBase{
			RelID: rel,
			Type:  typ,
			Value: raw,
		},
	}
	if err := cfg.Create(d); err != nil {
		return nil, err
	}
	return cfg, nil
}

func CreateOrUpdateConfig(ctx context.Context, rel types.SFID, typ enums.ConfigType, raw []byte) (cfg *models.Config, err error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "CreateOrUpdateConfig")
	defer l.End()

	cfg = &models.Config{
		ConfigBase: models.ConfigBase{
			RelID: rel,
			Type:  typ,
		},
	}

	found := false

	err = sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			err = cfg.FetchByRelIDAndType(db)
			if err == nil {
				found = true
				return err // do update
			}
			if err != nil && sqlx.DBErr(err).IsNotFound() {
				found = false
				return nil
			}
			return err
		},
		func(db sqlx.DBExecutor) error {
			if !found {
				return nil
			}
			cfg.Value = raw
			return cfg.UpdateByRelIDAndType(db)
		},
		func(db sqlx.DBExecutor) error {
			if found {
				return nil
			}
			cfg.ConfigID, cfg.Value = idg.MustGenSFID(), raw
			return cfg.Create(db)
		},
	).Do()
	if err != nil {
		return nil, status.CheckDatabaseError(err)
	}
	return
}
