package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/h3ll0kitt1/observability/internal/config"
	"github.com/h3ll0kitt1/observability/internal/logger"
	"github.com/h3ll0kitt1/observability/internal/storage"
	"github.com/h3ll0kitt1/observability/internal/storage/inmemory"
)

type application struct {
	storage storage.Storage
	router  *chi.Mux
	logger  *zap.SugaredLogger
}

func main() {

	cfg := config.NewServerConfig()

	s := inmemory.NewStorage()
	r := chi.NewRouter()
	l := logger.NewLogger()

	defer l.Sync()

	app := &application{
		storage: s,
		router:  r,
		logger:  l,
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
