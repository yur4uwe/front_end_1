package server

import (
	"context"
	"log"
	"net/http"
)

type Server struct {
	srv *http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	log.Println("Starting server on", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("Stopping server...")
	return s.srv.Shutdown(ctx)
}
