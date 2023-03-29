// Package httpserver implements HTTP server.
package httpserver

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server http.Server
}

func (s *Server) Start() {
	s.server.ListenAndServe()
}

func (s *Server) ShutDown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func New(addr string, writeTimeout time.Duration,
	readTimeout time.Duration, handler http.Handler) *Server {
	return &Server{
		server: http.Server{
			Addr:         addr,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			Handler:      handler,
		},
	}
}
