package filesystem

import (
	"io"
	"os"
)

type FileSystemOp interface {
	Upload(key string, file []byte) error
	Read(key string) ([]byte, error)
	Delete(key string) error
}

type LocalFileSystem struct{}

func (l *LocalFileSystem) Upload(key string, file []byte) error {
	var (
		fw  io.WriteCloser
		err error
	)
	if fw, err = os.OpenFile(key, os.O_RDWR|os.O_CREATE, 0666); err != nil {
		return err
	}
	defer fw.Close()

	if _, err = fw.Write(file); err != nil {
		return err
	}

	return nil
}

func (l *LocalFileSystem) Read(key string) ([]byte, error) {
	data, err := os.ReadFile(key)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (l *LocalFileSystem) Delete(key string) error {
	return nil
}
