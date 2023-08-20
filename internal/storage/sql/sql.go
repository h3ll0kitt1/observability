package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/h3ll0kitt1/observability/internal/config"
	"github.com/h3ll0kitt1/observability/internal/models"
)

type SQLStorage struct {
	db      *sql.DB
	retrier retrier
}

type retrier struct {
	attempts   int
	time       time.Duration
	delta      time.Duration
	errToRetry map[string]bool
}

func NewStorage(cfg *config.ServerConfig) (*SQLStorage, error) {
	db, err := sql.Open("pgx", cfg.Database)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `CREATE TABLE IF NOT EXISTS counter(
		metric_id varchar(512) primary key, 
		metric_value bigint not null)`
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}

	query = `CREATE TABLE IF NOT EXISTS gauge(
		metric_id varchar(512) primary key, 
		metric_value double precision not null)`
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}

	errToRetry := map[string]bool{
		pgerrcode.ConnectionException:                     true,
		pgerrcode.ConnectionDoesNotExist:                  true,
		pgerrcode.ConnectionFailure:                       true,
		pgerrcode.SQLClientUnableToEstablishSQLConnection: true,
	}
	return &SQLStorage{
		db:      db,
		retrier: retrier{attempts: 1, time: 0, delta: 0, errToRetry: errToRetry},
	}, nil
}

func (s *SQLStorage) Get(ctx context.Context, metric models.MetricsWithValue) (models.MetricsWithValue, error) {
	get := func() (models.MetricsWithValue, error) { return s.get(ctx, metric) }
	metric, err := s.retryWithMetric(get)
	if err != nil {
		return metric, err
	}
	return metric, nil
}

func (s *SQLStorage) GetList(ctx context.Context) ([]models.MetricsWithValue, error) {
	getList := func() ([]models.MetricsWithValue, error) { return s.getList(ctx) }
	list, err := s.retryWithMetrics(getList)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *SQLStorage) Update(ctx context.Context, metric models.MetricsWithValue) error {
	update := func() error { return s.update(ctx, metric) }
	if err := s.retry(update); err != nil {
		return err
	}
	return nil
}

func (s *SQLStorage) UpdateList(ctx context.Context, list []models.MetricsWithValue) error {
	updateList := func() error { return s.updateList(ctx, list) }
	if err := s.retry(updateList); err != nil {
		return err
	}
	return nil
}

func (s *SQLStorage) Ping() error {
	ping := func() error { return s.ping() }
	if err := s.retry(ping); err != nil {
		return err
	}
	return nil
}

func (s *SQLStorage) SetRetryCount(attempts int) {
	s.retrier.attempts = attempts
}

func (s *SQLStorage) SetRetryStartWaitTime(sleep time.Duration) {
	s.retrier.time = sleep
}

func (s *SQLStorage) SetRetryIncreaseWaitTime(delta time.Duration) {
	s.retrier.delta = delta
}

func (s *SQLStorage) retry(f func() error) error {
	var err error
	for i := 0; i <= s.retrier.attempts; i++ {
		if i > 0 {
			time.Sleep(s.retrier.time)
			s.retrier.time += s.retrier.delta
		}
		err = f()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && !s.retrier.errToRetry[pgErr.Code] {
				break
			}
		}
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("after %d attempts, last error: %s", s.retrier.attempts, err)
}

func (s *SQLStorage) retryWithMetric(f func() (models.MetricsWithValue, error)) (models.MetricsWithValue, error) {
	var (
		err    error
		metric models.MetricsWithValue
	)
	for i := 0; i <= s.retrier.attempts; i++ {
		if i > 0 {
			time.Sleep(s.retrier.time)
			s.retrier.time += s.retrier.delta
		}
		metric, err = f()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && !s.retrier.errToRetry[pgErr.Code] {
				break
			}
		}
		if err == nil {
			return metric, nil
		}
	}
	return metric, fmt.Errorf("after %d attempts, last error: %s", s.retrier.attempts, err)
}

func (s *SQLStorage) retryWithMetrics(f func() ([]models.MetricsWithValue, error)) ([]models.MetricsWithValue, error) {
	var err error
	for i := 0; i <= s.retrier.attempts; i++ {
		if i > 0 {
			time.Sleep(s.retrier.time)
			s.retrier.time += s.retrier.delta
		}
		list, err := f()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && !s.retrier.errToRetry[pgErr.Code] {
				break
			}
		}
		if err == nil {
			return list, nil
		}
	}
	return nil, fmt.Errorf("after %d attempts, last error: %s", s.retrier.attempts, err)
}

func (s *SQLStorage) get(ctx context.Context, metric models.MetricsWithValue) (models.MetricsWithValue, error) {
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

func (s *SQLStorage) getList(ctx context.Context) ([]models.MetricsWithValue, error) {
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

func (s *SQLStorage) update(ctx context.Context, metric models.MetricsWithValue) error {
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

func (s *SQLStorage) updateList(ctx context.Context, list []models.MetricsWithValue) error {
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

func (s *SQLStorage) ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}
	return nil
}
