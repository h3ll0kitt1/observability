package client

import (
	"math/rand"
	"testing"
)

func TestMetrics_updateSpecificMemStats(t *testing.T) {
	m := newMetrics()
	m.updateSpecificMemStats()

	searched := []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle",
		"HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
		"StackSys", "Sys", "TotalAlloc"}

	want := true
	for _, id := range searched {
		if _, ok := m.mapMetrics.metrics[metricKey{id: id, mtype: "gauge"}]; ok != want {
			t.Errorf("updateSpecificMemStats() = %v, want %v", ok, want)
		}
	}
}

func TestMetrics_updateMemoryCPUInfo(t *testing.T) {
	m := newMetrics()
	m.updateMemoryCPUInfo()

	searched := []string{"Total", "Free", "UsedPercent"}

	want := true
	for _, id := range searched {
		if _, ok := m.mapMetrics.metrics[metricKey{id: id, mtype: "gauge"}]; ok != want {
			t.Errorf("updateSpecificMemStats() = %v, want %v", ok, want)
		}
	}
}

func TestMetrics_updateRandomValue(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	m := newMetrics()
	m.updateRandomValue(rng)

	want := float64(81)
	key := metricKey{id: "RandomValue", mtype: "gauge"}
	if got := m.mapMetrics.metrics[key]; *got.Value != want {
		t.Errorf("updateRandomValue() = %v, want %v", *got.Value, want)
	}
}

func TestMetrics_updateCounterValue(t *testing.T) {
	m := newMetrics()
	m.updateCounterValue()

	want := int64(1)
	key := metricKey{id: "PollCount", mtype: "counter"}
	if got := m.mapMetrics.metrics[key]; *got.Delta != want {
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
