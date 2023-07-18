package inmemory

import (
	"testing"
)

func TestMemStorage_Update(t *testing.T) {

	tests := []struct {
		name        string
		metricName  string
		metricType  string
		metricValue any
		wantValue   any
		wantStatus  bool
	}{
		{
			name:        "update existing gauge",
			metricName:  "testGauge",
			metricType:  "gauge",
			metricValue: float64(2.0),
			wantValue:   float64(2.0),
			wantStatus:  true,
		},
		{
			name:        "update existing counter",
			metricName:  "testCounter",
			metricType:  "counter",
			metricValue: int64(1.0),
			wantValue:   int64(2.0),
			wantStatus:  true,
		},
		{
			name:        "update new gauge",
			metricName:  "newGauge",
			metricType:  "gauge",
			metricValue: float64(3.0),
			wantValue:   float64(3.0),
			wantStatus:  true,
		},
		{
			name:        "update new counter",
			metricName:  "testCounter",
			metricType:  "counter",
			metricValue: int64(3.0),
			wantValue:   int64(3.0),
			wantStatus:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ms := NewStorage()
			ms.Counter.mem["testCounter"] = int64(1)
			ms.Gauge.mem["testGauge"] = float64(1.0)

			ms.Update(tt.metricName, tt.metricValue)

			if tt.metricType == "counter" {
				if got, ok := ms.Counter.mem[tt.metricName]; got != tt.wantValue && ok != tt.wantStatus {
					t.Errorf("MemStorage_Update() = %v, want %v , wantStatus %v", got, tt.wantValue, tt.wantStatus)
				}
			}

			if tt.metricType == "gauge" {
				if got, ok := ms.Gauge.mem[tt.metricName]; got != tt.wantValue && ok != tt.wantStatus {
					t.Errorf("Update() = %v, want %v , wantStatus %v", got, tt.wantValue, tt.wantStatus)
				}
			}

		})
	}
}

func TestMemStorage_GetList(t *testing.T) {

	tests := []struct {
		name string
		want string
	}{
		{
			name: "list all",
			want: "testCounter : 1\ntestGauge : 1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ms := NewStorage()
			ms.Counter.mem["testCounter"] = int64(1)
			ms.Gauge.mem["testGauge"] = float64(1.0)

			if got := ms.GetList(); got != tt.want {
				t.Errorf("GetList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetValue(t *testing.T) {

	tests := []struct {
		name       string
		mtype      string
		metricName string
		wantValue  string
		wantStatus bool
	}{
		{
			name:       "get existing gauge",
			mtype:      "gauge",
			metricName: "testGauge",
			wantValue:  "1",
			wantStatus: true,
		},
		{
			name:       "get existing counter",
			mtype:      "counter",
			metricName: "testCounter",
			wantValue:  "1",
			wantStatus: true,
		},
		{
			name:       "get wrong gauge",
			mtype:      "gauge",
			metricName: "wrongGauge",
			wantValue:  "0",
			wantStatus: false,
		},
		{
			name:       "get wrong counter",
			mtype:      "counter",
			metricName: "wrongCounter",
			wantValue:  "0",
			wantStatus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ms := NewStorage()
			ms.Counter.mem["testCounter"] = int64(1)
			ms.Gauge.mem["testGauge"] = float64(1.0)

			if gotValue, gotStatus := ms.GetValue(tt.mtype, tt.metricName); gotValue != tt.wantValue &&
				gotStatus != tt.wantStatus {
				t.Errorf("GetValue() = %v, want %v ", gotValue, tt.wantValue)
			}
		})
	}
}
