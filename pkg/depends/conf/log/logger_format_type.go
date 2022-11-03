package log

//go:generate toolkit gen enum LoggerFormatType
type LoggerFormatType uint8

const (
	LOGGER_FORMAT_TYPE_UNKNOWN LoggerFormatType = iota
	LOGGER_FORMAT_TYPE__JSON
	LOGGER_FORMAT_TYPE__TEXT
)
