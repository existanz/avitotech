package main

import (
	"avitotech/internal/server"
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setupLogger() {
	logLevel := new(slog.LevelVar)
	options := &slog.HandlerOptions{Level: logLevel}

	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel.Set(slog.LevelDebug)
		options.AddSource = true
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, options))
	slog.SetDefault(logger)
}

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		slog.Info("Server forced to shutdown with", "error", err)
	}

	slog.Info("Server exiting")

	done <- true
}

func main() {
	setupLogger()
	srv := server.NewServer()

	done := make(chan bool, 1)

	go gracefulShutdown(srv, done)

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("http server error: %s", err)
	}

	<-done
	slog.Info("Graceful shutdown complete.")
}
