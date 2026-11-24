package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"threadly-api-v1/internal/config"
	"threadly-api-v1/internal/handler"
	"threadly-api-v1/internal/infrastructure/database"
	sqlcdb "threadly-api-v1/internal/infrastructure/database/sqlc"
	"threadly-api-v1/internal/infrastructure/logger"
	"threadly-api-v1/internal/repository"
	"threadly-api-v1/internal/routes"
	"threadly-api-v1/internal/server"
	"threadly-api-v1/internal/service"
)

// Run is the composition root: it builds every dependency, wires them
// together and serves HTTP until ctx is cancelled.
func Run(ctx context.Context) error {
	applicationConfig, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	appLogger, err := logger.New(applicationConfig.App.Environment, applicationConfig.Log.Level)
	if err != nil {
		return fmt.Errorf("create logger: %w", err)
	}
	defer func() { _ = appLogger.Sync() }()

	if err := run(ctx, applicationConfig, appLogger); err != nil {
		appLogger.Error("server exited with error", zap.Error(err))
		return err
	}

	return nil
}

func run(ctx context.Context, applicationConfig *config.Config, appLogger *zap.Logger) error {
	pool, err := database.NewPool(ctx, applicationConfig.Databases.Primary)
	if err != nil {
		return fmt.Errorf("connect primary database: %w", err)
	}
	defer pool.Close()

	if applicationConfig.App.Environment == config.EnvironmentProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	queries := sqlcdb.New(pool)
	userService := service.NewUserService(repository.NewUserRepository(queries))

	router, err := routes.NewRouter(applicationConfig.HTTP, appLogger, routes.Handlers{
		Health: handler.NewHealthHandler(pool, appLogger),
		User:   handler.NewUserHandler(userService),
	})
	if err != nil {
		return fmt.Errorf("build router: %w", err)
	}

	appLogger.Info("starting server",
		zap.String("app", applicationConfig.App.Name),
		zap.String("environment", applicationConfig.App.Environment),
	)

	return server.Run(ctx, server.New(applicationConfig.HTTP, router), appLogger, applicationConfig.HTTP.ShutdownTimeout)
}
