package controller

import (
	"context"
	"time"

	"github.com/h3ll0kitt1/observability/internal/models"
)

type AsyncController struct {
	time    time.Duration
	backup  BackupStorage
	storage MainStorage
}

func (c *AsyncController) Load() error {
	list, err := c.backup.GetList(context.Background())
	if err != nil {
		return err
	}

	if err := c.storage.UpdateList(context.Background(), list); err != nil {
		return err
	}
	return nil
}

func (c *AsyncController) Run() {
	ticker := time.NewTicker(c.time)
	for range ticker.C {
		c.flush()
	}
}

func (c *AsyncController) Set(newMainStorage MainStorage) {
	c.storage = newMainStorage
}

func (c *AsyncController) Get(ctx context.Context, metric models.MetricsWithValue) (models.MetricsWithValue, error) {
	metric, err := c.storage.Get(ctx, metric)
	if err != nil {
		return metric, err
	}
	return metric, nil
}

func (c *AsyncController) GetList(ctx context.Context) ([]models.MetricsWithValue, error) {
	metrics, err := c.storage.GetList(ctx)
	if err != nil {
		return metrics, err
	}
	return metrics, nil
}

func (c *AsyncController) Update(ctx context.Context, metric models.MetricsWithValue) error {
	if err := c.storage.Update(ctx, metric); err != nil {
		return err
	}
	return nil
}

func (c *AsyncController) UpdateList(ctx context.Context, list []models.MetricsWithValue) error {
	if err := c.storage.UpdateList(ctx, list); err != nil {
		return err
	}
	return nil
}

func (c *AsyncController) Ping() error {
	return c.storage.Ping()
}

func (c *AsyncController) SetRetryCount(attempts int) {
	c.storage.SetRetryCount(attempts)
	c.backup.SetRetryCount(attempts)
}

func (c *AsyncController) SetRetryStartWaitTime(sleep time.Duration) {
	c.storage.SetRetryStartWaitTime(sleep)
	c.backup.SetRetryStartWaitTime(sleep)
}

func (c *AsyncController) SetRetryIncreaseWaitTime(delta time.Duration) {
	c.storage.SetRetryIncreaseWaitTime(delta)
	c.backup.SetRetryIncreaseWaitTime(delta)
}

func (c *AsyncController) flush() error {
	list, err := c.storage.GetList(context.Background())
	if err != nil {
		return err
	}

	if err := c.backup.UpdateList(context.Background(), list); err != nil {
		return err
	}
	return nil
}
