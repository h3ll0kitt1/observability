package models

import (
	"encoding/json"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m Metrics) Convert2JSON() ([]byte, error) {
	res, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return res, nil
}
