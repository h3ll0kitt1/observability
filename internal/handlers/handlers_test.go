package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/h3ll0kitt1/observability/internal/storage/inmemory"
)

func TestUpdate(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name: "test no metric name",
			want: want{
				code: http.StatusNotFound,
			},
			request: "/update/counter",
		},
		{
			name: "test wrong metric type",
			want: want{
				code: http.StatusBadRequest,
			},
			request: "/update/wrongtype",
		},
		{
			name: "test no metric value",
			want: want{
				code: http.StatusBadRequest,
			},
			request: "/update/gauge/counter1",
		},
		{
			name: "test wrong metric value",
			want: want{
				code: http.StatusBadRequest,
			},
			request: "/update/counter/counter1/2.0",
		},
		{
			name: "test right counter request",
			want: want{
				code: http.StatusOK,
			},
			request: "/update/gauge/counter1/2",
		},
		{
			name: "test right gauge request",
			want: want{
				code: http.StatusOK,
			},
			request: "/update/gauge/gauge1/2.0",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var ms inmemory.MemStorage
			ms.Counter = make(map[string]int64)
			ms.Gauge = make(map[string]float64)

			request := httptest.NewRequest(http.MethodPost, test.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(Update(&ms))
			h(w, request)

			res := w.Result()
			if res.StatusCode != test.want.code {
				t.Errorf("Update() = %v, want %v", res.StatusCode, test.want.code)
			}
			defer res.Body.Close()
		})
	}

}

func TestUpdateNotSpecified(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "test /update",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/UpdateNotSpecified", nil)
			w := httptest.NewRecorder()
			UpdateNotSpecified(w, request)

			res := w.Result()
			if res.StatusCode != test.want.code {
				t.Errorf("UpdateNotSpecified() = %v, want %v", res.StatusCode, test.want.code)
			}
			defer res.Body.Close()
		})
	}
}

func TestValidateStringIsInt64(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		valueWant  int64
		statusWant bool
	}{
		{name: "pass int64", value: "2", valueWant: 2, statusWant: true},
		{name: "pass float64", value: "2.1", valueWant: -1, statusWant: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := validateStringIsInt64(tt.value); got != tt.valueWant &&
				ok != tt.statusWant {
				t.Errorf("validateStringIsInt64() = %v, want value %v, want status %v", got, tt.valueWant, tt.statusWant)
			}
		})
	}
}

func TestValidateStringIsFloat64(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		valueWant  float64
		statusWant bool
	}{
		{name: "pass float64 with dot", value: "2.1", valueWant: 2.1, statusWant: true},
		{name: "pass float64 without dot", value: "2", valueWant: 2, statusWant: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := validateStringIsFloat64(tt.value); got != tt.valueWant &&
				ok != tt.statusWant {
				t.Errorf("validateStringIsFloat64() = %v, want value %v, want status %v", got, tt.valueWant, tt.statusWant)
			}
		})
	}
}
