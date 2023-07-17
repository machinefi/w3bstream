package access_key

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateReqBase struct {
	// Name access token name
	Name string `json:"name"`
	// ExpirationDays access token valid in ExpirationDays, if 0 means token will not be expired.
	ExpirationDays int `json:"expirationDays,omitempty"`
	// Description access token description
	Desc string `json:"desc,omitempty"`
	// Privileges operator group access privileges
	Privileges GroupAccessPrivileges `json:"privileges,omitempty"`
}

type CreateAccountAccessKeyReq = CreateReqBase

type CreateReq struct {
	// IdentityID associated with a publisher, an account or other application
	IdentityID types.SFID `json:"identityID,omitempty"`
	// IdentityType associated type, default associated current account
	IdentityType enums.AccessKeyIdentityType `json:"identityType,default='1'"`
	CreateReqBase
}

type CreateRsp struct {
	Name         string                      `json:"name"`
	IdentityType enums.AccessKeyIdentityType `json:"identityType"`
	IdentityID   types.SFID                  `json:"identityID"`
	AccessKey    string                      `json:"accessKey"`
	ExpiredAt    *types.Timestamp            `json:"expiredAt,omitempty"`
	LastUsed     *types.Timestamp            `json:"lastUsed,omitempty"`
	Desc         string                      `json:"desc,omitempty"`
}

type UpdateReq struct {
	ExpirationDays int                   `json:"expirationDays,omitempty"`
	Desc           string                `json:"desc,omitempty"`
	Privileges     GroupAccessPrivileges `json:"privileges,omitempty"`
}

type ListData struct {
	Name      string           `json:"name"`
	ExpiredAt *types.Timestamp `json:"expiredAt,omitempty"`
	LastUsed  *types.Timestamp `json:"lastUsed,omitempty"`
	Desc      string           `json:"desc,omitempty"`
	datatypes.OperationTimes
}

type CondArgs struct {
	AccountID      types.SFID                    `name:"-"`
	Names          []string                      `in:"query" name:"name,omitempty"`
	ExpiredAtBegin types.Timestamp               `in:"query" name:"expiredAtBegin,omitempty"`
	ExpiredAtEnd   types.Timestamp               `in:"query" name:"expiredAtEnd,omitempty"`
	IdentityIDs    types.SFIDs                   `in:"query" name:"identityID,omitempty"`
	IdentityTypes  []enums.AccessKeyIdentityType `in:"query" name:"identityType,omitempty"`
}

func (r *CondArgs) Condition() builder.SqlCondition {
	var (
		m = &models.AccessKey{}
		c []builder.SqlCondition
	)
	if r.AccountID != 0 {
		c = append(c, m.ColAccountID().Eq(r.AccountID))
	}
	if len(r.Names) > 0 {
		c = append(c, m.ColName().In(r.Names))
	}
	if !r.ExpiredAtEnd.IsZero() {
		c = append(c, m.ColExpiredAt().Lte(r.ExpiredAtEnd))
	}
	if !r.ExpiredAtBegin.IsZero() {
		c = append(c, m.ColExpiredAt().Gte(r.ExpiredAtBegin))
	}
	if len(r.IdentityIDs) > 0 {
		c = append(c, m.ColIdentityID().In(r.IdentityIDs))
	}
	if len(r.IdentityTypes) > 0 {
		c = append(c, m.ColIdentityType().In(r.IdentityTypes))
	}
	return builder.And(c...)
}

type ListReq struct {
	CondArgs
	datatypes.Pager
}

type ListRsp struct {
	Data  []*ListData `json:"data"`
	Total int64       `json:"total"`
}

var (
	gAccessKeyVersion    = 1
	gAccessKeyRandLength = 12
	gBase64Encoding      = base64.RawURLEncoding
)

func NewAccessKeyContext(version int) *AccessKeyContext {
	return &AccessKeyContext{
		Version: version,
		Rand:    stringsx.GenRandomVisibleString(gAccessKeyRandLength, '_'),
		GenTS:   time.Now(),
	}
}

func NewDefaultAccessKeyContext() *AccessKeyContext {
	return NewAccessKeyContext(gAccessKeyVersion)
}

type AccessKeyContext struct {
	// Version integer version fixed field
	Version int
	// GenTS access key generated timestamp (utc,seconds)
	GenTS time.Time
	// Rand random part, length: gAccessKeyRandLength
	Rand string
}

func (c *AccessKeyContext) Regenerate() {
	c.Rand = stringsx.GenRandomVisibleString(gAccessKeyRandLength, '_')
}

func (c AccessKeyContext) MarshalText() ([]byte, error) {
	return []byte("w3b_" + gBase64Encoding.EncodeToString([]byte(
		fmt.Sprintf("%d_%d_%s", c.Version, c.GenTS.UTC().Unix(), c.Rand),
	))), nil
}

func (c *AccessKeyContext) UnmarshalText(data []byte) (err error) {
	sep := []byte{'_'}

	parts := bytes.SplitN(data, sep, 2)
	if !bytes.Equal(parts[0], []byte("w3b")) {
		return ErrInvalidPrefixOrPartCount
	}
	var raw []byte
	raw, err = gBase64Encoding.DecodeString(string(parts[1]))
	if err != nil {
		err = errors.Wrap(ErrBase64DecodeFailed, err.Error())
		return
	}

	contents := bytes.SplitN(raw, sep, 3)
	if len(contents) != 3 {
		return ErrInvalidContentsPartCount
	}

	c.Version, err = strconv.Atoi(string(contents[0]))
	if err != nil {
		return errors.Wrap(ErrParseVersionFailed, err.Error())
	}
	genTS, err := strconv.ParseInt(string(contents[1]), 10, 64)
	if err != nil {
		return errors.Wrap(ErrParseGenTsFailed, err.Error())
	}
	c.GenTS = time.Unix(genTS, 0)
	switch c.Version {
	case 1:
		c.Rand = string(contents[2])
		return
	default:
		return ErrInvalidVersion
	}
}

func (c *AccessKeyContext) Equal(v *AccessKeyContext) bool {
	return c.Version == v.Version &&
		c.GenTS.UTC().Second() == v.GenTS.UTC().Second() &&
		c.Rand == v.Rand
}

var (
	ErrInvalidPrefixOrPartCount = errors.New("invalid prefix or part count")
	ErrBase64DecodeFailed       = errors.New("base64 decode failed")
	ErrInvalidContentsPartCount = errors.New("invalid part count of contents")
	ErrParseVersionFailed       = errors.New("parse version failed")
	ErrParseGenTsFailed         = errors.New("parse generate ts failed")
	ErrInvalidVersion           = errors.New("invalid version")
)
