package enums

//go:generate toolkit gen enum ConfigType
type ConfigType uint8

const (
	CONFIG_TYPE_UNKNOWN ConfigType = iota
	CONFIG_TYPE__PROJECT_SCHEMA
	CONFIG_TYPE__INSTANCE_CACHE
	CONFIG_TYPE__PROJECT_ENV
)
