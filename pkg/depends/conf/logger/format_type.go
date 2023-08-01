package logger

//go:generate toolkit gen enum FormatType
type FormatType uint8

const (
	FORMAT_TYPE_UNKNOWN FormatType = iota
	FORMAT_TYPE__JSON
	FORMAT_TYPE__TEXT
)
