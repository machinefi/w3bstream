package enums

//go:generate toolkit gen enum ApiOperatorAttr
type ApiOperatorAttr uint8

const (
	API_OPERATOR_ATTR_UNKNOWN     ApiOperatorAttr = iota
	API_OPERATOR_ATTR__COMMON                     // common operator for all users with authorization, used for default attr
	API_OPERATOR_ATTR__INTERNAL                   // internal for debugging or maintaining only
	API_OPERATOR_ATTR__ADMIN_ONLY                 // only admin can access
	API_OPERATOR_ATTR__PUBLIC                     // can access without authorization
)
