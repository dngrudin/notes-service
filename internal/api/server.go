package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dngrudin/notes-service/internal/config"
	"github.com/dngrudin/notes-service/pkg/logger"
)

type Server struct {
	config  config.HTTPServer
	handler http.Handler
}

func NewServer(config config.HTTPServer, handler http.Handler) *Server {
	return &Server{config: config, handler: handler}
}

func (s *Server) Start() {
	done := make(chan struct{})
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         s.config.Address,
		Handler:      s.handler,
		ReadTimeout:  s.config.Timeout,
		WriteTimeout: s.config.Timeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	go func() {
		defer close(done)

		<-sigterm
		slog.Info("Stopping server")

		ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Failed to stop server", logger.Err(err))
		}
	}()

	slog.Info("Starting server", slog.String("address", s.config.Address))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		signal.Stop(sigterm)
		slog.Error("Failed to start server", logger.Err(err))
		return
	}

	<-done

	slog.Info("Server stopped")
}
