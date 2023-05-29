package local

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/depends/util"
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
	return l.UploadWithMD5(key, "", data)
}

func (l *LocalFileSystem) UploadWithMD5(key, md5 string, data []byte) error {
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

	if md5 != "" {
		sum, err := util.FileMD5(path)
		if err != nil {
			return err
		}
		if sum != md5 {
			return errors.New("md5 not match")
		}
	}

	return nil
}

func (l *LocalFileSystem) Read(key string) ([]byte, error) {
	return l.ReadWithMD5(key, "")
}

func (l *LocalFileSystem) ReadWithMD5(key, md5 string) ([]byte, error) {
	data, err := os.ReadFile(l.path(key))
	if err != nil {
		return nil, err
	}

	if md5 != "" {
		sum, err := util.ByteMD5(data)
		if err != nil {
			return nil, err
		}
		if sum != md5 {
			return nil, errors.New("md5 not match")
		}
	}

	return data, err
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
