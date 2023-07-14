package access_key

import (
	"context"
	"time"

	"github.com/pkg/errors"

	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
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
			Privileges:   r.Privileges.ConvToPrivilegeModel(),
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
		Privileges:   ConvToGroupMetaWithPrivileges(m.Privileges),
		Desc:         r.Desc,
	}
	if !exp.IsZero() {
		rsp.ExpiredAt = &types.Timestamp{Time: exp}
	}
	return rsp, nil
}

func UpdateByName(ctx context.Context, name string, r *UpdateReq) (*UpdateRsp, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	acc := types.MustAccountFromContext(ctx)

	m := &models.AccessKey{
		RelAccount:    models.RelAccount{AccountID: acc.AccountID},
		AccessKeyInfo: models.AccessKeyInfo{Name: name},
	}

	err := sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			if err := m.FetchByAccountIDAndName(d); err != nil {
				if sqlx.DBErr(err).IsNotFound() {
					return status.AccessKeyNotFound
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if r.Desc != "" {
				m.Description = r.Desc
			}
			m.ExpiredAt = base.Timestamp{}
			if r.ExpirationDays > 0 {
				m.ExpiredAt = base.Timestamp{
					Time: time.Now().UTC().Add(time.Hour * 24 * time.Duration(r.ExpirationDays)),
				}
			}
			m.Privileges = r.Privileges.ConvToPrivilegeModel()
			_, err := d.Exec(
				builder.Update(d.T(m)).Set(
					m.ColDescription().ValueBy(m.Description),
					m.ColExpiredAt().ValueBy(m.ExpiredAt),
					m.ColPrivileges().ValueBy(m.Privileges),
				).Where(
					builder.And(
						m.ColDeletedAt().Eq(0),
						m.ColAccountID().Eq(acc.AccountID),
						m.ColName().Eq(name),
					),
				),
			)
			if err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()
	if err != nil {
		return nil, err
	}
	rsp := &UpdateRsp{
		Name:         m.Name,
		IdentityType: m.IdentityType,
		IdentityID:   m.IdentityID,
		Privileges:   ConvToGroupMetaWithPrivileges(m.Privileges),
		Desc:         m.Description,
	}
	if !m.ExpiredAt.IsZero() {
		rsp.ExpiredAt = &m.ExpiredAt
	}
	if !m.LastUsed.IsZero() {
		rsp.LastUsed = &m.LastUsed
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
		ret.Data = append(ret.Data, NewListDataByModel(&lst[i]))
	}

	return ret, nil
}

func GetByName(ctx context.Context, name string) (*ListData, error) {
	acc := types.MustAccountFromContext(ctx)

	k := &models.AccessKey{
		RelAccount:    models.RelAccount{AccountID: acc.AccountID},
		AccessKeyInfo: models.AccessKeyInfo{Name: name},
	}

	err := k.FetchByAccountIDAndName(types.MustMgrDBExecutorFromContext(ctx))
	if err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.AccessKeyNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return NewListDataByModel(k), nil
}

func Validate(ctx context.Context, key string) (interface{}, error, bool) {
	opId := httptransport.OperationIDFromContext(ctx)
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

	groupName, ok := gOperators[opId]
	if !ok {
		err = errors.Errorf("operator id is not registered: operator[%s]", opId)
		return nil, status.AccessKeyPermissionDenied.StatusErr().WithDesc(err.Error()), true
	}
	perm, ok := m.Privileges[groupName]
	if !ok {
		err = errors.Errorf("no group permission: group[%s]", groupName)
		return nil, status.AccessKeyPermissionDenied.StatusErr().WithDesc(err.Error()), true
	}
	opMeta := gOperatorGroups[groupName].Operators[opId]
	if perm < opMeta.MinimalPerm {
		err = errors.Errorf("no operator permission: operator[%s] group[%s] have perm[%s] need perm[%s]",
			opId, groupName, perm, opMeta.MinimalPerm)
		return nil, status.AccessKeyPermissionDenied.StatusErr().WithDesc(err.Error()), true
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

	return m, nil, true
}
