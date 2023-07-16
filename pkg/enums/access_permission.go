package enums

//go:generate toolkit gen enum AccessPermission
type AccessPermission uint8

const (
	ACCESS_PERMISSION_UNKNOWN AccessPermission = iota
	ACCESS_PERMISSION__NO_ACCESS
	ACCESS_PERMISSION__READONLY
	ACCESS_PERMISSION__READ_WRITE
)
