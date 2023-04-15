package stringsx

import (
	"math/rand"
	"time"
)

var (
	visibleChars    = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	visibleCharsLen = len(visibleChars)
)

func GenRandomVisibleString(length int) string {
	if length == 0 {
		length = 16
	}
	result := make([]byte, 0, length)
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		idx := rand.Intn(visibleCharsLen)
		result = append(result, visibleChars[idx])
	}
	return string(result)
}

func SubStringWithLength(str string, start, length int) string {
	if start > len(str) {
		return ""
	}
	if length < 0 {
		length = 0
	}

	var i, n, leng, end int

	if start < 0 {
		end = len(str) + start
		if end > 0 {
			start = end - length
		} else {
			return ""
		}
	}
	if start >= 0 {
		leng = start + length
	} else {
		start = 0
		leng = end
	}
	for i = range str {
		if n == leng {
			break
		}
		n++
	}

	return str[start:i]
}
