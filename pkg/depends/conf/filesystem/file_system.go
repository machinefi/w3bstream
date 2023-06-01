package filesystem

import "github.com/pkg/errors"

type FileSystemOp interface {
	Upload(key string, file []byte) error
	UploadWithChecksum(key, sum, algorithm string, file []byte) error
	Read(key string) ([]byte, error)
	ReadWithChecksum(key, sum, algorithm string) ([]byte, error)
	Delete(key string) error
}

var ErrChecksumNotMatch = errors.New("checksum not match")
