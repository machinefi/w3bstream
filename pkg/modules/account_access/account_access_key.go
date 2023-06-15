package account_access

import (
	"bytes"
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func init() {
	jwt.SetBuiltInTokenFn(Validate)
}

func Create(ctx context.Context, r *CreateReq) (*models.AccountAccessKey, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	acc := types.MustAccountFromContext(ctx)

	rand, _, ts := GenAccessKey(acc.AccountID)

	exp := time.Time{}
	if r.ExpirationDays > 0 {
		exp = time.Now().UTC().Add(time.Hour * 24 * time.Duration(r.ExpirationDays))
	}

	m := &models.AccountAccessKey{
		RelAccount: models.RelAccount{AccountID: acc.AccountID},
		AccountAccessKeyInfo: models.AccountAccessKeyInfo{
			Name:      r.Name,
			AccessKey: rand,
			ExpiredAt: types.Timestamp{Time: exp},
		},
		OperationTimesWithDeleted: datatypes.OperationTimesWithDeleted{
			OperationTimes: datatypes.OperationTimes{
				CreatedAt: types.Timestamp{Time: ts},
			},
		},
	}

	if err := m.Create(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return nil, status.AccountKeyNameConflict
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func DeleteByName(ctx context.Context, name string) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	acc := types.MustAccountFromContext(ctx)

	m := &models.AccountAccessKey{
		RelAccount:           models.RelAccount{AccountID: acc.AccountID},
		AccountAccessKeyInfo: models.AccountAccessKeyInfo{Name: name},
	}

	if err := m.DeleteByAccountIDAndName(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return status.AccountKeyNotFound
		}
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func Validate(ctx context.Context, key string) (interface{}, error, bool) {
	if !strings.HasPrefix(key, "w3b_") {
		return nil, nil, false
	}
	id, rand, _, err := ParseAccessKey(key)
	if err != nil {
		return nil, err, true
	}

	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.AccountAccessKey{
		AccountAccessKeyInfo: models.AccountAccessKeyInfo{
			AccessKey: rand,
		},
	}
	if err = m.FetchByAccessKey(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.AccountKeyNotFound, true
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error()), true
	}

	if id != m.AccountID {
		return nil, status.InvalidAccountAccessKey, true
	}

	if time.Now().After(m.ExpiredAt.Time) {
		return nil, status.AccountAccessKeyExpired, true
	}

	return id, nil, true
}

// GenAccessKey key contains token owner, random string and generate time
func GenAccessKey(id types.SFID) (rand, key string, ts time.Time) {
	rand = uuid.New().String()
	ts = time.Now().UTC()
	key = "w3b_" + base64.StdEncoding.EncodeToString([]byte(
		id.String()+"_"+rand+"_"+ts.Format(time.RFC3339Nano),
	))
	return
}

var (
	ErrMsgAccessKeyInvalidPartCountOrPrefix = errors.New("invalid part count or prefix")
	ErrMsgAccessKeyBase64Decode             = errors.New("base64 decode")
	ErrMsgAccessKeyInvalidPartCount         = errors.New("invalid part count of contents")
	ErrMsgAccessKeyInvalidAccountID         = errors.New("invalid account id")
	ErrMsgAccessKeyInvalidTimestamp         = errors.New("invalid timestamp")
)

// ParseAccessKey parse access key
func ParseAccessKey(key string) (id types.SFID, rand string, ts time.Time, err error) {
	defer func() {
		if err != nil {
			err = status.InvalidAccountAccessKey.StatusErr().WithDesc(err.Error())
		}
	}()

	parts := strings.Split(key, "_")
	if len(parts) != 2 || parts[0] != "w3b" {
		err = ErrMsgAccessKeyInvalidPartCountOrPrefix
		return
	}
	var raw []byte
	raw, err = base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		err = errors.Wrap(err, ErrMsgAccessKeyBase64Decode.Error())
		return
	}

	contents := bytes.Split(raw, []byte("_"))
	if len(contents) != 3 {
		err = ErrMsgAccessKeyInvalidPartCount
		return
	}

	if err = id.UnmarshalText(contents[0]); err != nil {
		err = errors.Wrap(err, ErrMsgAccessKeyInvalidAccountID.Error())
		return
	}
	rand = string(contents[1])

	ts, err = time.ParseInLocation(time.RFC3339Nano, string(contents[2]), time.Local)
	if err != nil {
		err = errors.Wrap(err, ErrMsgAccessKeyInvalidTimestamp.Error())
		return
	}
	return
}
