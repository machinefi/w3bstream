package local

import (
	"io"
	"os"
	"path/filepath"
)

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
