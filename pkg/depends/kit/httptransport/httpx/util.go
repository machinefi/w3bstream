package httpx

import (
	"net"
	"net/http"
	"strings"
)

func ClientIP(r *http.Request) string {
	clientIP := ClientIPByHeaderForwardedFor(r.Header.Get(HeaderForwardedFor))
	if clientIP != "" {
		return clientIP
	}

	clientIP = ClientIPByHeaderRealIP(r.Header.Get(HeaderRealIP))
	if clientIP != "" {
		return clientIP
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// ClientIPByHeaderForwardedFor
// X-Forwarded-For: client, proxy1, proxy2
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For
func ClientIPByHeaderForwardedFor(headerForwardedFor string) string {
	if index := strings.IndexByte(headerForwardedFor, ','); index >= 0 {
		return headerForwardedFor[0:index]
	}
	return strings.TrimSpace(headerForwardedFor)
}

// ClientIPByHeaderRealIP
// X-Forwarded-For: client, proxy1, proxy2
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For
func ClientIPByHeaderRealIP(headerRealIP string) string {
	return strings.TrimSpace(headerRealIP)
}
