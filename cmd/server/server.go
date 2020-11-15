package main

import (
	"context"
	"net/http"
	"time"
)

// Server is a common http server
type Server struct {
	server *http.Server
}

// New creates a new Server
func New(port string, router http.Handler, timeout time.Duration) *Server {
	return &Server{
		server: &http.Server{
			Handler:      router,
			Addr:         ":" + port,
			WriteTimeout: timeout,
			ReadTimeout:  timeout,
		},
	}
}

// Run starts the server
func (s *Server) Run(ctx context.Context) error {
	ec := make(chan error)
	go func() {
		ec <- s.server.ListenAndServe()
	}()

	var err error
	select {
	case err = <-ec:
	case <-ctx.Done():
		err = s.server.Shutdown(ctx)
	}
	return err
}
