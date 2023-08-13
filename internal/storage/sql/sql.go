package sql

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/h3ll0kitt1/observability/internal/config"
	"github.com/h3ll0kitt1/observability/internal/models"
)

type SQLStorage struct {
	DB *sql.DB
}

func NewStorage(cfg *config.ServerConfig) *SQLStorage {
	db, err := sql.Open("pgx", cfg.Database)
	if err != nil {
		return nil
	}
	return &SQLStorage{DB: db}
}

func (s *SQLStorage) Update(metric models.MetricsWithValue) {

}

func (s *SQLStorage) GetList() []*models.MetricsWithValue {
	list := make([]*models.MetricsWithValue, 0)
	return list
}

func (s *SQLStorage) GetValue(metric models.MetricsWithValue) (models.MetricsWithValue, bool) {
	return models.MetricsWithValue{}, false
}

func (s *SQLStorage) Ping() bool {
	if err := s.DB.Ping(); err != nil {
		return false
	}
	return true
}
