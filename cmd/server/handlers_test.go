package main

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-resty/resty/v2"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"

// 	"github.com/h3ll0kitt1/observability/internal/config"
// 	"github.com/h3ll0kitt1/observability/internal/logger"
// 	"github.com/h3ll0kitt1/observability/internal/mocks"
// 	"github.com/h3ll0kitt1/observability/internal/models"
// )

// func TestHandler_getList(t *testing.T) {

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	sm := mocks.NewMockStorageManager(ctrl)

// 	r := chi.NewRouter()
// 	c := config.NewServerConfig()
// 	l := logger.NewLogger()

// 	app := &application{
// 		storageManager: sm,
// 		router:         r,
// 		logger:         l,
// 		config:         c,
// 	}
// 	app.setRouters()

// 	handler := http.HandlerFunc(app.getList)
// 	srv := httptest.NewServer(handler)
// 	defer srv.Close()

// 	list := []models.MetricsWithValue{
// 		{
// 			ID:    "testCounter",
// 			MType: "counter",
// 			Delta: int64(1),
// 		},
// 		{
// 			ID:    "testGauge",
// 			MType: "gauge",
// 			Value: float64(2),
// 		},
// 	}

// 	sm.EXPECT().
// 		GetList(gomock.Any()).
// 		Return(list, nil)

// 	testCases := []struct {
// 		name         string
// 		method       string
// 		body         string
// 		expectedCode int
// 		expectedBody string
// 	}{
// 		{
// 			name:         "method_get",
// 			method:       http.MethodGet,
// 			expectedCode: http.StatusOK,
// 			expectedBody: "testCounter: 1\ntestGauge: 2",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.method, func(t *testing.T) {
// 			req := resty.New().R()
// 			req.Method = tc.method
// 			req.URL = srv.URL

// 			if len(tc.body) > 0 {
// 				req.SetHeader("Content-Type", "application/json")
// 				req.SetBody(tc.body)
// 			}

// 			resp, err := req.Send()
// 			assert.NoError(t, err, "error making HTTP request")

// 			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
// 			if tc.expectedBody != "" {
// 				assert.Regexp(t, tc.expectedBody, string(resp.Body()))
// 			}
// 		})
// 	}
// }

// func TestHandler_getValue(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	sm := mocks.NewMockStorageManager(ctrl)

// 	r := chi.NewRouter()
// 	l := logger.NewLogger()

// 	app := &application{
// 		storageManager: sm,
// 		router:         r,
// 		logger:         l,
// 	}
// 	app.setRouters()

// 	handler := http.HandlerFunc(app.getValue)
// 	srv := httptest.NewServer(handler)
// 	defer srv.Close()

// 	counterMetric := models.MetricsWithValue{
// 		ID:    "testCounter",
// 		MType: "counter",
// 		Delta: int64(1),
// 	}

// 	sm.EXPECT().
// 		Get(gomock.Any(), gomock.Any()).
// 		Return(counterMetric, nil)

// 	testCases := []struct {
// 		name         string
// 		path         string
// 		method       string
// 		body         string
// 		expectedCode int
// 		expectedBody string
// 	}{
// 		{
// 			name:         "method_post",
// 			path:         "/value/",
// 			method:       http.MethodPost,
// 			body:         `{"id":"testCounter","type":"counter"}`,
// 			expectedCode: http.StatusOK,
// 			expectedBody: "1",
// 		},
// 		{
// 			name:         "wrong format",
// 			path:         "/value/",
// 			method:       http.MethodPost,
// 			body:         `id:"testCounter", type:"counter"`,
// 			expectedCode: http.StatusInternalServerError,
// 			expectedBody: "",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.method, func(t *testing.T) {
// 			req := resty.New().R()
// 			req.Method = tc.method
// 			req.URL = srv.URL
// 			req.URL += tc.path

// 			if len(tc.body) > 0 {
// 				req.SetHeader("Content-Type", "application/json")
// 				req.SetBody(tc.body)
// 			}

// 			resp, err := req.Send()
// 			assert.NoError(t, err, "error making HTTP request")

// 			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
// 			if tc.expectedBody != "" {
// 				assert.Regexp(t, tc.expectedBody, string(resp.Body()))
// 			}
// 		})
// 	}
// }

// func TestHandler_getCounter(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	sm := mocks.NewMockStorageManager(ctrl)

// 	r := chi.NewRouter()
// 	c := config.NewServerConfig()
// 	l := logger.NewLogger()

// 	app := &application{
// 		storageManager: sm,
// 		router:         r,
// 		logger:         l,
// 		config:         c,
// 	}
// 	app.setRouters()

// 	handler := http.HandlerFunc(app.getCounter)
// 	srv := httptest.NewServer(handler)
// 	defer srv.Close()

// 	counterMetric := models.MetricsWithValue{
// 		ID:    "testCounter",
// 		MType: "counter",
// 		Delta: int64(1),
// 	}

// 	sm.EXPECT().
// 		Get(gomock.Any(), gomock.Any()).
// 		Return(counterMetric, nil)

// 	testCases := []struct {
// 		name         string
// 		method       string
// 		body         string
// 		expectedCode int
// 		expectedBody string
// 	}{
// 		{
// 			name:         "method_get",
// 			method:       http.MethodGet,
// 			expectedCode: http.StatusOK,
// 			expectedBody: "1",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.method, func(t *testing.T) {
// 			req := resty.New().R()
// 			req.Method = tc.method
// 			req.URL = srv.URL

