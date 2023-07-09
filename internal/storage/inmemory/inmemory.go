package inmemory

import (
	"fmt"
)

type MemStorage struct {
	Counter map[string]int64
	Gauge   map[string]float64
}

func NewStorage() *MemStorage {
	var ms MemStorage
	ms.Counter = make(map[string]int64)
	ms.Gauge = make(map[string]float64)
	return &ms
}

func (ms *MemStorage) Update(metricName string, metricValue any) {
	switch mv := metricValue.(type) {
	case int64:
		ms.Counter[metricName] += mv
	case float64:
		ms.Gauge[metricName] = mv
	}
}

func (ms MemStorage) GetList() string {
	list := ""
	for name, value := range ms.Counter {
		list += name + fmt.Sprintf(" : %d\n", value)
	}
	for name, value := range ms.Gauge {
		list += name + fmt.Sprintf(" : %f\n", value)
	}
	return list
}

func (ms MemStorage) GetValue(mtype, name string) (string, bool) {
	switch mtype {
	case "counter":
		value, ok := ms.Counter[name]
		return fmt.Sprintf("%d", value), ok
	case "gauge":
		value, ok := ms.Gauge[name]
		return fmt.Sprintf("%f", value), ok
	default:
		return "-1", false
	}
}
