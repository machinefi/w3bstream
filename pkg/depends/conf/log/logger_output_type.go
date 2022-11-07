package log

//go:generate toolkit gen enum LoggerOutputType
type LoggerOutputType uint8

const (
	LOGGER_OUTPUT_TYPE_UNKNOWN LoggerOutputType = iota
	LOGGER_OUTPUT_TYPE__ALWAYS
	LOGGER_OUTPUT_TYPE__ON_FAILURE
	LOGGER_OUTPUT_TYPE__NEVER
)
