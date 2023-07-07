package api

import (
	"net/http"
	"net/http/httptest"
)

type Server struct {
	srv *http.ServeMux
}

func (s *Server) Serve(req *http.Request) *http.Response {
	resp := httptest.NewRecorder()
	s.srv.ServeHTTP(resp, req)
	return resp.Result()
}

func NewServer() *Server {
	srv := http.NewServeMux()

	srv.HandleFunc("/system/hello", hello)

	return &Server{srv}
}
