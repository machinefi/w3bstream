package filesystem

import (
	"io"
	"os"
	"path/filepath"

	"github.com/machinefi/w3bstream/pkg/depends/conf/filesystem/amazonS3"
	"github.com/machinefi/w3bstream/pkg/enums"
)

type FileSystem struct {
	Type              enums.FileSystemMode `env:""`
	LocalRoot         string               `env:""`
	FileSizeLimit     int64                `env:""`
	S3Region          string               `env:""`
	S3AccessKeyID     string               `env:""`
	S3SecretAccessKey string               `env:""`
	S3SessionToken    string               `env:""`
	S3BucketName      string               `env:""`

	FileCli FileSystemOp
}

func (f *FileSystem) SetDefault() {
	if f.Type > enums.FILE_SYSTEM_MODE__S3 || f.Type < 0 {
		f.Type = enums.FILE_SYSTEM_MODE__LOCAL
	}
}

func (f *FileSystem) Init() error {
	switch f.Type {
	case enums.FILE_SYSTEM_MODE__S3:
		f.FileCli = amazonS3.NewAmazonS3(f.S3Region, f.S3AccessKeyID, f.S3SecretAccessKey, f.S3SessionToken, f.S3BucketName)
	default:
		f.FileCli = &LocalFileSystem{}
	}
	return nil
}

type FileSystemOp interface {
	Upload(key string, file []byte) error
	Read(key string) ([]byte, error)
	Delete(key string) error
}

type LocalFileSystem struct{}

// Upload key full path with filename
func (l *LocalFileSystem) Upload(key string, data []byte) error {
	var (
		fw  io.WriteCloser
		err error
	)
	dir, _ := filepath.Split(key)
	if !isDirExists(dir) {
		if err = os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}

	if fw, err = os.OpenFile(key, os.O_RDWR|os.O_CREATE, 0666); err != nil {
		return err
	}
	defer fw.Close()

	if _, err = fw.Write(data); err != nil {
		return err
	}

	return nil
}

func (l *LocalFileSystem) Read(key string) ([]byte, error) {
	return os.ReadFile(key)
}

func (l *LocalFileSystem) Delete(key string) error {
	return os.Remove(key)
}

func isDirExists(path string) bool {
	info, err := os.Stat(path)
	return (err == nil || os.IsNotExist(err)) && (info != nil && info.IsDir())
}
