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

func Marshal(c wasm.Configuration) (data []byte, err error) {
	data, err = json.Marshal(c)
	if err != nil {
		err = status.ConfigParseFailed.StatusErr().WithDesc(err.Error())
	}
	return
}

func Unmarshal(data []byte, typ enums.ConfigType) (c wasm.Configuration, err error) {
	c, err = wasm.NewConfigurationByType(typ)
	if err != nil {
		return nil, status.InvalidConfigType.StatusErr().WithDesc(err.Error())
	}
	if err = json.Unmarshal(data, c); err != nil {
		return nil, status.ConfigParseFailed.StatusErr().WithDesc(err.Error())
	}
	return c, nil
}

func GetBySFID(ctx context.Context, id types.SFID) (*models.Config, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Config{RelConfig: models.RelConfig{ConfigID: id}}

	if err := m.FetchByConfigID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.ConfigNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func GetValueBySFID(ctx context.Context, id types.SFID) (wasm.Configuration, error) {
	m, err := GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	return Unmarshal(m.Value, m.Type)
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

func GetValueByRelAndType(ctx context.Context, rel types.SFID, t enums.ConfigType) (wasm.Configuration, error) {
	m, err := GetByRelAndType(ctx, rel, t)
	if err != nil {
		return nil, err
	}
	return Unmarshal(m.Value, t)
}

func Upsert(ctx context.Context, rel types.SFID, c wasm.Configuration) (*models.Config, error) {
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		idg = confid.MustSFIDGeneratorFromContext(ctx)
		v   = &Detail{rel, c}
		err error
		m   *models.Config
		old wasm.Configuration
	)

	err = sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			m = &models.Config{
				ConfigBase: models.ConfigBase{RelID: rel, Type: c.ConfigType()},
			}
			if err = m.FetchByRelIDAndType(d); err != nil {
				if sqlx.DBErr(err).IsNotFound() {
					return nil
				}
				return status.DatabaseError.StatusErr().WithDesc(v.Log(err))
			}
			if old, err = Unmarshal(m.Value, c.ConfigType()); err != nil {
				return err
			}
			if err = wasm.UninitConfiguration(ctx, old); err != nil {
				return status.ConfigUninitFailed.StatusErr().WithDesc(v.Log(err))
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			var raw []byte
			raw, err = Marshal(c)
			if err != nil {
				return err
			}
			if old == nil {
				m = &models.Config{
					RelConfig: models.RelConfig{ConfigID: idg.MustGenSFID()},
					ConfigBase: models.ConfigBase{
						Type: c.ConfigType(), RelID: rel, Value: raw,
					},
				}
				err = m.Create(d)
			} else {
				m.Value = raw
				err = m.UpdateByRelIDAndType(d)
			}
			if err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.ConfigConflict.StatusErr().WithDesc(v.Log(err))
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if err = wasm.InitConfiguration(ctx, c); err != nil {
				return status.ConfigInitFailed.StatusErr().WithDesc(v.Log(err))
			}
			return nil
		},
	).Do()

	if err != nil {
		return nil, err
	}
	return m, nil
}

func List(ctx context.Context, r *CondArgs) ([]*Detail, error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.Config{}
		c wasm.Configuration
	)

	lst, err := m.List(d, r.Condition())
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}

	configs := make([]*Detail, 0, len(lst))
	for _, cfg := range lst {
		c, err = Unmarshal(cfg.Value, cfg.Type)
		if err != nil {
			return nil, err
		}
		configs = append(configs, &Detail{RelID: cfg.RelID, Configuration: c})
	}
	return configs, nil
}

func Remove(ctx context.Context, r *CondArgs) error {
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		m   = &models.Config{}
		lst []*Detail
		err error
	)

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, d)
			lst, err = List(ctx, r)
			return err
		},
		func(d sqlx.DBExecutor) error {
			if _, err = d.Exec(
				builder.Delete().From(d.T(m), builder.Where(r.Condition())),
			); err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			summary := make(statusx.ErrorFields, 0, len(lst))
			for _, c := range lst {
				err2 := wasm.UninitConfiguration(ctx, c.Configuration)
				if err2 != nil {
					summary = append(summary, &statusx.ErrorField{
						Field: c.String(), Msg: err2.Error(),
					})
				}
			}
			if len(summary) > 0 {
				return status.ConfigUninitFailed.StatusErr().
					AppendErrorFields(summary...)
			}
			return nil
		},
	).Do()
}

func Create(ctx context.Context, id types.SFID, c wasm.Configuration) (*models.Config, error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		v = &Detail{id, c}
	)

	raw, err := Marshal(c)
	if err != nil {
		return nil, err
	}

	m := &models.Config{
		RelConfig: models.RelConfig{
			ConfigID: confid.MustSFIDGeneratorFromContext(ctx).MustGenSFID(),
		},
		ConfigBase: models.ConfigBase{
			RelID: id, Type: c.ConfigType(), Value: raw,
		},
	}

	err = sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			if err = m.Create(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.ConfigConflict.StatusErr().WithDesc(v.Log(err))
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
	if err != nil {
		return nil, err
	}
	return m, nil
}
