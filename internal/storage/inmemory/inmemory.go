package inmemory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/h3ll0kitt1/observability/internal/models"
)

type MemStorage struct {
	Counter *MemCounter
	Gauge   *MemGauge
}

type MemCounter struct {
	mem map[string]int64
	sync.Mutex
}

type MemGauge struct {
	mem map[string]float64
	sync.Mutex
}

func NewStorage() *MemStorage {
	MemCounter := NewMemCounter()
	MemGauge := NewMemGauge()
	return &MemStorage{
		Counter: MemCounter,
		Gauge:   MemGauge,
	}
}

func NewMemCounter() *MemCounter {
	var mc MemCounter
	mc.mem = make(map[string]int64)
	return &mc
}

func NewMemGauge() *MemGauge {
	var mg MemGauge
	mg.mem = make(map[string]float64)
	return &mg
}

func (ms *MemStorage) Get(ctx context.Context, metric models.MetricsWithValue) (models.MetricsWithValue, error) {
	var status bool

	switch metric.MType {
	case "counter":
		ms.Counter.Lock()
		value, ok := ms.Counter.mem[metric.ID]
		if ok {
			metric.Delta = value
		}
		status = ok
		ms.Counter.Unlock()
	case "gauge":
		ms.Gauge.Lock()
		value, ok := ms.Gauge.mem[metric.ID]
		if ok {
			metric.Value = value
		}
		status = ok
		ms.Gauge.Unlock()
	}

	if !status {
		return metric, errors.New("unknown metric")
	}
	return metric, nil
}

func (ms *MemStorage) GetList(ctx context.Context) ([]models.MetricsWithValue, error) {
	ms.Counter.Lock()
	ms.Gauge.Lock()

	list := make([]models.MetricsWithValue, 0)

	for name, value := range ms.Counter.mem {
		metric := models.MetricsWithValue{
			ID:    name,
			MType: "counter",
			Delta: value,
		}
		list = append(list, metric)
	}

	for name, value := range ms.Gauge.mem {
		metric := models.MetricsWithValue{
			ID:    name,
			MType: "gauge",
			Value: value,
		}
		list = append(list, metric)
	}
	ms.Counter.Unlock()
	ms.Gauge.Unlock()

	return list, nil
}

func (ms *MemStorage) Update(ctx context.Context, metric models.MetricsWithValue) error {
	switch metric.MType {
	case "counter":
		ms.Counter.Lock()
		ms.Counter.mem[metric.ID] += metric.Delta
		ms.Counter.Unlock()
	case "gauge":
		ms.Gauge.Lock()
		ms.Gauge.mem[metric.ID] = metric.Value
		ms.Gauge.Unlock()
	}
	return nil
}

func (ms *MemStorage) UpdateList(ctx context.Context, list []models.MetricsWithValue) error {
	for _, metric := range list {
		ms.Update(ctx, metric)
	}
	return nil
}

func (ms *MemStorage) Ping() error { return nil }

func (ms *MemStorage) SetRetryCount(attempts int) {}

func (ms *MemStorage) SetRetryStartWaitTime(sleep time.Duration) {}

func (ms *MemStorage) SetRetryIncreaseWaitTime(delta time.Duration) {}
