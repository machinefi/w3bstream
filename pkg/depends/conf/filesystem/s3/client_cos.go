package confs3

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"hash"
	"net/url"
	"strings"
	"time"
)

func COSPresignedValues(db *ObjectDB, key string, exp time.Duration) url.Values {
	authTime := NewAuthTime(exp)
	signTime := authTime.signString()
	keyTime := authTime.keyString()
	signKey := calSignKey(db.SecretAccessKey.String(), keyTime)
	formatString := genFormatString("get", "/"+key, "", "")
	stringToSign := calStringToSign(sha1SignAlgorithm, keyTime, formatString)
	signature := calSignature(signKey, stringToSign)
	signedHeaderList := make([]string, 0)
	signedParameterList := make([]string, 0)

	values := url.Values{}

	values.Set("q-sign-algorithm", sha1SignAlgorithm)
	values.Set("q-ak", db.AccessKeyID)
	values.Set("q-sign-time", signTime)
	values.Set("q-objectKey-time", keyTime)
	values.Set("q-header-list", strings.Join(signedHeaderList, ";"))
	values.Set("q-url-param-list", strings.Join(signedParameterList, ";"))
	values.Set("q-signature", signature)

	return values
}

func NewAuthTime(du time.Duration) *AuthTime {
	if du == time.Duration(0) {
		du = defaultAuthExpire
	}
	signStartTime := time.Now()
	keyStartTime := signStartTime
	signEndTime := signStartTime.Add(du)
	keyEndTime := signEndTime
	return &AuthTime{
		SignStartTime: signStartTime,
		SignEndTime:   signEndTime,
		KeyStartTime:  keyStartTime,
		KeyEndTime:    keyEndTime,
	}
}

const (
	sha1SignAlgorithm = "sha1"
	md5SignAlgorithm  = "md5"
)
const defaultAuthExpire = time.Hour

// AuthTime is a struct storing the q-signSearch-time and q-key-time which are needed to generate signature
type AuthTime struct {
	SignStartTime time.Time
	SignEndTime   time.Time
	KeyStartTime  time.Time
	KeyEndTime    time.Time
}

func (a *AuthTime) signString() string {
	return fmt.Sprintf("%d;%d", a.SignStartTime.Unix(), a.SignEndTime.Unix())
}

func (a *AuthTime) keyString() string {
	return fmt.Sprintf("%d;%d", a.KeyStartTime.Unix(), a.KeyEndTime.Unix())
}

func calSignKey(secretKey, keyTime string) string {
	digest := HMAC(secretKey, keyTime, sha1SignAlgorithm)
	return fmt.Sprintf("%x", digest)
}

func calStringToSign(signAlgorithm, signTime, formatString string) string {
	h := sha1.New()
	h.Write([]byte(formatString))
	return fmt.Sprintf("%s\n%s\n%x\n", signAlgorithm, signTime, h.Sum(nil))
}

func calSignature(signKey, stringToSign string) string {
	digest := HMAC(signKey, stringToSign, sha1SignAlgorithm)
	return fmt.Sprintf("%x", digest)
}

func genFormatString(method string, url string, formatParameters, formatHeaders string) string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n", method, url,
		formatParameters, formatHeaders,
	)
}

func HMAC(key, msg, signMethod string) []byte {
	var hashFn func() hash.Hash
	switch signMethod {
	case sha1SignAlgorithm:
		hashFn = sha1.New
	case md5SignAlgorithm:
		hashFn = md5.New
	default:
		hashFn = sha1.New
	}
	h := hmac.New(hashFn, []byte(key))
	h.Write([]byte(msg))
	return h.Sum(nil)
}
