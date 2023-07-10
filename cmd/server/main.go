package main

import (
	"log"

	"github.com/go-chi/chi/v5"

	"github.com/h3ll0kitt1/observability/internal/config"
	"github.com/h3ll0kitt1/observability/internal/router"
	"github.com/h3ll0kitt1/observability/internal/server"
	"github.com/h3ll0kitt1/observability/internal/storage/inmemory"
)

func main() {

	cfg := config.NewServerConfig()

	s := inmemory.NewStorage()
	r := chi.NewRouter()

	router.SetRouters(s, r)

	if err := server.Run(cfg, r); err != nil {
		log.Fatal(err)
		return
	}
}
