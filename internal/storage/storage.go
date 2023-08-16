package storage

import (
	"context"

	"github.com/h3ll0kitt1/observability/internal/models"
)

type Storage interface {
	Get(ctx context.Context, metric models.MetricsWithValue) (models.MetricsWithValue, error)
	GetList(ctx context.Context) ([]models.MetricsWithValue, error)

	Update(ctx context.Context, metric models.MetricsWithValue) error
	UpdateList(ctx context.Context, list []models.MetricsWithValue) error

	Ping() error
}
