package main

import (
	"context"
	"errors"
	"inforce_task/internal/config"
	"inforce_task/internal/database"
	"inforce_task/internal/handler"
	"inforce_task/internal/logger"
	"inforce_task/internal/repo"
	"inforce_task/internal/server"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log := logger.InitLogger()
	slog.SetDefault(log)

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.MustLoad()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	pool, err := database.NewPostgresPool(mainCtx, cfg.DB)
	if err != nil {
		slog.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	eventRepo := repo.NewPostgresEventRepository(pool)

	eventHandler := handler.NewEventHandler(eventRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /events", eventHandler.CreateEvent)
	mux.HandleFunc("GET /events", eventHandler.GetEvents)

	srv, err := server.New(&cfg.Server, mux)
	if err != nil {
		slog.Error("failed to init server", "error", err)
		os.Exit(1)
	}

	go func() {
		slog.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server HTTP listen error", "error", err)
			stop()
		}
	}()

	<-mainCtx.Done()
	slog.Info("shutting down server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server stopped gracefully")
}
