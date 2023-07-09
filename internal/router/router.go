package router

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/h3ll0kitt1/observability/internal/storage"
)

func SetRouters(s storage.Storage, mux *chi.Mux) {

	mux.Route("/", func(mux chi.Router) {

		mux.Get("/", getAllHandle(s))

		mux.Route("/value", func(mux chi.Router) {
			mux.Get("/counter/{name}", getCounterHandle(s))
			mux.Get("/gauge/{name}", getGaugeHandle(s))
			mux.Get("/{other}/{name}", errorUnknownHandle)
		})

		mux.Route("/update", func(mux chi.Router) {

			mux.Route("/counter", func(mux chi.Router) {
				mux.Post("/", errorNoNameHandle)
				mux.Post("/{name}/{value}", updateCounterHandle(s))
			})

			mux.Route("/gauge", func(mux chi.Router) {
				mux.Post("/", errorNoNameHandle)
				mux.Post("/{name}/{value}", updateGaugeHandle(s))
			})
		})
	})

	mux.NotFound(errorNotFoundHandler)
}

func getAllHandle(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list := s.GetList()
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(list))
	}
}

func getCounterHandle(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		value, ok := s.GetValue("counter", name)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(value))
	}
}

func getGaugeHandle(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		value, ok := s.GetValue("gauge", name)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(value))
	}
}

func errorUnknownHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func errorNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func errorNoNameHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func updateCounterHandle(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		name := chi.URLParam(r, "name")
		valueStr := chi.URLParam(r, "value")

		value, ok := validateStringIsInt64(valueStr)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.Update(name, value)
		w.WriteHeader(http.StatusOK)
	}
}

func updateGaugeHandle(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		name := chi.URLParam(r, "name")
		valueStr := chi.URLParam(r, "value")

		value, ok := validateStringIsFloat64(valueStr)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.Update(name, value)
		w.WriteHeader(http.StatusOK)
	}
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
