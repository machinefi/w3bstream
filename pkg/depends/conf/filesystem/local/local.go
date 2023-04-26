package local

import (
	"io"
	"os"
	"path/filepath"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
)

type LocalFileSystem struct {
	Root          string `env:""`
	FilesizeLimit int64  `env:""`
	DiskReserve   int64  `env:""`
}

func (l *LocalFileSystem) Init() error {
	err := os.MkdirAll(l.Root, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (l *LocalFileSystem) SetDefault() {
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
	if l.FilesizeLimit == 0 {
		l.FilesizeLimit = 1024 * 1024
	}
	if l.DiskReserve == 0 {
		l.DiskReserve = 20 * 1024 * 1024
	}
}

// Upload key full path with filename
func (l *LocalFileSystem) Upload(md5 string, data []byte) error {
	var (
		fw  io.WriteCloser
		err error
	)

	path := filepath.Join(l.Root, md5)
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

func (l *LocalFileSystem) Read(md5 string) ([]byte, error) {
	return os.ReadFile(filepath.Join(l.Root, md5))
}

func (l *LocalFileSystem) Delete(md5 string) error {
	return os.Remove(filepath.Join(l.Root, md5))
}

func isPathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
