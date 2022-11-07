package jwt

import "net/http"

//go:generate toolkit gen status Error
type Error int

func (Error) ServiceCode() int {
	return 999 * 1e3
}

const (
	Unauthorized Error = http.StatusUnauthorized*1e6 + iota + 1
	InvalidToken
	InvalidClaim
)
