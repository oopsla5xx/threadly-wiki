package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"threadly-api-v1/internal/config"
	"threadly-api-v1/internal/infrastructure/database"
	"threadly-api-v1/internal/infrastructure/logger"
	"threadly-api-v1/internal/routes"
)

func main() {
	applicationConfig, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	appLogger, err := logger.New(applicationConfig.App.Environment, applicationConfig.Log.Level)
	if err != nil {
		log.Fatalf("create logger: %v", err)
	}
	defer func() { _ = appLogger.Sync() }()

	if err := run(applicationConfig, appLogger); err != nil {
		appLogger.Error("server exited with error", zap.Error(err))
		_ = appLogger.Sync()
		os.Exit(1)
	}
}

func run(applicationConfig *config.Config, appLogger *zap.Logger) error {
	shutdownSignalContext, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	pool, err := database.NewPool(shutdownSignalContext, applicationConfig.Databases.Primary)
	if err != nil {
		return fmt.Errorf("connect primary database: %w", err)
	}
	defer pool.Close()

	if applicationConfig.App.Environment == config.EnvironmentProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router, err := routes.NewRouter(applicationConfig.HTTP, pool, appLogger)
	if err != nil {
		return fmt.Errorf("build router: %w", err)
	}

	serverAddress := fmt.Sprintf(
		"%s:%d",
		applicationConfig.HTTP.Host,
		applicationConfig.HTTP.Port,
	)

	httpServer := &http.Server{
		Addr:         serverAddress,
		Handler:      router,
		ReadTimeout:  applicationConfig.HTTP.ReadTimeout,
		WriteTimeout: applicationConfig.HTTP.WriteTimeout,
		IdleTimeout:  applicationConfig.HTTP.IdleTimeout,
	}

	appLogger.Info("starting server",
		zap.String("app", applicationConfig.App.Name),
		zap.String("environment", applicationConfig.App.Environment),
		zap.String("address", serverAddress),
	)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}

	case <-shutdownSignalContext.Done():
		appLogger.Info("shutdown signal received")
	}

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		applicationConfig.HTTP.ShutdownTimeout,
	)
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
