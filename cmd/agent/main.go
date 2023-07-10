package main

import (
	"flag"
	"time"

	"github.com/h3ll0kitt1/observability/internal/client"
)

var (
	flagReportInterval int
	flagPollInterval   int
	flagRunAddr        string
)

func main() {

	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run client")
	flag.IntVar(&flagReportInterval, "r", 10, "number of seconds to report to server")
	flag.IntVar(&flagPollInterval, "p", 2, "number of seconds to update metrics")
	flag.Parse()

	pollInterval := time.Duration(flagPollInterval) * time.Second
	reportInterval := time.Duration(flagReportInterval) * time.Second

	client.Run(flagRunAddr, pollInterval, reportInterval)
}
