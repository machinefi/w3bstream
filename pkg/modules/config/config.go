package config

import (
	"context"
	"encoding/json"
	"fmt"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func GetConfigValue(ctx context.Context, rel types.SFID, v wasm.Configuration) error {
	l := types.MustLoggerFromContext(ctx).WithValues("rel", rel)

	_, l = l.Start(ctx, "GetConfigValue")
	defer l.End()

	typ := v.ConfigType()

	m, err := GetConfigByRelIdAndType(ctx, rel, typ)
	if err != nil {
		l.Error(err)
		return err
	}
	if err = json.Unmarshal(m.Value, v); err != nil {
		l.Error(err)
		return status.InternalServerError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func FetchConfigValuesByRelIDs(ctx context.Context, relIDs ...types.SFID) ([]wasm.Configuration, error) {
	l := types.MustLoggerFromContext(ctx)
	d := types.MustMgrDBExecutorFromContext(ctx)

	_, l = l.Start(ctx, "FetchConfigsByRelIDs")
	defer l.End()

	ms := make([]models.Config, 0)
	m := &models.Config{}
	err := d.QueryAndScan(
		builder.Select(nil).From(
			d.T(m),
			builder.Where(m.ColRelID().In(relIDs)),
		),
		&ms,
	)
	if err != nil {
		return nil, status.CheckDatabaseError(err)
	}

	configs := make([]wasm.Configuration, 0, len(ms))
	for _, cfg := range ms {
		v, err := wasm.NewConfigurationByType(cfg.Type)
		if err != nil {
			l.Error(err)
			continue
		}
		if err = json.Unmarshal(cfg.Value, v); err != nil {
			return nil, status.InternalServerError.StatusErr().WithDesc(err.Error())
		}
		configs = append(configs, v)
	}
	return configs, nil
}

func GetConfigByRelIdAndType(ctx context.Context, rel types.SFID, typ enums.ConfigType) (*models.Config, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "GetConfigByRelIdAndType")
	defer l.End()

	cfg := &models.Config{
		ConfigBase: models.ConfigBase{
			RelID: rel,
			Type:  typ,
		},
	}

	err := cfg.FetchByRelIDAndType(d)
	if err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.ConfigNotFound.StatusErr().
				WithDesc(fmt.Sprintf("rel:%v type: %v", rel, typ))
		} else {
			return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
		}
	}

	return cfg, nil
}

func CreateConfig(ctx context.Context, rel types.SFID, cfg wasm.Configuration) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "CreateConfig")
	defer l.End()

	raw, err := json.Marshal(cfg)
	if err != nil {
		l.Error(err)
		return err
	}

	m := &models.Config{
		RelConfig: models.RelConfig{ConfigID: idg.MustGenSFID()},
		ConfigBase: models.ConfigBase{
			RelID: rel,
			Type:  cfg.ConfigType(),
			Value: raw,
		},
	}
	if err = m.Create(d); err != nil {
		return err
	}
	return nil
}

func CreateOrUpdateConfig(ctx context.Context, rel types.SFID, typ enums.ConfigType, raw []byte) (cfg *models.Config, err error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
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
