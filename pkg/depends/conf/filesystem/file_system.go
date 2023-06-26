package filesystem

import "errors"

type FileSystemOp interface {
	Upload(key string, file []byte) error
	Read(key string) ([]byte, error)
	Delete(key string) error
	StatObject(key string) (*ObjectMeta, error)
}

var (
	ErrInvalidObjectKey  = errors.New("invalid object key")
	ErrNotExistObjectKey = errors.New("not exist object key")
)
