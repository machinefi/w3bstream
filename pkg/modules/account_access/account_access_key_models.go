package account_access

type CreateReq struct {
	// Name access token name
	Name string `json:"name"`
	// ExpirationDays access token valid in ExpirationDays, if 0 means token will not be expired.
	ExpirationDays int `json:"expirationDays,omitempty"`
	// Description access token description
	Description string `json:"description"`
	// TODO _Privileges access token privileges
	_Privileges interface{}
}
