package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/h3ll0kitt1/observability/internal/config"
	"github.com/h3ll0kitt1/observability/internal/storage"
	"github.com/h3ll0kitt1/observability/internal/storage/inmemory"
)

type application struct {
	storage storage.Storage
	router  *chi.Mux
}

func main() {

	cfg := config.NewServerConfig()

	s := inmemory.NewStorage()
	r := chi.NewRouter()

	app := &application{
		storage: s,
		router:  r,
	}

	app.setRouters()

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: app.router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
		return
	}
}
