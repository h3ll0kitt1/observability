package storage

import (
	"github.com/h3ll0kitt1/observability/internal/models"
)

type Storage interface {
	Update(metric models.MetricsWithValue)
	UpdateList(metric []models.MetricsWithValue)
	GetValue(metric models.MetricsWithValue) (models.MetricsWithValue, bool)
	GetList() []*models.MetricsWithValue
	Ping() bool
}
