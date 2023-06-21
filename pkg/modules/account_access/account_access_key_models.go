package account_access

import "github.com/machinefi/w3bstream/pkg/types"

type CreateReq struct {
	// Name access token name
	Name string `json:"name"`
	// ExpirationDays access token valid in ExpirationDays, if 0 means token will not be expired.
	ExpirationDays int `json:"expirationDays,omitempty"`
	// Description access token description
	Desc string `json:"desc"`
	// TODO _Privileges access token privileges
	_Privileges interface{}
}

type CreateRsp struct {
	Name      string           `json:"name"`
	AccessKey string           `json:"accessKey"`
	ExpiredAt *types.Timestamp `json:"expiredAt,omitempty"`
	Desc      string           `json:"desc,omitempty"`
}
