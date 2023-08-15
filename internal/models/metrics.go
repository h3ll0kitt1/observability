package models

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type MetricsWithValue struct {
	ID    string
	MType string
	Delta int64
	Value float64
}

func ToMetricWithValue(metric Metrics) MetricsWithValue {
	var m MetricsWithValue
	m.ID = metric.ID
	m.MType = metric.MType
	if metric.Delta != nil {
		m.Delta = *(metric.Delta)
	}
	if metric.Value != nil {
		m.Value = *(metric.Value)
	}
	return m
}

func ToMetric(metric MetricsWithValue) Metrics {
	var m Metrics
	m.ID = metric.ID
	m.MType = metric.MType
	switch m.MType {
	case "counter":
		m.Delta = &metric.Delta
	case "gauge":
		m.Value = &metric.Value
	}
	return m
}
