package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func FileMD5(path string) (string, error) {
	f, err := os.Open(path)
	if nil != err {
		return "", err
	}
	defer f.Close()

	hash := md5.New()
	_, err = io.Copy(hash, f)
	if nil != err {
		return "", err
	}
	sum := hash.Sum(nil)
	return fmt.Sprintf("%x", sum), nil
}
