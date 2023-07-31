package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/h3ll0kitt1/observability/internal/config"
	"github.com/h3ll0kitt1/observability/internal/disk"
	"github.com/h3ll0kitt1/observability/internal/logger"
	"github.com/h3ll0kitt1/observability/internal/storage"
	"github.com/h3ll0kitt1/observability/internal/storage/inmemory"
)

type application struct {
	storage    storage.Storage
	router     *chi.Mux
	logger     *zap.SugaredLogger
	backupFile string
	backupTime time.Duration
}

func (app *application) loadFromDisk() error {
	if err := disk.Load(app.backupFile, app.storage); err != nil {
		return err
	}
	return nil
}

func (app *application) flushToDisk() {

	ticker := time.NewTicker(app.backupTime)

	for range ticker.C {
		disk.Flush(app.backupFile, app.storage)
	}
}

func main() {

	cfg := config.NewServerConfig()

	s := inmemory.NewStorage()
	r := chi.NewRouter()
	l := logger.NewLogger()

	defer l.Sync()

	app := &application{
		storage:    s,
		router:     r,
		logger:     l,
		backupFile: cfg.FileStoragePath,
		backupTime: cfg.StoreInterval,
	}

	app.setRouters()

	if cfg.Restore {
		if err := app.loadFromDisk(); err != nil {
			log.Fatal(err)
			return
		}
	}

	go app.flushToDisk()

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: app.router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
		return
	}
}
