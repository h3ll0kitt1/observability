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
	Addr            string
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
	Database        string
}

func NewClientConfig() *ClientConfig {

	var (
		flagReportInterval int
		flagPollInterval   int
		flagRunAddr        string
		flagDatabase       string
	)

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run client")
	flag.StringVar(&flagDatabase, "d", "", "database to store metrics")
	flag.IntVar(&flagReportInterval, "r", 10, "number of seconds to report to server")
	flag.IntVar(&flagPollInterval, "p", 2, "number of seconds to update metrics")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	if envDatabase := os.Getenv("DATABASE_DSN"); envDatabase != "" {
		flagDatabase = envDatabase
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

	var (
		flagRunAddr         string
		flagFileStoragePath string
		flagDatabasePath    string
		flagStoreInterval   int
		flagRestore         bool
	)

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagFileStoragePath, "f", "/tmp/metrics-db.json", "full name of file to save metrics")
	flag.StringVar(&flagDatabasePath, "d", "", "sql database to store metrics")
	flag.IntVar(&flagStoreInterval, "i", 300, "interval in seconds to store metric values to file")
	flag.BoolVar(&flagRestore, "r", true, "bool value to show if previosly saved metrics should be loaded into server memory")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		flagFileStoragePath = envFileStoragePath
	}

	if envDatabasePath := os.Getenv("DATABASE_DSN"); envDatabasePath != "" {
		flagDatabasePath = envDatabasePath
	}

	envRestore, err := strconv.ParseBool(os.Getenv("RESTORE"))
	if err == nil {
		flagRestore = envRestore
	}

	envStoreInterval, err := strconv.Atoi(os.Getenv("STORE_INTERVAL"))
	if err == nil {
		flagStoreInterval = envStoreInterval
	}

	addr := flagRunAddr
	file := flagFileStoragePath
	storeInterval := time.Duration(flagStoreInterval) * time.Second
	restore := flagRestore
	database := flagDatabasePath

	return &ServerConfig{
		Addr:            addr,
		StoreInterval:   storeInterval,
		FileStoragePath: file,
		Restore:         restore,
		Database:        database,
	}
}
