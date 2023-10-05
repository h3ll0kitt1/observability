package main

import (
	"github.com/h3ll0kitt1/observability/internal/client"
	"github.com/h3ll0kitt1/observability/internal/config"
)

func main() {
	cfg := config.NewClientConfig()
	cfg.Parse()
	client.Run(cfg)
}