// 			if len(tc.body) > 0 {
// 				req.SetHeader("Content-Type", "application/json")
// 				req.SetBody(tc.body)
// 			}

// 			resp, err := req.Send()
// 			assert.NoError(t, err, "error making HTTP request")

// 			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
// 			if tc.expectedBody != "" {
// 				assert.Regexp(t, tc.expectedBody, string(resp.Body()))
// 			}
// 		})
// 	}
// }

// func TestHandler_getGauge(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	sm := mocks.NewMockStorageManager(ctrl)

// 	r := chi.NewRouter()
// 	l := logger.NewLogger()

// 	app := &application{
// 		storageManager: sm,
// 		router:         r,
// 		logger:         l,
// 	}
// 	app.setRouters()

// 	handler := http.HandlerFunc(app.getGauge)
// 	srv := httptest.NewServer(handler)
// 	defer srv.Close()

// 	gaugeMetric := models.MetricsWithValue{
// 		ID:    "testGauge",
// 		MType: "gauge",
// 		Value: float64(1),
// 	}

// 	sm.EXPECT().
// 		Get(gomock.Any(), gomock.Any()).
// 		Return(gaugeMetric, nil)

// 	testCases := []struct {
// 		name         string
// 		method       string
// 		body         string
// 		expectedCode int
// 		expectedBody string
// 	}{
// 		{
// 			name:         "method_get",
// 			method:       http.MethodGet,
// 			expectedCode: http.StatusOK,
// 			expectedBody: "1",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.method, func(t *testing.T) {
// 			req := resty.New().R()
// 			req.Method = tc.method
// 			req.URL = srv.URL

// 			if len(tc.body) > 0 {
// 				req.SetHeader("Content-Type", "application/json")
// 				req.SetBody(tc.body)
// 			}

// 			resp, err := req.Send()
// 			assert.NoError(t, err, "error making HTTP request")

// 			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
// 			if tc.expectedBody != "" {
// 				assert.Regexp(t, tc.expectedBody, string(resp.Body()))
// 			}
// 		})
// 	}
// }

// func TestHandler_updateList(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	sm := mocks.NewMockStorageManager(ctrl)

// 	r := chi.NewRouter()
// 	l := logger.NewLogger()

// 	app := &application{
// 		storageManager: sm,
// 		router:         r,
// 		logger:         l,
// 	}
// 	app.setRouters()

// 	handler := http.HandlerFunc(app.updateList)
// 	srv := httptest.NewServer(handler)
// 	defer srv.Close()

// 	sm.EXPECT().
// 		UpdateList(gomock.Any(), gomock.Any()).
// 		Return(nil)

// 	testCases := []struct {
// 		name         string
// 		path         string
// 		method       string
// 		body         string
// 		expectedCode int
// 		expectedBody string
// 	}{
// 		{
// 			name:         "method_post",
// 			path:         "/updates/",
// 			method:       http.MethodPost,
// 			body:         `[{"id":"testCounter","type":"counter","delta":1},{"id":"testGauge","type":"gauge","value":1}]`,
// 			expectedCode: http.StatusOK,
// 			expectedBody: "",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.method, func(t *testing.T) {
// 			req := resty.New().R()
// 			req.Method = tc.method
// 			req.URL = srv.URL
// 			req.URL += tc.path

// 			if len(tc.body) > 0 {
// 				req.SetHeader("Content-Type", "application/json")
// 				req.SetBody(tc.body)
// 			}

// 			resp, err := req.Send()
// 			assert.NoError(t, err, "error making HTTP request")

// 			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
// 			if tc.expectedBody != "" {
// 				assert.Regexp(t, tc.expectedBody, string(resp.Body()))
// 			}
// 		})
// 	}
// }

// func TestHandler_updateValue(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	sm := mocks.NewMockStorageManager(ctrl)

// 	r := chi.NewRouter()
// 	l := logger.NewLogger()

// 	app := &application{
// 		storageManager: sm,
// 		router:         r,
// 		logger:         l,
// 	}
// 	app.setRouters()

// 	handler := http.HandlerFunc(app.updateValue)
// 	srv := httptest.NewServer(handler)
// 	defer srv.Close()

// 	sm.EXPECT().
// 		Update(gomock.Any(), gomock.Any()).
// 		Return(nil)

// 	testCases := []struct {
// 		name         string
// 		method       string
// 		body         string
// 		expectedCode int
// 		expectedBody string
// 	}{
// 		{
// 			name:         "method_post",
// 			method:       http.MethodPost,
// 			body:         `{"id":"testCounter","type":"counter","delta":1}`,
// 			expectedCode: http.StatusOK,
// 			expectedBody: "",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.method, func(t *testing.T) {
// 			req := resty.New().R()
// 			req.Method = tc.method
// 			req.URL = srv.URL

// 			if len(tc.body) > 0 {
// 				req.SetHeader("Content-Type", "application/json")
// 				req.SetBody(tc.body)
// 			}

// 			resp, err := req.Send()
// 			assert.NoError(t, err, "error making HTTP request")

// 			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
// 			if tc.expectedBody != "" {
// 				assert.Regexp(t, tc.expectedBody, string(resp.Body()))
// 			}
// 		})
// 	}

// }
