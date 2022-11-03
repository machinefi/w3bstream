package mws

import "net/http"

func HealthCheckHandler() func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return &healthCheckHandler{next: handler}
	}
}

type healthCheckHandler struct{ next http.Handler }

func (h *healthCheckHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if (req.Method == http.MethodHead || req.Method == http.MethodGet) && req.URL.Path == "/" {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	h.next.ServeHTTP(rw, req)
}
