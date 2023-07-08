package client

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type customClient struct {
	httpClient *http.Client
	addr       string
}

type metrics struct {
	mapGauge  map[string]float64
	pollCount int64
}

func Run(addr string) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second
	metrics := newMetrics()

	metrics.pollCount = 0

	client := newCustomClient(addr)

	for {
		metrics.updateSpecificMemStats()
		metrics.updateRandomValue(rng)
		metrics.pollCount++

		metrics.sendToServer(client, reportInterval-pollInterval)

		time.Sleep(pollInterval)
	}
}

func newMetrics() metrics {
	mapGauge := make(map[string]float64)
	return metrics{mapGauge, 0}
}

func newCustomClient(addr string) customClient {
	timeout := 5 * time.Second
	httpClient := &http.Client{
		Timeout: timeout,
	}

	req, _ := http.NewRequest(http.MethodPost, addr, http.NoBody)
	req.Header.Set("Content-Type", "text/plain")

	client := http.Client{}

	for {
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			break
		}
		log.Printf("Error making POST request with path %s\n", addr)
		time.Sleep(5 * time.Second)
	}
	return customClient{httpClient, addr}
}

func (m *metrics) updateSpecificMemStats() {

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	m.mapGauge["Alloc"] = float64(ms.Alloc)
	m.mapGauge["BuckHashSys"] = float64(ms.BuckHashSys)
	m.mapGauge["Frees"] = float64(ms.Frees)
	m.mapGauge["GCCPUFraction"] = ms.GCCPUFraction
	m.mapGauge["GCSys"] = float64(ms.GCSys)
	m.mapGauge["HeapAlloc"] = float64(ms.HeapAlloc)
	m.mapGauge["HeapIdle"] = float64(ms.HeapIdle)
	m.mapGauge["HeapInuse"] = float64(ms.HeapInuse)
	m.mapGauge["HeapObjects"] = float64(ms.HeapObjects)
	m.mapGauge["HeapReleased"] = float64(ms.HeapReleased)
	m.mapGauge["HeapSys"] = float64(ms.HeapSys)
	m.mapGauge["LastGC"] = float64(ms.LastGC)
	m.mapGauge["Lookups"] = float64(ms.Lookups)
	m.mapGauge["MCacheInuse"] = float64(ms.MCacheInuse)
	m.mapGauge["MCacheSys"] = float64(ms.MCacheSys)
	m.mapGauge["MSpanInuse"] = float64(ms.MSpanInuse)
	m.mapGauge["MSpanSys"] = float64(ms.MSpanSys)
	m.mapGauge["Mallocs"] = float64(ms.Mallocs)
	m.mapGauge["NextGC"] = float64(ms.NextGC)
	m.mapGauge["NumForcedGC"] = float64(ms.NumForcedGC)
	m.mapGauge["NumGC"] = float64(ms.NumGC)
	m.mapGauge["OtherSys"] = float64(ms.OtherSys)
	m.mapGauge["PauseTotalNs"] = float64(ms.PauseTotalNs)
	m.mapGauge["StackInuse"] = float64(ms.StackInuse)
	m.mapGauge["StackSys"] = float64(ms.StackSys)
	m.mapGauge["Sys"] = float64(ms.Sys)
	m.mapGauge["TotalAlloc"] = float64(ms.TotalAlloc)
}

func (m *metrics) updateRandomValue(rng *rand.Rand) {
	value := rng.Intn(100)
	m.mapGauge["Random"] = float64(value)
}

func (m *metrics) sendToServer(client customClient, reportInterval time.Duration) {
	for name, value := range m.mapGauge {
		doRequestPOST(client, name, value)
	}
	doRequestPOST(client, "MyCounter", m.pollCount)
	time.Sleep(reportInterval)
}

func doRequestPOST(client customClient, metricName string, metricValue any) {

	requestURL := constructURL(client.addr, metricName, metricValue)
	req, err := http.NewRequest(http.MethodPost, requestURL, http.NoBody)
	if err != nil {
		log.Fatalf("Error constructing POST request with path %s\n", requestURL)
	}
	req.Header.Set("Content-Type", "text/plain")

	var resp *http.Response
	resp, err = client.httpClient.Do(req)
	if err != nil {
		log.Printf("Error making POST request with path %s\n", requestURL)
	}
	defer resp.Body.Close()
	log.Printf("Request SEND: %s\n", requestURL)
}

func constructURL(addr string, metricName string, metricValue any) string {
	requestURL := "/update"
	switch mv := metricValue.(type) {
	case int64:
		metricValueStr := strconv.FormatInt(mv, 10)
		requestURL += joinPath("counter", metricName, metricValueStr)
	case float64:
		metricValueStr := fmt.Sprintf("%f", mv)
		requestURL += joinPath("gauge", metricName, metricValueStr)
	}
	return addr + requestURL
}

func joinPath(metricType, metricName, metricValueStr string) string {
	pathParts := []string{"", metricType, metricName, metricValueStr}
	path := strings.Join(pathParts, "/")
	return path
}
