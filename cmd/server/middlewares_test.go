package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/h3ll0kitt1/observability/internal/logger"
	"github.com/h3ll0kitt1/observability/internal/models"
	"github.com/h3ll0kitt1/observability/internal/storage/inmemory"
)

func TestGzipper(t *testing.T) {

	s := inmemory.NewStorage()
	r := chi.NewRouter()
	l := logger.NewLogger()

	app := &application{
		storage: s,
		router:  r,
		logger:  l,
	}

	gaugeValue := float64(2.0)
	m1 := models.Metrics{
		ID:    "testGauge",
		MType: "gauge",
		Value: &gaugeValue,
	}
	s.Update(&m1)

	counterValue := int64(2)
	m2 := models.Metrics{
		ID:    "testCounter",
		MType: "counter",
		Delta: &counterValue,
	}
	s.Update(&m2)

	app.setRouters()

	ts := httptest.NewServer(app.router)
	defer ts.Close()

	testCases := []struct {
		name                string
		method              string
		contentEncodingGzip bool
		acceptEncodingGzip  bool
		body                string
		expectedCode        int
		expectedBody        string
	}{
		{
			name:                "application/json, content-encoding not gzip, accept-encoding not gzip",
			method:              http.MethodPost,
			contentEncodingGzip: false,
			acceptEncodingGzip:  false,
			body:                `{"id": "testGauge", "type": "gauge", "value": 3}`,
			expectedCode:        http.StatusOK,
			expectedBody:        `{"id": "testGauge", "type": "gauge", "value": 3}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = ts.URL + "/update/"

			req.SetHeader("Content-Type", "application/json")
			req.SetBody(tc.body)

			if tc.contentEncodingGzip {
				req.Header.Set("Content-Encoding", "gzip")
			}

			if !tc.acceptEncodingGzip {
				req.Header.Set("Accept-Encoding", "identity")
			}

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(resp.Body()))
			}
		})
	}
}
