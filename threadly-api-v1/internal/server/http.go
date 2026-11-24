package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"threadly-api-v1/internal/config"
)

func New(httpConfig config.HTTPConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         net.JoinHostPort(httpConfig.Host, strconv.Itoa(httpConfig.Port)),
		Handler:      handler,
		ReadTimeout:  httpConfig.ReadTimeout,
		WriteTimeout: httpConfig.WriteTimeout,
		IdleTimeout:  httpConfig.IdleTimeout,
	}
}

// Run serves until ctx is cancelled or the server fails, then shuts down
// gracefully within shutdownTimeout, force-closing if that deadline passes.
func Run(ctx context.Context, httpServer *http.Server, appLogger *zap.Logger, shutdownTimeout time.Duration) error {
	listener, err := net.Listen("tcp", httpServer.Addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", httpServer.Addr, err)
	}

	appLogger.Info("http server listening", zap.String("address", listener.Addr().String()))

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- httpServer.Serve(listener)
	}()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}

		return nil

	case <-ctx.Done():
		appLogger.Info("shutdown signal received")
	}

	shutdownContext, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownContext); err != nil {
		appLogger.Error("graceful shutdown failed", zap.Error(err))

		if closeErr := httpServer.Close(); closeErr != nil {
			appLogger.Error("force close server", zap.Error(closeErr))
		}
	}

	appLogger.Info("server stopped")

	return nil
}
