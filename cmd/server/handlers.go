package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/h3ll0kitt1/observability/internal/models"
)

func (app *application) getList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var list strings.Builder

	// здесь надо добавить retrieble
	metrics, err := app.storage.GetList(ctx)
	if err != nil {
		app.logger.Infow("error",
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
	}

	if err := app.storage.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) getValue(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
	// здесь надо добавить retrieble
	metricWithValue, err = app.storage.Get(ctx, metricWithValue)
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	name := chi.URLParam(r, "name")

	metric := models.MetricsWithValue{
		ID:    name,
		MType: "counter",
	}

	// здесь надо добавить retrieble
	metric, err := app.storage.Get(ctx, metric)
	if err != nil {
		app.logger.Infow("error",
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	name := chi.URLParam(r, "name")

	metric := models.MetricsWithValue{
		ID:    name,
		MType: "gauge",
	}

	// здесь надо добавить retrieble
	metric, err := app.storage.Get(ctx, metric)
	if err != nil {

		app.logger.Infow("error",
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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

	if err := app.storage.UpdateList(ctx, listWithValue); err != nil {

		app.logger.Infow("error",
			"update list", err,
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if app.config.StoreInterval == 0 {
		app.flush()
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) updateValue(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var metric models.Metrics
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metricWithValue := models.ToMetricWithValue(metric)
	if err := app.storage.Update(ctx, metricWithValue); err != nil {

		app.logger.Infow("error",
			"update value", err,
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if app.config.StoreInterval == 0 {
		app.flush()
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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

	if err := app.storage.Update(ctx, metric); err != nil {
		app.logger.Infow("error",
			"update counter", err,
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if app.config.StoreInterval == 0 {
		app.flush()
	}
	w.WriteHeader(http.StatusOK)
}

func (app *application) updateGauge(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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

	if err := app.storage.Update(ctx, metric); err != nil {
		app.logger.Infow("error",
			"update gauge", err,
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if app.config.StoreInterval == 0 {
		app.flush()
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
