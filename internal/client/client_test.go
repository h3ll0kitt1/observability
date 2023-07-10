package client

import (
	"math/rand"
	"testing"
)

func TestMetrics_updateSpecificMemStats(t *testing.T) {
	m := newMetrics()
	m.updateSpecificMemStats()

	want := 27
	if got := len(m.mapMetrics["gauge"]); got != want {
		t.Errorf("updateSpecificMemStats() = %v, want %v", got, want)
	}
}

func TestMetrics_updateRandomValue(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	m := newMetrics()
	m.updateRandomValue(rng)

	want := "81.000000"
	if got := m.mapMetrics["gauge"]["Random"]; got != want {
		t.Errorf("updateRandomValue() = %v, want %v", got, want)
	}
}

func TestMetrics_updateCounterValue(t *testing.T) {
	m := newMetrics()
	m.updateCounterValue()

	want := "1"
	if got := m.mapMetrics["counter"]["Counter"]; got != want {
		t.Errorf("updateCounterValue() = %v, want %v", got, want)
	}
}

func TestConvertToString(t *testing.T) {

	tests := []struct {
		name  string
		value any
		want  string
	}{
		{
			name:  "test counter type",
			value: uint64(1),
			want:  "1.000000",
		},
		{
			name:  "test gauge type",
			value: float64(2.55500),
			want:  "2.555000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertToString(tt.value); got != tt.want {
				t.Errorf("convertToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
