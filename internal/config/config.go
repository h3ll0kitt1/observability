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
	Key              string
	ReportInterval   time.Duration
	PollInterval     time.Duration
	RateLimit        int
	RetryCount       int
	RetryWaitTime    time.Duration
	RetryMaxWaitTime time.Duration
}

type ServerConfig struct {
	Addr            string
	Key             string
	Database        string
	FileStoragePath string
	Restore         bool
	StoreInterval   time.Duration
}

func NewClientConfig() *ClientConfig {
	var cc ClientConfig
	return &cc
}

func (cc *ClientConfig) Parse() {
	var (
		flagReportInterval int
		flagPollInterval   int
		flagRunAddr        string
		flagDatabase       string
		flagKey            string
		flagRateLimit      int
	)

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run client")
	flag.StringVar(&flagDatabase, "d", "", "database to store metrics")
	flag.StringVar(&flagKey, "k", "", "symmetrical key for SHA256 hash function")
	flag.IntVar(&flagReportInterval, "r", 10, "number of seconds to report to server")
	flag.IntVar(&flagPollInterval, "p", 2, "number of seconds to update metrics")
	flag.IntVar(&flagRateLimit, "l", 2, "number of concurrent post requests to server")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	if envDatabase := os.Getenv("DATABASE_DSN"); envDatabase != "" {
		flagDatabase = envDatabase
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		flagKey = envKey
	}

	envReportInterval, err := strconv.Atoi(os.Getenv("REPORT_INTERVAL"))
	if err == nil {
		flagReportInterval = envReportInterval
	}

	envPollInterval, err := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	if err == nil {
		flagPollInterval = envPollInterval
	}

	envRateLimit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if err == nil {
		flagRateLimit = envRateLimit
	}

	protocol := "http://"
	addr := flagRunAddr
	endpoint := protocol + addr
	key := flagKey
	pollInterval := time.Duration(flagPollInterval) * time.Second
	reportInterval := time.Duration(flagReportInterval) * time.Second
	rateLimit := flagRateLimit
	retryCount := 6
	retryWaitTime := 3 * time.Second
	retryMaxWaitTime := 90 * time.Second

	cc.Protocol = protocol
	cc.Addr = addr
	cc.Endpoint = endpoint
	cc.Key = key
	cc.ReportInterval = reportInterval
	cc.PollInterval = pollInterval
	cc.RateLimit = rateLimit
	cc.RetryCount = retryCount
	cc.RetryWaitTime = retryWaitTime
	cc.RetryMaxWaitTime = retryMaxWaitTime
}

func NewServerConfig() *ServerConfig {
	var sc ServerConfig
	return &sc
}

func (sc *ServerConfig) Parse() {
	var (
		flagRunAddr         string
		flagFileStoragePath string
		flagDatabasePath    string
		flagKey             string
		flagStoreInterval   int
		flagRestore         bool
	)

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagFileStoragePath, "f", "/tmp/metrics-db.json", "full name of file to save metrics")
	flag.StringVar(&flagDatabasePath, "d", "", "sql database to store metrics")
	flag.StringVar(&flagKey, "k", "", "symmetrical key for SHA256 hash function")
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

	if envKey := os.Getenv("KEY"); envKey != "" {
		flagKey = envKey
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
	key := flagKey

	sc.Addr = addr
	sc.StoreInterval = storeInterval
	sc.FileStoragePath = file
	sc.Restore = restore
	sc.Database = database
	sc.Key = key
}
