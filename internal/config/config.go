package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type ClientConfig struct {
	Protocol         string
	Addr             string
	Endpoint         string
	ReportInterval   time.Duration
	PollInterval     time.Duration
	RetryCount       int
	RetryWaitTime    time.Duration
	RetryMaxWaitTime time.Duration
}

type ServerConfig struct {
	Addr string
}

func NewClientConfig() *ClientConfig {

	var (
		flagReportInterval int
		flagPollInterval   int
		flagRunAddr        string
	)

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run client")
	flag.IntVar(&flagReportInterval, "r", 10, "number of seconds to report to server")
	flag.IntVar(&flagPollInterval, "p", 2, "number of seconds to update metrics")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	envReportInterval, err := strconv.Atoi(os.Getenv("REPORT_INTERVAL"))
	if err == nil {
		flagReportInterval = envReportInterval
	}

	envPollInterval, err := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	if err == nil {
		flagPollInterval = envPollInterval
	}

	protocol := "http://"
	addr := flagRunAddr
	endpoint := protocol + addr
	pollInterval := time.Duration(flagPollInterval) * time.Second
	reportInterval := time.Duration(flagReportInterval) * time.Second
	retryCount := 6
	retryWaitTime := 3 * time.Second
	retryMaxWaitTime := 90 * time.Second

	return &ClientConfig{
		Protocol:         protocol,
		Addr:             addr,
		Endpoint:         endpoint,
		ReportInterval:   reportInterval,
		PollInterval:     pollInterval,
		RetryCount:       retryCount,
		RetryWaitTime:    retryWaitTime,
		RetryMaxWaitTime: retryMaxWaitTime,
	}
}

func NewServerConfig() *ServerConfig {

	var flagRunAddr string

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	addr := flagRunAddr
	return &ServerConfig{
		Addr: addr,
	}
}
