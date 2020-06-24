package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Integrity-178B/url-fetcher/internal/pkg/log"
)

const shutdownTime = 5 * time.Second

// Config keeps server configuration
type Config struct {
	Host string
	Port string
}

func (c Config) addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// Server is http server implementation
type Server struct {
	*http.Server
	logger *log.Logger
}

// NewServer creates new server instance
func NewServer(conf *Config, router http.Handler) *Server {
	logger := log.NewLogger("[server] ")

	return &Server{
		Server: &http.Server{
			Addr:     conf.addr(),
			Handler:  router,
			ErrorLog: logger.Logger,
		},
		logger: logger,
	}
}

// ListenAndServe listens on the tcp address and serves the requests
func (s Server) ListenAndServe(ctx context.Context) {
	go func() {
		if err := s.Server.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Print(err)
		}
	}()

	s.logger.Printf("started on %s", s.Addr)
	<-ctx.Done()

	s.logger.Printf("gracefully stopping")
	s.Shutdown()
}

// Shutdown gracefully shutdown the server
func (s Server) Shutdown() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTime)
	defer cancel()

	if err := s.Server.Shutdown(shutdownCtx); err == nil {
		s.logger.Print("stopped gracefully")
	} else {
		s.logger.Printf("shutdown failed: %s", err)
	}
}
