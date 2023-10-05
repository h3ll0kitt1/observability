package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/h3ll0kitt1/observability/internal/config"
	"github.com/h3ll0kitt1/observability/internal/controller"
	"github.com/h3ll0kitt1/observability/internal/logger"
	"github.com/h3ll0kitt1/observability/internal/storage/sql"
)

type application struct {
	config         *config.ServerConfig
	storageManager controller.StorageManager
	router         *chi.Mux
	logger         *zap.SugaredLogger
}

func main() {

	cfg := config.NewServerConfig()
	cfg.Parse()

	sm := controller.NewStorageManager(cfg)
	sm.SetRetryCount(3)
	sm.SetRetryStartWaitTime(1)
	sm.SetRetryIncreaseWaitTime(2)

	if cfg.Database != "" {
		db, err := sql.NewStorage(cfg)
		if err != nil {
			log.Fatalf("Error %s open database", err)
		}
		sm.Set(db)
	}

	if cfg.Restore {
		if err := sm.Load(); err != nil {
			log.Fatalf("Error %s loading from disk", err)
		}
	}

	r := chi.NewRouter()
	l := logger.NewLogger()
	defer l.Sync()

	app := &application{
		config:         cfg,
		storageManager: sm,
		router:         r,
		logger:         l,
	}
	app.setRouters()

	go app.storageManager.Run()

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: app.router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error %s launching server", err)
	}
}
