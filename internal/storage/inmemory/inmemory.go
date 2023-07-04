package inmemory

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
	switch metricValue.(type) {
	case int64:
		ms.Counter[metricName] += metricValue.(int64)
	case float64:
		ms.Gauge[metricName] = metricValue.(float64)
	}
}
