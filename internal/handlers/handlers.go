package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/h3ll0kitt1/observability/internal/storage/inmemory"
)

func Update(s *inmemory.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPart := strings.Split(r.URL.Path, "/")
		w.Header().Set("Content-Type", "text/plain")
		if len(urlPart) < 3 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if urlPart[2] != "counter" && urlPart[2] != "gauge" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metricType := urlPart[2]

		if len(urlPart) < 4 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metricName := urlPart[3]

		if len(urlPart) < 5 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metricValueStr := urlPart[4]

		if metricType == "counter" {
			metricValue, ok := validateStringIsInt64(metricValueStr)
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			s.Update(metricName, metricValue)
		} else if metricType == "gauge" {
			metricValue, ok := validateStringIsFloat64(metricValueStr)
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			s.Update(metricName, metricValue)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateNotSpecified(w http.ResponseWriter, r *http.Request) {
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
