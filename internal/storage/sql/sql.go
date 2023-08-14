package sql

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/h3ll0kitt1/observability/internal/config"
	"github.com/h3ll0kitt1/observability/internal/models"
)

type SQLStorage struct {
	db *sql.DB
}

func NewStorage(cfg *config.ServerConfig) *SQLStorage {
	db, err := sql.Open("pgx", cfg.Database)
	if err != nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `CREATE TABLE IF NOT EXISTS counter(
		metric_id varchar(512) primary key, 
		metric_value bigint not null)`
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating counter table", err)
		return nil
	}

	query = `CREATE TABLE IF NOT EXISTS gauge(
		metric_id varchar(512) primary key, 
		metric_value double precision not null)`
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating gauge table", err)
		return nil
	}

	return &SQLStorage{db: db}
}

func (s *SQLStorage) Update(metric models.MetricsWithValue) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	switch metric.MType {
	case "counter":
		_, err := s.db.ExecContext(ctx, `INSERT INTO counter (metric_id, metric_value) 
										VALUES ($1, $2) 
										ON CONFLICT (metric_id)
										DO
										UPDATE SET metric_value = EXCLUDED.metric_value;`, metric.ID, metric.Delta)
		if err != nil {
			log.Printf("Error %s when inserting in counter table", err)
		}
	case "gauge":
		_, err := s.db.ExecContext(ctx, `INSERT INTO gauge (metric_id, metric_value) 
										VALUES ($1, $2) 
										ON CONFLICT (metric_id) DO UPDATE 
										SET metric_value = EXCLUDED.metric_value;`, metric.ID, metric.Value)
		if err != nil {
			log.Printf("Error %s when inserting in gauge table", err)
		}
	}
}

func (s *SQLStorage) GetList() []*models.MetricsWithValue {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	list := make([]*models.MetricsWithValue, 0)

	rows, err := s.db.QueryContext(ctx, "SELECT metric_id, metric_value FROM counter")
	if err != nil {
		//return nil, err
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var metric models.MetricsWithValue
		metric.MType = "counter"

		err = rows.Scan(&metric.ID, &metric.Delta)
		if err != nil {
			//return nil, err
			return nil
		}
		list = append(list, &metric)
	}

	err = rows.Err()
	if err != nil {
		//return nil, err
		return nil
	}

	rows, err = s.db.QueryContext(ctx, "SELECT metric_id, metric_value FROM gauge")
	if err != nil {
		//return nil, err
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var metric models.MetricsWithValue
		metric.MType = "gauge"

		err = rows.Scan(&metric.ID, &metric.Value)
		if err != nil {
			//return nil, err
			return nil
		}
		list = append(list, &metric)
	}

	err = rows.Err()
	if err != nil {
		//return nil, err
		return nil
	}

	return list
}

func (s *SQLStorage) GetValue(metric models.MetricsWithValue) (models.MetricsWithValue, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	switch metric.MType {

	case "counter":
		var value int64

		row := s.db.QueryRowContext(ctx, "SELECT metric_value FROM counter WHERE metric_id = $1", metric.ID)

		if err := row.Scan(&value); err != nil {
			return metric, false
		}
		metric.Delta = value

	case "gauge":
		var value float64

		row := s.db.QueryRowContext(ctx, "SELECT metric_value FROM gauge WHERE metric_id = $1", metric.ID)

		if err := row.Scan(&value); err != nil {
			return metric, false
		}
		metric.Value = value
	}

	return metric, true
}

func (s *SQLStorage) Ping() bool {
	if err := s.db.Ping(); err != nil {
		return false
	}
	return true
}
