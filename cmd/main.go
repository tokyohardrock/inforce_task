package main

import (
	"context"
	"errors"
	"inforce_task/internal/config"
	"inforce_task/internal/handler"
	"inforce_task/internal/logger"
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
		slog.Error(err.Error())
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /events", handler.CreateEvent)

	server, err := server.New(&cfg.Server, mux)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	fmt.Println("Server runs on ", server.Addr)

	<-mainCtx.Done()
	slog.Info("shutting down server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server stopped gracefully")
}
