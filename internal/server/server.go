package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Run(endpoint string, r *chi.Mux) error {
	server := &http.Server{
		Handler: r,
		Addr:    endpoint,
	}

	return server.ListenAndServe()
}
