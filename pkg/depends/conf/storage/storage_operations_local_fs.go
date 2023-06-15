package storage

import (
	"io"
	"os"
	"path/filepath"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
)

type LocalFs struct {
	Root string `env:""`
}

func (l *LocalFs) Init() error {
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

func (l *LocalFs) Type() StorageType { return STORAGE_TYPE__FILESYSTEM }

func (l *LocalFs) SetDefault() {}

// Upload key full path with filename
func (l *LocalFs) Upload(key string, data []byte, chk ...HmacAlgType) error {
	var (
		fw  io.WriteCloser
		err error
	)

	path := filepath.Join(l.Root, key)
	if IsPathExists(path) {
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

func (l *LocalFs) Read(key string, chk ...HmacAlgType) (data []byte, sum []byte, err error) {
	data, err = os.ReadFile(filepath.Join(l.Root, key))
	if err != nil {
		return
	}
	t := HMAC_ALG_TYPE__MD5
	if len(chk) > 0 && chk[0] != 0 {
		t = chk[0]
	}
	sum = t.Sum(data)
	return
}

func (l *LocalFs) Delete(key string) error {
	return os.Remove(filepath.Join(l.Root, key))
}

func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
