package sql

import (
	"context"
	"database/sql"
	"errors"
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
		log.Fatalf("Error %s open database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `CREATE TABLE IF NOT EXISTS counter(
		metric_id varchar(512) primary key, 
		metric_value bigint not null)`
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalf("Error %s when creating counter table", err)
	}

	query = `CREATE TABLE IF NOT EXISTS gauge(
		metric_id varchar(512) primary key, 
		metric_value double precision not null)`
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalf("Error %s when creating gauge table", err)
	}
	return &SQLStorage{db: db}
}

func (s *SQLStorage) Get(ctx context.Context, metric models.MetricsWithValue) (models.MetricsWithValue, error) {
	switch metric.MType {
	case "counter":
		var value int64
		row := s.db.QueryRowContext(ctx, "SELECT metric_value FROM counter WHERE metric_id = $1", metric.ID)

		if err := row.Scan(&value); err != nil {
			return metric, errors.New("unknown metric name")
		}
		metric.Delta = value

	case "gauge":
		var value float64
		row := s.db.QueryRowContext(ctx, "SELECT metric_value FROM gauge WHERE metric_id = $1", metric.ID)

		if err := row.Scan(&value); err != nil {
			return metric, errors.New("unknown metric name")
		}
		metric.Value = value
	}
	return metric, nil
}

func (s *SQLStorage) GetList(ctx context.Context) ([]models.MetricsWithValue, error) {
	list := make([]models.MetricsWithValue, 0)

	rows, err := s.db.QueryContext(ctx, "SELECT metric_id, metric_value FROM counter")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var metric models.MetricsWithValue
		metric.MType = "counter"

		err = rows.Scan(&metric.ID, &metric.Delta)
		if err != nil {
			return nil, err
		}
		list = append(list, metric)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	rows, err = s.db.QueryContext(ctx, "SELECT metric_id, metric_value FROM gauge")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var metric models.MetricsWithValue
		metric.MType = "gauge"

		err = rows.Scan(&metric.ID, &metric.Value)
		if err != nil {
			return nil, err
		}
		list = append(list, metric)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *SQLStorage) Update(ctx context.Context, metric models.MetricsWithValue) error {
	switch metric.MType {
	case "counter":
		_, err := s.db.ExecContext(ctx,
			`INSERT INTO counter (metric_id, metric_value) 
			VALUES ($1, $2) 
			ON CONFLICT (metric_id) DO UPDATE 
			SET metric_value = EXCLUDED.metric_value + counter.metric_value;`, metric.ID, metric.Delta)
		if err != nil {
			return err
		}

	case "gauge":
		_, err := s.db.ExecContext(ctx,
			`INSERT INTO gauge (metric_id, metric_value) 
			VALUES ($1, $2) 
			ON CONFLICT (metric_id) DO UPDATE 
			SET metric_value = EXCLUDED.metric_value;`, metric.ID, metric.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SQLStorage) UpdateList(ctx context.Context, list []models.MetricsWithValue) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, metric := range list {
		switch metric.MType {
		case "counter":
			_, err := tx.ExecContext(ctx,
				`INSERT INTO counter (metric_id, metric_value) 
				VALUES ($1, $2) 
				ON CONFLICT (metric_id) DO UPDATE 
				SET metric_value = EXCLUDED.metric_value + counter.metric_value;`, metric.ID, metric.Delta)
			if err != nil {
				return err
			}

		case "gauge":
			_, err := tx.ExecContext(ctx,
				`INSERT INTO gauge (metric_id, metric_value) 
				VALUES ($1, $2) 
				ON CONFLICT (metric_id) DO UPDATE 
				SET metric_value = EXCLUDED.metric_value;`, metric.ID, metric.Value)

			if err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

func (s *SQLStorage) Ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}
	return nil
}
