package inmemory

import (
	"testing"

	"github.com/h3ll0kitt1/observability/internal/models"
)

func TestMemStorage_Update(t *testing.T) {
	tests := []struct {
		name       string
		metric     models.MetricsWithValue
		wantValue  any
		wantStatus bool
	}{
		{
			name: "update existing gauge",
			metric: models.MetricsWithValue{
				ID:    "testGauge",
				MType: "gauge",
				Value: float64(2.12346),
			},
			wantValue:  float64(3.35802),
			wantStatus: true,
		},
		{
			name: "update existing counter",
			metric: models.MetricsWithValue{
				ID:    "testCounter",
				MType: "counter",
				Delta: int64(1),
			},
			wantValue:  int64(2),
			wantStatus: true,
		},
		{
			name: "update new gauge",
			metric: models.MetricsWithValue{
				ID:    "newGauge",
				MType: "gauge",
				Value: float64(3.123456),
			},
			wantValue:  float64(3.123456),
			wantStatus: true,
		},
		{
			name: "update new counter",
			metric: models.MetricsWithValue{
				ID:    "testCounter",
				MType: "counter",
				Delta: int64(3),
			},
			wantValue:  int64(3),
			wantStatus: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ms := NewStorage()
			ms.Counter.mem["testCounter"] = int64(1)
			ms.Gauge.mem["testGauge"] = float64(1.23456)

			ms.Update(tt.metric)

			if tt.metric.MType == "counter" {
				if got, ok := ms.Counter.mem[tt.metric.ID]; got != tt.wantValue && ok != tt.wantStatus {
					t.Errorf("MemStorage_Update() = %v, want %v , wantStatus %v", got, tt.wantValue, tt.wantStatus)
				}
			}

			if tt.metric.MType == "gauge" {
				if got, ok := ms.Gauge.mem[tt.metric.ID]; got != tt.wantValue && ok != tt.wantStatus {
					t.Errorf("Update() = %v, want %v , wantStatus %v", got, tt.wantValue, tt.wantStatus)
				}
			}

		})
	}
}

func TestMemStorage_GetValue(t *testing.T) {

	tests := []struct {
		name       string
		mtype      string
		metric     models.MetricsWithValue
		wantValue  any
		wantStatus bool
	}{
		{
			name: "get existing gauge",
			metric: models.MetricsWithValue{
				ID:    "testGauge",
				MType: "gauge",
			},
			wantValue:  "1",
			wantStatus: true,
		},
		{
			name: "get existing counter",
			metric: models.MetricsWithValue{
				ID:    "testCounter",
				MType: "counter",
			},
			wantValue:  "1",
			wantStatus: true,
		},
		{
			name: "get wrong gauge",
			metric: models.MetricsWithValue{
				ID:    "wrongGauge",
				MType: "gauge",
			},
			wantValue:  "0",
			wantStatus: false,
		},
		{
			name: "get wrong counter",
			metric: models.MetricsWithValue{
				ID:    "wrongCounter",
				MType: "counter",
			},
			wantValue:  "0",
			wantStatus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ms := NewStorage()
			ms.Counter.mem["testCounter"] = int64(1)
			ms.Gauge.mem["testGauge"] = float64(1.0)

			if gotValue, gotStatus := ms.GetValue(tt.metric); gotValue != tt.wantValue &&
				gotStatus != tt.wantStatus {
				t.Errorf("GetValue() = %v, want %v ", gotValue, tt.wantValue)
			}
		})
	}
}
