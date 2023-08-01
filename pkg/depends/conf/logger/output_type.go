package logger

//go:generate toolkit gen enum OutputType
type OutputType uint8

const (
	OUTPUT_TYPE_UNKNOWN OutputType = iota
	OUTPUT_TYPE__ALWAYS
	OUTPUT_TYPE__ON_FAILURE
	OUTPUT_TYPE__NEVER
)
