package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/h3ll0kitt1/observability/internal/config"
)

func Run(cfg *config.ServerConfig, r *chi.Mux) error {
	server := &http.Server{
		Handler: r,
		Addr:    cfg.Addr,
	}

	return server.ListenAndServe()
}
