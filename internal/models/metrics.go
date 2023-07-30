package models

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func NewMetric(mtype string, name string, mvalue any) Metrics {
	var m Metrics
	m.ID = name
	m.MType = mtype

	if mtype == "counter" {
		v := mvalue.(int64)
		m.Delta = &v
	}

	if mtype == "gauge" {
		v := mvalue.(float64)
		m.Value = &v
	}

	return m
}
