package types

import "net/http"

type Server interface {
	Call(projectName string, data []byte) *http.Response
}
