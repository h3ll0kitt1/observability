package client

import (
	"math/rand"
	"testing"
)

func TestMetrics_updateSpecificMemStats(t *testing.T) {
	m := newMetrics()
	m.updateSpecificMemStats()

	want := 24
	if got := len(m.mapMetrics); got != want {
		t.Errorf("updateSpecificMemStats() = %v, want %v", got, want)
	}
}

func TestMetrics_updateRandomValue(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	m := newMetrics()
	m.updateRandomValue(rng)

	want := float64(81)
	key := metricKey{id: "RandomValue", mtype: "gauge"}
	if got := m.mapMetrics[key]; *got.Value != want {
		t.Errorf("updateRandomValue() = %v, want %v", *got.Value, want)
	}
}

func TestMetrics_updateCounterValue(t *testing.T) {
	m := newMetrics()
	m.updateCounterValue()

	want := int64(1)
	key := metricKey{id: "PollCount", mtype: "counter"}
	if got := m.mapMetrics[key]; *got.Delta != want {
		t.Errorf("updateCounterValue() = %v, want %v", *got.Delta, want)
	}
}

func TestGetFloat64(t *testing.T) {

	tests := []struct {
		name  string
		value any
		want  float64
	}{
		{
			name:  "test float64 type",
			value: float64(1),
			want:  float64(1),
		},
		{
			name:  "test float32 type",
			value: float32(1),
			want:  float64(1),
		},
		{
			name:  "test int64 type",
			value: int64(1),
			want:  float64(1),
		},
		{
			name:  "test int32 type",
			value: int32(1),
			want:  float64(1),
		},
		{
			name:  "test uint64 type",
			value: uint64(1),
			want:  float64(1),
		},
		{
			name:  "test uint32 type",
			value: uint32(1),
			want:  float64(1),
		},
		{
			name:  "test not specified type",
			value: "1",
			want:  float64(-1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFloat64(tt.value); got != tt.want {
				t.Errorf("getFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}
