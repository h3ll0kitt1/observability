package main

import (
	"net/http"

	"github.com/h3ll0kitt1/observability/internal/handlers"
	"github.com/h3ll0kitt1/observability/internal/storage/inmemory"
)

func main() {

	storage := inmemory.NewStorage()

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, handlers.Update(storage))
	mux.HandleFunc(`/`, handlers.UpdateNotSpecified)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
