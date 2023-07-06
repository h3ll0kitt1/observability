package main

import (
	"github.com/h3ll0kitt1/observability/internal/client"
)

func main() {

	protocol := `http://`
	addr := `localhost:8080`
	reuestURL := protocol + addr
	client.Run(reuestURL)
}
