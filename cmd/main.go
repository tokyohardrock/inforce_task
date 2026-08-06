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
)

func main() {
	log := logger.InitLogger()

	slog.SetDefault(log)

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

	if err = server.ListenAndServe(); err != nil {
		slog.Error(err.Error())
		return
	}
}
