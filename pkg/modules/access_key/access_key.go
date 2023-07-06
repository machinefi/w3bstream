package access_key

import (
	"context"
	"time"

	"github.com/pkg/errors"

	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func init() {
	jwt.SetBuiltInTokenFn(Validate)
}

func Create(ctx context.Context, r *CreateReq) (*CreateRsp, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	acc := types.MustAccountFromContext(ctx)

	switch r.IdentityType {
	case enums.ACCESS_KEY_IDENTITY_TYPE__ACCOUNT, enums.ACCESS_KEY_IDENTITY_TYPE_UNKNOWN:
		r.IdentityID = acc.AccountID
		r.IdentityType = enums.ACCESS_KEY_IDENTITY_TYPE__ACCOUNT
	case enums.ACCESS_KEY_IDENTITY_TYPE__PUBLISHER:
	default:
		return nil, status.InvalidAccessKeyIdentityType
	}

	kctx := NewDefaultAccessKeyContext()

	exp := time.Time{}
	if r.ExpirationDays > 0 {
		exp = time.Now().UTC().Add(time.Hour * 24 * time.Duration(r.ExpirationDays))
	}

	m := &models.AccessKey{
		RelAccount: models.RelAccount{AccountID: acc.AccountID},
		AccessKeyInfo: models.AccessKeyInfo{
			IdentityID:   r.IdentityID,
			IdentityType: r.IdentityType,
			Name:         r.Name,
			Rand:         kctx.Rand,
			ExpiredAt:    types.Timestamp{Time: exp},
			Description:  r.Desc,
		},
		OperationTimesWithDeleted: datatypes.OperationTimesWithDeleted{
			OperationTimes: datatypes.OperationTimes{
				CreatedAt: base.Timestamp{Time: kctx.GenTS},
			},
		},
	}

	err := sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			for {
				err := m.FetchByRand(d)
				if err != nil {
					if sqlx.DBErr(err).IsNotFound() {
						return nil
					} else {
						return status.DatabaseError.StatusErr().WithDesc(err.Error())
					}
				} else {
					kctx.Regenerate()
					m.Rand = kctx.Rand
				}
			}
		},
		func(d sqlx.DBExecutor) error {
			if err := m.Create(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.AccessKeyNameConflict
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()
	if err != nil {
		return nil, err
	}

	key, _ := kctx.MarshalText()

	rsp := &CreateRsp{
		Name:         r.Name,
		IdentityType: r.IdentityType,
		IdentityID:   r.IdentityID,
		AccessKey:    string(key),
		ExpiredAt:    &types.Timestamp{Time: exp},
		Desc:         r.Desc,
	}
	if rsp.ExpiredAt.IsZero() {
		rsp.ExpiredAt = nil
	}
	return rsp, nil
}

func DeleteByName(ctx context.Context, name string) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	acc := types.MustAccountFromContext(ctx)

	m := &models.AccessKey{
		RelAccount:    models.RelAccount{AccountID: acc.AccountID},
		AccessKeyInfo: models.AccessKeyInfo{Name: name},
	}

	if err := m.DeleteByAccountIDAndName(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return status.AccessKeyNotFound
		}
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func List(ctx context.Context, r *ListReq) (*ListRsp, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.AccessKey{}

	cond := r.Condition()
	adds := r.Addition()

	lst, err := m.List(d, cond, adds)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	cnt, err := m.Count(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}

	ret := &ListRsp{Total: cnt}

	for i := range lst {
		v := &ListData{
			Name: lst[i].Name,
			Desc: lst[i].Description,
		}
		if !lst[i].ExpiredAt.IsZero() {
			v.ExpiredAt = &base.Timestamp{Time: lst[i].ExpiredAt.Time}
		}
		if !lst[i].LastUsed.IsZero() {
			v.LastUsed = &base.Timestamp{Time: lst[i].LastUsed.Time}
		}
		v.OperationTimes = lst[i].OperationTimes
		ret.Data = append(ret.Data, v)
	}

	return ret, nil
}

func Validate(ctx context.Context, key string) (interface{}, error, bool) {
	kctx := &AccessKeyContext{}

	err := kctx.UnmarshalText([]byte(key))
	if err != nil {
		return nil, err, err != ErrInvalidPrefixOrPartCount
	}

	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.AccessKey{
		AccessKeyInfo: models.AccessKeyInfo{
			Rand: kctx.Rand,
		},
	}
	if err = m.FetchByRand(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.AccessKeyNotFound, true
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error()), true
	}

	if kctx.GenTS.UTC().Second() != m.CreatedAt.UTC().Second() {
		return nil, status.InvalidAccessKey, true
	}

	if !m.ExpiredAt.IsZero() && time.Now().UTC().After(m.ExpiredAt.Time) {
		return nil, status.AccessKeyExpired, true
	}

	ts := base.Timestamp{Time: time.Now().UTC()}
	if _, err = d.Exec(
		builder.Update(d.T(m)).Set(
			m.ColUpdatedAt().ValueBy(ts),
			m.ColLastUsed().ValueBy(ts),
		).Where(
			m.ColRand().Eq(kctx.Rand),
		),
	); err != nil {
		conflog.FromContext(ctx).Warn(errors.Wrap(err, "update access key last used"))
	}

	// TODO check privileges
	return m.IdentityID, nil, true
}
