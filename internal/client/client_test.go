package client

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestnewMetrics(t *testing.T) {

}

func TestnewCustomClient(t *testing.T) {

}

func TestMetrics_updateSpecificMemStats(t *testing.T) {
	mapGauge := make(map[string]float64)
	m := metrics{mapGauge, 0}
	m.updateSpecificMemStats()

	want := 27
	if got := len(m.mapGauge); got != want {
		t.Errorf("updateSpecificMemStats() = %v, want %v", got, want)
	}
}

func TestMetircs_updateRandomValue(t *testing.T) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	mapGauge := make(map[string]float64)
	m := metrics{mapGauge, 0}
	m.updateRandomValue(rng)

	want := float64(2.0)
	if got := m.mapGauge["Random"]; reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("updateRandomValue() = %v, want value %v", reflect.TypeOf(got), reflect.TypeOf(want))
	}
}

func TestdoRequestPOST(t *testing.T) {

}

func TestConstructURL(t *testing.T) {
	tests := []struct {
		name  string
		addr  string
		names string
		value any
		want  string
	}{
		{
			name:  "pass counter int64",
			addr:  "https//localhost:8080",
			names: "counter1",
			value: int64(1),
			want:  "https//localhost:8080/update/counter/counter1/1",
		},
		{
			name:  "pass gauge float64",
			addr:  "https//localhost:8080",
			names: "gauge1",
			value: 2.0,
			want:  "https//localhost:8080/update/gauge/gauge1/2.000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := constructURL(tt.addr, tt.names, tt.value); got != tt.want {
				t.Errorf("joinPath() = %v, want value %v", got, tt.want)
			}
		})
	}
}

func TestJoinPath(t *testing.T) {
	tests := []struct {
		name     string
		types    string
		names    string
		valueStr string
		want     string
	}{
		{
			name:     "pass counter int64",
			types:    "counter",
			names:    "counter1",
			valueStr: "1",
			want:     "/counter/counter1/1",
		},
		{
			name:     "pass gauge float64",
			types:    "gauge",
			names:    "gauge1",
			valueStr: "2.0",
			want:     "/gauge/gauge1/2.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinPath(tt.types, tt.names, tt.valueStr); got != tt.want {
				t.Errorf("joinPath() = %v, want value %v", got, tt.want)
			}
		})
	}
}
