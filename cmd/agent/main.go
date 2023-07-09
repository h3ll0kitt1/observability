package main

import (
	"time"

	"github.com/h3ll0kitt1/observability/internal/client"
)

func main() {

	reportInterval := 10 * time.Second
	pollInterval := 2 * time.Second
	endpoint := "http://localhost:8080"

	client.Run(endpoint, reportInterval, pollInterval)
}
