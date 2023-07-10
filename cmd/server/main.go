package main

import (
	"flag"
	"log"

	"github.com/go-chi/chi/v5"

	"github.com/h3ll0kitt1/observability/internal/router"
	"github.com/h3ll0kitt1/observability/internal/server"
	"github.com/h3ll0kitt1/observability/internal/storage/inmemory"
)

var flagRunAddr string

func main() {

	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()

	s := inmemory.NewStorage()
	r := chi.NewRouter()

	router.SetRouters(s, r)

	if err := server.Run(flagRunAddr, r); err != nil {
		log.Fatal(err)
		return
	}
}
