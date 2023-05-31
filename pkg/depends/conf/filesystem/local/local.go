package local

import (
	"io"
	"os"
	"path/filepath"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
)

type LocalFileSystem struct {
	Root string `env:""`
}

func (l *LocalFileSystem) Init() error {
	if l.Root == "" {
		tmp := os.Getenv("TMPDIR")
		if tmp == "" {
			tmp = "/tmp"
		}
		serviceName := os.Getenv(consts.EnvProjectName)
		if serviceName == "" {
			serviceName = "service_tmp"
		}
		l.Root = filepath.Join(tmp, serviceName)
	}
	return os.MkdirAll(filepath.Join(l.Root, os.Getenv(consts.EnvResourceGroup)), 0777)
}

func (l *LocalFileSystem) SetDefault() {}

// Upload key full path with filename
func (l *LocalFileSystem) Upload(key string, data []byte) error {
	var (
		fw  io.WriteCloser
		err error
	)

	path := filepath.Join(l.Root, key)
	if isPathExists(path) {
		return nil
	}

	if fw, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666); err != nil {
		return err
	}
	defer fw.Close()

	if _, err = fw.Write(data); err != nil {
		return err
	}

	return nil
}

func (l *LocalFileSystem) Read(key string) ([]byte, error) {
	return os.ReadFile(l.path(key))
}

func (l *LocalFileSystem) Delete(key string) error {
	return os.Remove(l.path(key))
}

func (l *LocalFileSystem) path(name string) string {
	return filepath.Join(l.Root, name)
}

func isPathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
