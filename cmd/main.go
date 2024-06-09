package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"zota-dev-challenge/internal"
)

func startServer(lc fx.Lifecycle, logger *zap.Logger, router *chi.Mux) {
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting server on :8080")
			// Start server in a new goroutine
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal("Failed to start server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping server")
			return server.Shutdown(ctx)
		},
	})
}

// @title			Merchant Server
// @version		1.0
// @description	This is a simple merchant's server implementing Zota payment gateway.
// @host			localhost:8080
// @BasePath		/api/v1
func main() {
	app := fx.New(
		internal.AppModules,
		fx.Invoke(startServer),
	)

	// Handle SIGINT and SIGTERM signals for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		_ = app.Stop(context.Background())
	}()

	app.Run()
}
