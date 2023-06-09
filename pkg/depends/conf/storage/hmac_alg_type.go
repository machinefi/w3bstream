package storage

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

//go:generate toolkit gen enum HmacAlgType
type HmacAlgType uint8

const (
	HMAC_ALG_TYPE_UNKNOWN HmacAlgType = iota
	HMAC_ALG_TYPE__MD5
	HMAC_ALG_TYPE__SHA1
	HMAC_ALG_TYPE__SHA256
)

func (v HmacAlgType) Sum(key, content []byte) []byte {
	var hashFn = md5.New
	switch v {
	case HMAC_ALG_TYPE__MD5:
		hashFn = md5.New
	case HMAC_ALG_TYPE__SHA1:
		hashFn = sha1.New
	case HMAC_ALG_TYPE__SHA256:
		hashFn = sha256.New
	}

	h := hmac.New(hashFn, key)
	h.Write(content)
	return h.Sum(nil)
}

func (v HmacAlgType) HexSum(key, content []byte) string {
	return fmt.Sprintf("%x", v.Sum(key, content))
}

func (v HmacAlgType) Base64Sum(key, content []byte) string {
	return base64.StdEncoding.EncodeToString(v.Sum(key, content))
}

func (v HmacAlgType) Type() string { return strings.ToLower(v.String()) }
