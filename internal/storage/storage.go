package storage

import (
	"github.com/h3ll0kitt1/observability/internal/models"
)

type Storage interface {
	Update(metric *models.Metrics)
	GetList() []*models.Metrics
	GetValue(metricType, metricName string) (string, bool)
}
