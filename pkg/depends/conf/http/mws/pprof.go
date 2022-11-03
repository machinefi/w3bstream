package mws

import (
	"net/http"
	"net/http/pprof"
	"strings"
)

func PProfHandler(enabled bool) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return &pprofHandler{
			enabled: enabled,
			next:    handler,
		}
	}
}

type pprofHandler struct {
	enabled bool
	next    http.Handler
}

func (h *pprofHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if h.enabled && strings.HasPrefix(req.URL.Path, "/debug/pprof") {
		switch req.URL.Path {
		case "/debug/pprof/cmdline":
			pprof.Cmdline(rw, req)
			return
		case "/debug/pprof/profile":
			pprof.Profile(rw, req)
			return
		case "/debug/pprof/symbol":
			pprof.Symbol(rw, req)
			return
		case "/debug/pprof/trace":
			pprof.Trace(rw, req)
			return
		default:
			pprof.Index(rw, req)
			return
		}
	}
	h.next.ServeHTTP(rw, req)
}
