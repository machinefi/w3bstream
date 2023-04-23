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
