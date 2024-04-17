package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"nishojib/gotrello/internal/option"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nedpals/supabase-go"
	"github.com/uptrace/bun"
)

// Server represents an HTTP server.
type Server struct {
	*http.Server
	wg sync.WaitGroup
}

type limiter struct {
	rps     int
	enabled bool
}

func NewLimiter(rps int, enabled bool) limiter {
	return limiter{rps, enabled}
}

// New creates a new server with the provided options.
func New(
	db *bun.DB,
	sbClient *supabase.Client,
	limiter limiter,
	opts ...option.Option[Server],
) *Server {
	s := &Server{
		Server: &http.Server{
			Addr:         ":8080",
			Handler:      Routes(db, sbClient, limiter.rps, limiter.enabled),
			IdleTimeout:  1 * time.Minute,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			ErrorLog:     slog.NewLogLogger(slog.Default().Handler(), slog.LevelError),
		},
	}

	return s.Append(opts...)
}

// Serve starts the server and blocks until the server is stopped.
func (s *Server) Serve(environment string) error {
	shutdownError := make(chan error, 1)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		slog.Info("shutting down server", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := s.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		slog.Info("completing background tasks", "addr", s.Addr)

		s.wg.Wait()
		shutdownError <- nil
	}()

	slog.Info("starting server", "addr", s.Addr, "env", environment)

	err := s.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	slog.Info("stopped server", "addr", s.Addr)

	return nil

}

// Append applies the provided options to the server.
func (s *Server) Append(opts ...option.Option[Server]) *Server {
	for _, opt := range opts {
		opt.Apply(s)
	}

	return s
}

// WithPort sets the address for the server.
func WithPort(port int) option.Option[Server] {
	return option.OptionFunc[Server](func(s *Server) {
		s.Server.Addr = fmt.Sprintf(":%d", port)
	})
}

// WithHandler sets the handler for the server.
func WithHandler(handler http.Handler) option.Option[Server] {
	return option.OptionFunc[Server](func(s *Server) {
		s.Server.Handler = handler
	})
}

// WithIdleTimeout sets the idle timeout for the server.
func WithIdleTimeout(timeout time.Duration) option.Option[Server] {
	return option.OptionFunc[Server](func(s *Server) {
		s.Server.IdleTimeout = timeout
	})
}

// WithReadTimeout sets the read timeout for the server.
func WithReadTimeout(timeout time.Duration) option.Option[Server] {
	return option.OptionFunc[Server](func(s *Server) {
		s.Server.ReadTimeout = timeout
	})
}

// WithWriteTimeout sets the write timeout for the server.
func WithWriteTimeout(timeout time.Duration) option.Option[Server] {
	return option.OptionFunc[Server](func(s *Server) {
		s.Server.WriteTimeout = timeout
	})
}

// WithErrorLog sets the error log for the server.
func WithErrorLog(logger *log.Logger) option.Option[Server] {
	return option.OptionFunc[Server](func(s *Server) {
		s.Server.ErrorLog = logger
	})
}
