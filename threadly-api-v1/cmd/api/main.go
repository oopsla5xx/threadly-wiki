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
	"time"

	"github.com/gin-gonic/gin"

	"threadly-api-v1/internal/config"
)

func main() {
	applicationConfig, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if applicationConfig.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("configure trusted proxies: %v", err)
	}

	router.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

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

	log.Printf(
		"starting %s in %s mode",
		applicationConfig.App.Name,
		applicationConfig.App.Environment,
	)
	log.Printf("listening on %s", serverAddress)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- httpServer.ListenAndServe()
	}()

	shutdownSignalContext, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server: %v", err)
		}

	case <-shutdownSignalContext.Done():
		log.Println("shutdown signal received")
	}

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := httpServer.Shutdown(shutdownContext); err != nil {
		log.Printf("graceful shutdown failed: %v", err)

		if closeErr := httpServer.Close(); closeErr != nil {
			log.Printf("force close server: %v", closeErr)
		}
	}

	log.Println("server stopped")
}
