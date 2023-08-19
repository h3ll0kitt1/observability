package controller

import (
	"context"

	"github.com/h3ll0kitt1/observability/internal/models"
)

type MainStorage interface {
	Get(ctx context.Context, metric models.MetricsWithValue) (models.MetricsWithValue, error)
	Update(ctx context.Context, metric models.MetricsWithValue) error
	Ping() error

	BackupStorage
}

type BackupStorage interface {
	GetList(ctx context.Context) ([]models.MetricsWithValue, error)
	UpdateList(ctx context.Context, list []models.MetricsWithValue) error
}
