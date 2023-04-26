package enums

//go:generate toolkit gen enum FileSystemMode
type FileSystemMode uint8

const (
	FILE_SYSTEM_MODE_UNKNOWN FileSystemMode = iota
	FILE_SYSTEM_MODE__LOCAL
	FILE_SYSTEM_MODE__S3
)
