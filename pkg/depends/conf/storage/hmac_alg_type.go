package storage

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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

func (v HmacAlgType) Sum(content []byte) []byte {
	switch v {
	case HMAC_ALG_TYPE__SHA1:
		sum := sha1.Sum(content)
		return sum[:]
	case HMAC_ALG_TYPE__SHA256:
		sum := sha256.Sum256(content)
		return sum[:]
	default:
		sum := md5.Sum(content)
		return sum[:]
	}
}

func (v HmacAlgType) HexSum(content []byte) string {
	return hex.EncodeToString(v.Sum(content))
}

func (v HmacAlgType) Base64Sum(content []byte) string {
	return base64.StdEncoding.EncodeToString(v.Sum(content))
}

func (v HmacAlgType) Type() string { return strings.ToLower(v.String()) }
