package server

import (
	"fmt"
	"inforce_task/internal/config"
	"net/http"
)

func New(cfg *config.HTTPServer, handler http.Handler) (*http.Server, error) {
	const fn = "server.New"

	if cfg == nil {
		return nil, fmt.Errorf("%s: config is nil", fn)
	}
	if handler == nil {
		return nil, fmt.Errorf("%s: handler is nil", fn)
	}

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
	}, nil
}
