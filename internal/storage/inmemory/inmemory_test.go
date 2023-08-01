package inmemory

// import (
// 	"testing"

// 	"github.com/h3ll0kitt1/observability/internal/models"
// )

// func TestMemStorage_Update(t *testing.T) {

// 	gaugeValue := float64(2.12346)
// 	newGaugeValue := float64(3.123456)
// 	counterValue := int64(1)
// 	newCounterValue := int64(3)

// 	tests := []struct {
// 		name       string
// 		metric     models.Metrics
// 		wantValue  any
// 		wantStatus bool
// 	}{
// 		{
// 			name: "update existing gauge",
// 			metric: models.Metrics{
// 				ID:    "testGauge",
// 				MType: "gauge",
// 				Value: &gaugeValue,
// 			},
// 			wantValue:  float64(3.35802),
// 			wantStatus: true,
// 		},
// 		{
// 			name: "update existing counter",
// 			metric: models.Metrics{
// 				ID:    "testCounter",
// 				MType: "counter",
// 				Delta: &counterValue,
// 			},
// 			wantValue:  int64(2),
// 			wantStatus: true,
// 		},
// 		{
// 			name: "update new gauge",
// 			metric: models.Metrics{
// 				ID:    "newGauge",
// 				MType: "gauge",
// 				Value: &newGaugeValue,
// 			},
// 			wantValue:  float64(3.123456),
// 			wantStatus: true,
// 		},
// 		{
// 			name: "update new counter",
// 			metric: models.Metrics{
// 				ID:    "testCounter",
// 				MType: "counter",
// 				Delta: &newCounterValue,
// 			},
// 			wantValue:  int64(3),
// 			wantStatus: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			ms := NewStorage()
// 			ms.Counter.mem["testCounter"] = int64(1)
// 			ms.Gauge.mem["testGauge"] = float64(1.23456)

// 			ms.Update(&tt.metric)

// 			if tt.metric.MType == "counter" {
// 				if got, ok := ms.Counter.mem[tt.metric.ID]; got != tt.wantValue && ok != tt.wantStatus {
// 					t.Errorf("MemStorage_Update() = %v, want %v , wantStatus %v", got, tt.wantValue, tt.wantStatus)
// 				}
// 			}

// 			if tt.metric.MType == "gauge" {
// 				if got, ok := ms.Gauge.mem[tt.metric.ID]; got != tt.wantValue && ok != tt.wantStatus {
// 					t.Errorf("Update() = %v, want %v , wantStatus %v", got, tt.wantValue, tt.wantStatus)
// 				}
// 			}

// 		})
// 	}
// }

// func TestMemStorage_GetValue(t *testing.T) {

// 	tests := []struct {
// 		name       string
// 		mtype      string
// 		metricName string
// 		wantValue  string
// 		wantStatus bool
// 	}{
// 		{
// 			name:       "get existing gauge",
// 			mtype:      "gauge",
// 			metricName: "testGauge",
// 			wantValue:  "1",
// 			wantStatus: true,
// 		},
// 		{
// 			name:       "get existing counter",
// 			mtype:      "counter",
// 			metricName: "testCounter",
// 			wantValue:  "1",
// 			wantStatus: true,
// 		},
// 		{
// 			name:       "get wrong gauge",
// 			mtype:      "gauge",
// 			metricName: "wrongGauge",
// 			wantValue:  "0",
// 			wantStatus: false,
// 		},
// 		{
// 			name:       "get wrong counter",
// 			mtype:      "counter",
// 			metricName: "wrongCounter",
// 			wantValue:  "0",
// 			wantStatus: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			ms := NewStorage()
// 			ms.Counter.mem["testCounter"] = int64(1)
// 			ms.Gauge.mem["testGauge"] = float64(1.0)

// 			if gotValue, gotStatus := ms.GetValue(tt.mtype, tt.metricName); gotValue != tt.wantValue &&
// 				gotStatus != tt.wantStatus {
// 				t.Errorf("GetValue() = %v, want %v ", gotValue, tt.wantValue)
// 			}
// 		})
// 	}
// }
