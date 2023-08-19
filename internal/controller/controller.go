package controller

import (
	"github.com/h3ll0kitt1/observability/internal/config"
	"github.com/h3ll0kitt1/observability/internal/storage/file"
	"github.com/h3ll0kitt1/observability/internal/storage/inmemory"
)

type StorageManager interface {
	Load() error
	Run()
	Set(MainStorage)

	MainStorage
}

func NewStorageManager(cfg *config.ServerConfig) StorageManager {
	s := inmemory.NewStorage()
	b := file.NewStorage(cfg.FileStoragePath)

	if cfg.StoreInterval == 0 {
		return &SyncController{
			storage: s,
			backup:  b,
		}
	}

	return &AsyncController{
		time:    cfg.StoreInterval,
		storage: s,
		backup:  b,
	}
}
