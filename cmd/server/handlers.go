package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/h3ll0kitt1/observability/internal/models"
)

func (app *application) getAll(w http.ResponseWriter, r *http.Request) {
	list := app.storage.GetList()
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

	value, ok := app.storage.GetValue("counter", name)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

func (app *application) getGauge(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	value, ok := app.storage.GetValue("gauge", name)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
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
	app.storage.Update(name, value)
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
	app.storage.Update(name, value)
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
