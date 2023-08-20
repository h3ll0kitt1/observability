package controller

import (
	"context"
	"time"

	"github.com/h3ll0kitt1/observability/internal/models"
)

type SyncController struct {
	backup  BackupStorage
	storage MainStorage
}

func (c *SyncController) Load() error {
	list, err := c.backup.GetList(context.Background())
	if err != nil {
		return err
	}

	if err := c.storage.UpdateList(context.Background(), list); err != nil {
		return err
	}
	return nil
}

func (c *SyncController) Run() {
}

func (c *SyncController) Set(newMainStorage MainStorage) {
	c.storage = newMainStorage
}

func (c *SyncController) Get(ctx context.Context, metric models.MetricsWithValue) (models.MetricsWithValue, error) {
	metric, err := c.storage.Get(ctx, metric)
	if err != nil {
		return metric, err
	}
	return metric, nil
}

func (c *SyncController) GetList(ctx context.Context) ([]models.MetricsWithValue, error) {
	metrics, err := c.storage.GetList(ctx)
	if err != nil {
		return metrics, err
	}
	return metrics, nil
}

func (c *SyncController) Update(ctx context.Context, metric models.MetricsWithValue) error {
	if err := c.storage.Update(ctx, metric); err != nil {
		return err
	}
	return c.flush()
}

func (c *SyncController) UpdateList(ctx context.Context, list []models.MetricsWithValue) error {
	if err := c.storage.UpdateList(ctx, list); err != nil {
		return err
	}
	return c.flush()
}

func (c *SyncController) Ping() error {
	return c.storage.Ping()
}

func (c *SyncController) SetRetryCount(attempts int) {
	c.storage.SetRetryCount(attempts)
	c.backup.SetRetryCount(attempts)
}

func (c *SyncController) SetRetryStartWaitTime(sleep time.Duration) {
	c.storage.SetRetryStartWaitTime(sleep)
	c.backup.SetRetryStartWaitTime(sleep)
}

func (c *SyncController) SetRetryIncreaseWaitTime(delta time.Duration) {
	c.storage.SetRetryIncreaseWaitTime(delta)
	c.backup.SetRetryIncreaseWaitTime(delta)
}

func (c *SyncController) flush() error {
	list, err := c.storage.GetList(context.Background())
	if err != nil {
		return err
	}

	if err := c.backup.UpdateList(context.Background(), list); err != nil {
		return err
	}
	return nil
}
