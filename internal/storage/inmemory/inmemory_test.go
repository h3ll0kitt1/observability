package inmemory

import (
	"context"
	"testing"
	"time"

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

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			ms := NewStorage()
			ms.Counter.mem["testCounter"] = int64(1)
			ms.Gauge.mem["testGauge"] = float64(1.23456)

			ms.Update(ctx, tt.metric)

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
		name      string
		mtype     string
		metric    models.MetricsWithValue
		wantValue any
	}{
		{
			name: "get existing gauge",
			metric: models.MetricsWithValue{
				ID:    "testGauge",
				MType: "gauge",
			},
			wantValue: float64(1),
		},
		{
			name: "get existing counter",
			metric: models.MetricsWithValue{
				ID:    "testCounter",
				MType: "counter",
			},
			wantValue: int64(1),
		},
		{
			name: "get wrong gauge",
			metric: models.MetricsWithValue{
				ID:    "wrongGauge",
				MType: "gauge",
			},
			wantValue: float64(0),
		},
		{
			name: "get wrong counter",
			metric: models.MetricsWithValue{
				ID:    "wrongCounter",
				MType: "counter",
			},
			wantValue: int64(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			ms := NewStorage()
			ms.Counter.mem["testCounter"] = int64(1)
			ms.Gauge.mem["testGauge"] = float64(1.0)

			gotMetric, _ := ms.Get(ctx, tt.metric)
			if gotMetric.MType == "counter" && gotMetric.Delta != tt.wantValue {
				t.Errorf("GetValue() = %v, want %v ", gotMetric.Delta, tt.wantValue)
			}

			if gotMetric.MType == "gauge" && gotMetric.Value != tt.wantValue {
				t.Errorf("GetValue() = %v, want %v ", gotMetric.Value, tt.wantValue)
			}

		})
	}
}
