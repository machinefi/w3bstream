package stringsx

import (
	"math/rand"
	"time"
)

var (
	visibleChars    []byte
	visibleCharsLen int

	src = rand.NewSource(time.Now().UnixNano())
)

var GenRandomVisibleString = GenRandomVisibleStringV2

const (
	letterIdxBits = 7                    // 7 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func init() {
	visibleCharsLen = 95 // all visible chars [32,126] 95=1011111b
	visibleChars = make([]byte, visibleCharsLen)
	for i := byte(0); int(i) < visibleCharsLen; i++ {
		visibleChars[i] = i + 32
	}
}

func GenRandomVisibleStringV1(length int, excluded ...byte) string {
	if length == 0 {
		return ""
	}
	result := make([]byte, 0, length)
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; {
	SKIP:
		idx := rand.Intn(visibleCharsLen)
		c := visibleChars[idx]
		if len(excluded) > 0 {
			for _, v := range excluded {
				if v == c {
					goto SKIP
				}
			}
		}
		result = append(result, c)
		i++
	}
	return string(result)
}

func GenRandomVisibleStringV2(n int, excluded ...byte) string {
	if n == 0 {
		return ""
	}
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(visibleChars) {
			c := visibleChars[idx]
			if len(excluded) > 0 {
				for _, v := range excluded {
					if v == c {
						goto SKIP
					}
				}
			}
			b[i] = c
			i--
		}
	SKIP:
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
