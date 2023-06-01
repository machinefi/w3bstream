package local

import (
	"io"
	"os"
	"path/filepath"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/depends/conf/filesystem"
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
	return l.UploadWithChecksum(key, "", "", data)
}

func (l *LocalFileSystem) UploadWithChecksum(key, sum, algorithm string, data []byte) error {
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

	if sum != "" && algorithm != "" {
		targetSum, err := util.ChecksumByFile(path, algorithm)
		if err != nil {
			return err
		}
		if sum != targetSum {
			return filesystem.ErrChecksumNotMatch
		}
	}

	return nil
}

func (l *LocalFileSystem) Read(key string) ([]byte, error) {
	return l.ReadWithChecksum(key, "", "")
}

func (l *LocalFileSystem) ReadWithChecksum(key, sum, algorithm string) ([]byte, error) {
	data, err := os.ReadFile(l.path(key))
	if err != nil {
		return nil, err
	}

	if sum != "" && algorithm != "" {
		targetSum := util.Checksum(data, algorithm)
		if sum != targetSum {
			return nil, filesystem.ErrChecksumNotMatch
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
