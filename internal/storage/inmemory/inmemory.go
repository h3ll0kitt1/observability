package inmemory

import (
	"fmt"
	"strconv"
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

func (ms *MemStorage) Update(metric *models.Metrics) {
	switch metric.MType {
	case "counter":
		ms.Counter.Lock()
		ms.Counter.mem[metric.ID] += *(metric.Delta)
		ms.Counter.Unlock()
	case "gauge":
		ms.Gauge.Lock()
		ms.Gauge.mem[metric.ID] = *(metric.Value)
		ms.Gauge.Unlock()
	}
}

func (ms *MemStorage) GetList() []*models.Metrics {
	list := make([]*models.Metrics, 0)

	ms.Counter.Lock()
	for name, value := range ms.Counter.mem {
		metric := &models.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &value,
		}
		list = append(list, metric)
	}
	ms.Counter.Unlock()

	ms.Gauge.Lock()
	for name, value := range ms.Gauge.mem {
		metric := &models.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		}
		list = append(list, metric)
	}
	ms.Gauge.Unlock()

	return list
}

func (ms *MemStorage) GetValue(mtype, name string) (string, bool) {
	switch mtype {
	case "counter":
		value, ok := ms.Counter.mem[name]
		return fmt.Sprintf("%d", value), ok
	case "gauge":
		value, ok := ms.Gauge.mem[name]
		valueStr := strconv.FormatFloat(value, 'f', -1, 64)
		return valueStr, ok
	default:
		return "-1", false
	}
}
