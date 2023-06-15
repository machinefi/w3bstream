package storage

//go:generate toolkit gen enum StorageType
type StorageType uint8

const (
	STORAGE_TYPE_UNKNOWN StorageType = iota
	STORAGE_TYPE__S3
	STORAGE_TYPE__FILESYSTEM
	STORAGE_TYPE__IPFS
)
