package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/h3ll0kitt1/observability/internal/storage/inmemory"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouterGet(t *testing.T) {
	s := inmemory.NewStorage()

	s.Update("testGauge", float64(2.0))
	s.Update("testCounter", int64(2))

	r := chi.NewRouter()
	SetRouters(s, r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	var testTable = []struct {
		url    string
		want   string
		status int
	}{
		// OK
		{"/", "testCounter : 2\ntestGauge : 2\n", http.StatusOK},
		{"/value/counter/testCounter", "2", http.StatusOK},
		{"/value/gauge/testGauge", "2", http.StatusOK},

		// WRONG
		{"/value/counter/unknownCounter", "", http.StatusNotFound},
		{"/value/gauge/unknownGauge", "", http.StatusNotFound},
	}

	for _, v := range testTable {
		resp, get := testRequest(t, ts, "GET", v.url)
		defer resp.Body.Close()
		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.want, get)
	}

}

func TestRouterPost(t *testing.T) {
	s := inmemory.NewStorage()
	r := chi.NewRouter()
	SetRouters(s, r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	var testTable = []struct {
		url    string
		status int
	}{
		// OK
		{"/update/gauge/testGauge/100", http.StatusOK},
		{"/update/counter/testCounter/100", http.StatusOK},

		// WRONG
		{"/update/counter/testCounter/100.0", http.StatusBadRequest},
		{"/update/counter/", http.StatusNotFound},
		{"/update/wrongtype/testCounter/100", http.StatusBadRequest},
	}

	for _, v := range testTable {
		resp, _ := testRequest(t, ts, "POST", v.url)
		defer resp.Body.Close()
		assert.Equal(t, v.status, resp.StatusCode)
	}

}