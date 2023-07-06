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

func Run(addr string) {
	rand.Seed(time.Now().UnixNano())

	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second
	metrics := newMetrics()

	metrics.pollCount = 0

	for {
		metrics.updateSpecificMemStats()
		metrics.updateRandomValue()
		metrics.pollCount++

		metrics.sendToServer(addr, reportInterval-pollInterval)

		time.Sleep(pollInterval)
	}
}

type metrics struct {
	mapGauge  map[string]float64
	pollCount int64
}

func newMetrics() metrics {
	mapGauge := make(map[string]float64)
	return metrics{mapGauge, 0}
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

func (m *metrics) updateRandomValue() {
	m.mapGauge["Random"] = rand.Float64() * 8
}

func (m *metrics) sendToServer(addr string, reportInterval time.Duration) {
	for name, value := range m.mapGauge {
		go doRequestPOST(addr, name, value)
	}
	go doRequestPOST(addr, "MyCounter", m.pollCount)
	time.Sleep(reportInterval)
}

func doRequestPOST(addr string, metricName string, metricValue any) {

	requestURL := constructURL(addr, metricName, metricValue)
	req, err := http.NewRequest(http.MethodPost, requestURL, http.NoBody)
	req.Header.Set("Content-Type", "text/plain")

	httpClient := &http.Client{}

	_, err = httpClient.Do(req)
	if err != nil {
		log.Fatalf("Error making POST request wit path %s\n", requestURL)
	}
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
