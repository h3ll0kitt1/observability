package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/h3ll0kitt1/observability/internal/config"
	"github.com/h3ll0kitt1/observability/internal/logger"
	"github.com/h3ll0kitt1/observability/internal/storage"
	"github.com/h3ll0kitt1/observability/internal/storage/file"
	"github.com/h3ll0kitt1/observability/internal/storage/inmemory"
	"github.com/h3ll0kitt1/observability/internal/storage/sql"
)

type application struct {
	config  *config.ServerConfig
	storage storage.Storage
	backup  storage.Storage
	router  *chi.Mux
	logger  *zap.SugaredLogger
}

func (app *application) loadFromDisk() error {
	list, err := app.backup.GetList(context.Background())
	if err != nil {
		return err
	}

	if err := app.storage.UpdateList(context.Background(), list); err != nil {
		return err
	}
	return nil
}

func (app *application) flushToDisk() {
	ticker := time.NewTicker(app.config.StoreInterval)

	for range ticker.C {
		app.flush()
	}
}

func (app *application) flush() {
	list, err := app.storage.GetList(context.Background())
	if err != nil {
		log.Print(err)
	}

	if err := app.backup.UpdateList(context.Background(), list); err != nil {
		log.Print(err)
	}
}

func main() {

	cfg := config.NewServerConfig()

	var s storage.Storage

	s = inmemory.NewStorage()
	if cfg.Database != "" {
		s = sql.NewStorage(cfg)
	}

	r := chi.NewRouter()
	l := logger.NewLogger()

	defer l.Sync()

	app := &application{
		config:  cfg,
		storage: s,
		router:  r,
		logger:  l,
	}

	app.setRouters()

	if cfg.Restore {
		app.backup = file.NewStorage(cfg.FileStoragePath)
		if err := app.loadFromDisk(); err != nil {
			log.Fatal(err)
		}
	}

	if app.config.StoreInterval > 0 {
		go app.flushToDisk()
	}

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: app.router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
		return
	}
}
