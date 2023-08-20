package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/h3ll0kitt1/observability/internal/models"
)

func (app *application) getList(w http.ResponseWriter, r *http.Request) {
	var list strings.Builder
	metrics, err := app.storageManager.GetList(r.Context())
	if err != nil {
		app.logger.Errorw("error",
			"get list", err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
		return
	}

	if err := app.storageManager.Ping(); err != nil {
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
	metricWithValue, err = app.storageManager.Get(r.Context(), metricWithValue)
	if err != nil {
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

	metric, err := app.storageManager.Get(r.Context(), metric)
	if err != nil {
		app.logger.Errorw("error",
			"get counter", err,
		)
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

	metric, err := app.storageManager.Get(r.Context(), metric)
	if err != nil {

		app.logger.Errorw("error",
			"get gauge", err,
		)

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
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	listWithValue := make([]models.MetricsWithValue, 0, len(list))
	for _, metric := range list {
		metricWithValue := models.ToMetricWithValue(metric)
		listWithValue = append(listWithValue, metricWithValue)
	}

	if err := app.storageManager.UpdateList(r.Context(), listWithValue); err != nil {

		app.logger.Errorw("error",
			"update list", err,
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
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
	if err := app.storageManager.Update(r.Context(), metricWithValue); err != nil {

		app.logger.Errorw("error",
			"update value", err,
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
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

	if err := app.storageManager.Update(r.Context(), metric); err != nil {
		app.logger.Errorw("error",
			"update counter", err,
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
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

	if err := app.storageManager.Update(r.Context(), metric); err != nil {
		app.logger.Errorw("error",
			"update gauge", err,
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
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
