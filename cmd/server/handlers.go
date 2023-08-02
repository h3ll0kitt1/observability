package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/h3ll0kitt1/observability/internal/models"
)

func (app *application) getAll(w http.ResponseWriter, r *http.Request) {
	list := ""
	metrics := app.storage.GetList()
	for _, metric := range metrics {
		if metric.MType == "counter" {
			list += fmt.Sprintf("%s: %d", metric.ID, metric.Delta)
			continue
		}
		list += fmt.Sprintf("%s: %f", metric.ID, metric.Value)
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(list))
}

func (app *application) getValue(w http.ResponseWriter, r *http.Request) {

	var metric models.Metrics

	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	app.logger.Infow("get value",
		"metric", metric,
	)

	metricWithValue := models.ToMetricWithValue(metric)
	m, ok := app.storage.GetValue(metricWithValue)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metric = models.ToMetric(m)
	jsonData, err := json.Marshal(metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
}

func (app *application) getCounter(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	metric := models.Metrics{
		ID:    name,
		MType: "counter",
	}

	metricWithValue := models.ToMetricWithValue(metric)
	m, ok := app.storage.GetValue(metricWithValue)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	valueStr := fmt.Sprintf("%d", m.Delta)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(valueStr))
}

func (app *application) getGauge(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	metric := models.Metrics{
		ID:    name,
		MType: "gauge",
	}

	metricWithValue := models.ToMetricWithValue(metric)
	m, ok := app.storage.GetValue(metricWithValue)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	valueStr := strconv.FormatFloat(m.Value, 'f', -1, 64)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(valueStr))
}

func (app *application) updateValue(w http.ResponseWriter, r *http.Request) {

	var metric models.Metrics

	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metricWithValue := models.ToMetricWithValue(metric)
	app.storage.Update(metricWithValue)

	app.logger.Infow("updated value",
		"metric", metricWithValue,
	)

	jsonData, err := json.Marshal(metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
}

func (app *application) updateCounter(w http.ResponseWriter, r *http.Request) {

	name := chi.URLParam(r, "name")
	valueStr := chi.URLParam(r, "value")

	value, ok := validateStringIsInt64(valueStr)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metric := models.MetricsWithValue{
		ID:    name,
		MType: "counter",
		Delta: value,
	}
	app.storage.Update(metric)
	w.WriteHeader(http.StatusOK)
}

func (app *application) updateGauge(w http.ResponseWriter, r *http.Request) {

	name := chi.URLParam(r, "name")
	valueStr := chi.URLParam(r, "value")

	value, ok := validateStringIsFloat64(valueStr)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metric := models.MetricsWithValue{
		ID:    name,
		MType: "gauge",
		Value: value,
	}
	app.storage.Update(metric)
	w.WriteHeader(http.StatusOK)
}

func errorUnknown(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func errorNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func errorNoName(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func validateStringIsInt64(s string) (int64, bool) {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1, false
	}
	return value, true
}

func validateStringIsFloat64(s string) (float64, bool) {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return -1, false
	}
	return value, true
}
