package util

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

const (
	Md5Algorithm    = "md5"
	Sha1Algorithm   = "sha1"
	Sha256Algorithm = "sha256"
)

func ChecksumByFile(path, algorithm string) (string, error) {
	f, err := os.Open(path)
	if nil != err {
		return "", err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return Checksum(data, algorithm), nil
}

func Checksum(bytes []byte, algorithm string) string {
	var sum string
	switch algorithm {
	case Md5Algorithm:
		s := md5.Sum(bytes)
		sum = fmt.Sprintf("%x", s)
	case Sha1Algorithm:
		sha1sum := sha1.Sum(bytes)
		sum = fmt.Sprintf("%x", sha1sum)
	case Sha256Algorithm:
		sha256sum := sha256.Sum256(bytes)
		sum = fmt.Sprintf("%x", sha256sum)
	}

	return sum
}
