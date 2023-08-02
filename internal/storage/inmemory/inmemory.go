package inmemory

import (
	"sync"

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

func (ms *MemStorage) Update(metric models.MetricsWithValue) {
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
}

func (ms *MemStorage) GetList() []*models.MetricsWithValue {
	ms.Counter.Lock()
	ms.Gauge.Lock()

	list := make([]*models.MetricsWithValue, 0)

	for name, value := range ms.Counter.mem {
		metric := &models.MetricsWithValue{
			ID:    name,
			MType: "counter",
			Delta: value,
		}
		list = append(list, metric)
	}

	for name, value := range ms.Gauge.mem {
		metric := &models.MetricsWithValue{
			ID:    name,
			MType: "gauge",
			Value: value,
		}
		list = append(list, metric)
	}
	ms.Counter.Unlock()
	ms.Gauge.Unlock()
	return list
}

func (ms *MemStorage) GetValue(metric models.MetricsWithValue) (models.MetricsWithValue, bool) {
	var status bool
	switch metric.MType {
	case "counter":
		value, ok := ms.Counter.mem[metric.ID]
		if ok {
			metric.Delta = value
		}
		status = ok
	case "gauge":
		value, ok := ms.Gauge.mem[metric.ID]
		if ok {
			metric.Value = value
		}
		status = ok
	}
	return metric, status
}
