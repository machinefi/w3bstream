package config

import (
	"context"
	"encoding/json"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func GetValue(ctx context.Context, id types.SFID, t enums.ConfigType) (wasm.Configuration, error) {
	m, err := GetByRelAndType(ctx, id, t)
	if err != nil {
		return nil, err
	}
	c, err := wasm.NewConfigurationByType(t)
	if err != nil {
		return nil, status.InvalidConfigType.StatusErr().WithDesc(err.Error())
	}
	if err = json.Unmarshal(m.Value, c); err != nil {
		return nil, status.InternalServerError.StatusErr().WithDesc(err.Error())
	}
	return c, nil
}

func GetByRelAndType(ctx context.Context, id types.SFID, t enums.ConfigType) (*models.Config, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Config{ConfigBase: models.ConfigBase{RelID: id, Type: t}}
	v := &Detail{id, t}

	if err := m.FetchByRelIDAndType(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.ConfigNotFound.StatusErr().WithDesc(v.Log(err))
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(v.Log(err))
	}
	return m, nil
}

func Upsert(ctx context.Context, id types.SFID, c wasm.Configuration) error {
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		m   *models.Config
		v   = &Detail{id, c}
		err error
		old wasm.Configuration
	)

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			m, err = GetByRelAndType(ctx, id, c.ConfigType())
			if err != nil {
				if se, ok := statusx.IsStatusErr(err); ok &&
					se.Code == status.ConfigNotFound.Code() {
					return nil
				}
				return err
			}
			if old, err = wasm.NewConfigurationByType(m.Type); err != nil {
				return status.ConfigInitFailed.StatusErr().WithDesc(v.Log(err))
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			if err = wasm.UninitConfiguration(ctx, old); err != nil {
				return status.ConfigUninitFailed.StatusErr().WithDesc(v.Log(err))
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if old == nil {
				return Create(ctx, id, c)
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if m != nil {
				if m.Value, err = json.Marshal(c); err != nil {
					return status.ConfigParsingFailed.StatusErr().WithDesc(v.Log(err))
				}
				if err = m.UpdateByConfigID(d); err != nil {
					if sqlx.DBErr(err).IsConflict() {
						return status.ConfigConflict.StatusErr().WithDesc(v.Log(err))
					}
					return status.DatabaseError.StatusErr().WithDesc(err.Error())
				}
			}
			return nil
		},
	).Do()
}

func List(ctx context.Context, r *DataListParam) ([]*Detail, error) {
	l := types.MustLoggerFromContext(ctx)
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Config{}

	lst, err := m.List(d, r.Condition())
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}

	configs := make([]*Detail, 0, len(lst))
	for _, cfg := range lst {
		c, err := wasm.NewConfigurationByType(cfg.Type)
		if err != nil {
			l.Warn(err)
			continue
		}
		if err = json.Unmarshal(cfg.Value, c); err != nil {
			return nil, status.InternalServerError.StatusErr().WithDesc(err.Error())
		}
		configs = append(configs, &Detail{
			RelID:         cfg.RelID,
			Configuration: c,
		})
	}
	return configs, nil
}

func BatchRemove(ctx context.Context, r *DataListParam) error {
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		m   = &models.Config{}
		lst []*Detail
		err error
	)

	sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			lst, err = List(ctx, r)
			return err
		},
		func(d sqlx.DBExecutor) error {
			if _, err = d.Exec(
				builder.Delete().From(d.T(m), builder.Where(r.Condition())),
			); err != nil {
				return status.InternalServerError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			summary := make(statusx.ErrorFields, 0, len(lst))
			for _, c := range lst {
				err2 := wasm.UninitConfiguration(ctx, c.Configuration)
				if err2 != nil {
					summary = append(summary, &statusx.ErrorField{
						Field: c.String(), Msg: err.Error(),
					})
				}
			}
			if len(summary) > 0 {
				return status.ConfigUninitFailed.StatusErr().
					AppendErrorFields(summary...)
			}
			return nil
		},
	)
	return nil
}

func Create(ctx context.Context, id types.SFID, c wasm.Configuration) error {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		v = &Detail{id, c}
	)

	raw, err := json.Marshal(c)
	if err != nil {
		return status.ConfigParsingFailed.StatusErr().WithDesc(err.Error())
	}

	m := &models.Config{
		RelConfig: models.RelConfig{
			ConfigID: confid.MustSFIDGeneratorFromContext(ctx).MustGenSFID(),
		},
		ConfigBase: models.ConfigBase{
			RelID: id, Type: c.ConfigType(), Value: raw,
		},
	}

	return sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			if err = m.Create(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.ConfigConflict.StatusErr().
						WithDesc(v.Log(err))
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			if err = wasm.InitConfiguration(ctx, c); err != nil {
				return status.ConfigInitFailed.StatusErr().
					WithDesc(v.Log(err))
			}
			return nil
		},
	).Do()
}
