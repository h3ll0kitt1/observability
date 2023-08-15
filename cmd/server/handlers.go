package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/h3ll0kitt1/observability/internal/disk"
	"github.com/h3ll0kitt1/observability/internal/models"
)

func (app *application) getAll(w http.ResponseWriter, r *http.Request) {
	var list strings.Builder
	metrics := app.storage.GetList()
	for _, metric := range metrics {
		if metric.MType == "counter" {
			fmt.Fprintf(&list, "%s: %d\n", metric.ID, metric.Delta)
			continue
		}
		fmt.Fprintf(&list, "%s: %f\n", metric.ID, metric.Value)
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(list.String()))
}

func (app *application) ping(w http.ResponseWriter, r *http.Request) {

	if app.config.Database == "" {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if ok := app.storage.Ping(); !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
	metricWithValue, ok := app.storage.GetValue(metricWithValue)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metric = models.ToMetric(metricWithValue)
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

	metric := models.MetricsWithValue{
		ID:    name,
		MType: "counter",
	}

	metric, ok := app.storage.GetValue(metric)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	valueStr := fmt.Sprintf("%d", metric.Delta)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(valueStr))
}

func (app *application) getGauge(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	metric := models.MetricsWithValue{
		ID:    name,
		MType: "gauge",
	}

	metric, ok := app.storage.GetValue(metric)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	valueStr := strconv.FormatFloat(metric.Value, 'f', -1, 64)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(valueStr))
}

func (app *application) updateList(w http.ResponseWriter, r *http.Request) {
	var list []models.Metrics

	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	listWithValue := make([]models.MetricsWithValue, 0, len(list))
	for _, metric := range list {
		metricWithValue := models.ToMetricWithValue(metric)
		listWithValue = append(listWithValue, metricWithValue)
	}

	app.storage.UpdateList(listWithValue)

	if app.backupTime == 0 {
		disk.Flush(app.backupFile, app.storage)
	}

	w.WriteHeader(http.StatusOK)
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

	if app.backupTime == 0 {
		disk.Flush(app.backupFile, app.storage)
	}

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
	if app.backupTime == 0 {
		disk.Flush(app.backupFile, app.storage)
	}
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
	if app.backupTime == 0 {
		disk.Flush(app.backupFile, app.storage)
	}
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
