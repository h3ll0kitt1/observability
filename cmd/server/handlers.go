package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/h3ll0kitt1/observability/internal/models"
)

func (app *application) getValue(w http.ResponseWriter, r *http.Request) {

	var metric models.Metrics

	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	value, ok := app.storage.GetValue(metric.MType, metric.ID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if metric.MType == "counter" {
		v, ok := validateStringIsInt64(value)
		if ok {
			metric.Delta = &v
		}
	}

	if metric.MType == "gauge" {
		v, ok := validateStringIsFloat64(value)
		if ok {
			metric.Value = &v
		}
	}

	metricJSON, err := metric.Convert2JSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metricJSON))
}

func (app *application) updateValue(w http.ResponseWriter, r *http.Request) {

	var metric models.Metrics

	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if metric.Delta != nil {
		app.storage.Update(metric.ID, *(metric.Delta))
	}

	if metric.Value != nil {
		app.storage.Update(metric.ID, *(metric.Value))
	}

	metricJSON, err := metric.Convert2JSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metricJSON))
}

func errorNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
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
