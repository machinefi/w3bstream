package config

import (
	"context"
	"encoding/json"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
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
		return status.CheckDatabaseError(err)
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

	if err := cfg.FetchByRelIDAndType(d); err != nil {
		return nil, err
	}

	return cfg, nil
}

func CreateConfig(ctx context.Context, rel types.SFID, cfg wasm.Configuration) (*models.Config, error) {
	_, l := conflog.FromContext(ctx).Start(ctx, "CreateConfig")
	defer l.End()

	if err := wasm.Init(ctx, cfg); err != nil {
		l.Error(err)
		return nil, status.ConfigInitFailed.StatusErr().WithDesc(err.Error())
	}

	d := types.MustMgrDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	raw, err := json.Marshal(cfg)
	if err != nil {
		l.Error(err)
		return nil, status.InvalidConfig.StatusErr().WithDesc(err.Error())
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
		return nil, err
	}
	return m, nil
}

func CreateOrUpdateConfig(ctx context.Context, rel types.SFID, cfg wasm.Configuration) (*models.Config, error) {
	_, l := conflog.FromContext(ctx).Start(ctx, "CreateOrUpdateConfig")
	defer l.End()

	if err := wasm.Init(ctx, cfg); err != nil {
		l.Error(err)
		return nil, status.ConfigInitFailed.StatusErr().WithDesc(err.Error())
	}

	d := types.MustMgrDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	raw, err := json.Marshal(cfg)
	if err != nil {
		l.Error(err)
		return nil, status.InvalidConfig.StatusErr().WithDesc(err.Error())
	}

	m := &models.Config{
		ConfigBase: models.ConfigBase{
			RelID: rel,
			Type:  cfg.ConfigType(),
		},
	}

	found := false

	err = sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			err = m.FetchByRelIDAndType(db)
			if err == nil {
				found = true
				return err
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
			m.Value = raw
			return m.UpdateByRelIDAndType(db)
		},
		func(db sqlx.DBExecutor) error {
			if found {
				return nil
			}
			m.ConfigID, m.Value = idg.MustGenSFID(), raw
			return m.Create(db)
		},
	).Do()
	if err != nil {
		return nil, status.CheckDatabaseError(err)
	}
	return m, nil
}
