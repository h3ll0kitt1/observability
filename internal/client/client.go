package client

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/h3ll0kitt1/observability/internal/config"
)

type customClient struct {
	httpClient *resty.Client
	endpoint   string
}

type metrics struct {
	mapMetrics map[string]map[string]string
	pollCount  int64
}

func Run(cfg *config.ClientConfig) {

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	pollTicker := time.NewTicker(cfg.PollInterval)
	sendTicker := time.NewTicker(cfg.ReportInterval)

	client := newCustomClient(cfg)
	metrics := newMetrics()

	for {
		select {
		case <-pollTicker.C:
			metrics.update(rng)
		case <-sendTicker.C:
			metrics.sendToServer(client)
		}
	}
}

func newCustomClient(cfg *config.ClientConfig) customClient {
	httpClient := resty.New()

	httpClient.
		SetRetryCount(cfg.RetryCount).
		SetRetryWaitTime(cfg.RetryWaitTime).
		SetRetryMaxWaitTime(cfg.RetryMaxWaitTime)

	return customClient{httpClient: httpClient, endpoint: cfg.Endpoint}
}

func (m *metrics) sendToServer(client customClient) {
	for mtype, mmap := range m.mapMetrics {
		for mname, mvalue := range mmap {
			client.doRequestPOST(mtype, mname, mvalue)
		}
	}
}

func (c customClient) doRequestPOST(mtype, mname, mvalue string) {
	c.httpClient.R().SetPathParams(map[string]string{
		"type":  mtype,
		"name":  mname,
		"value": mvalue,
	}).
		Post(c.endpoint + "/update/{type}/{name}/{value}")
}

func newMetrics() metrics {
	mapMetrics := map[string]map[string]string{
		"counter": make(map[string]string),
		"gauge":   make(map[string]string),
	}
	return metrics{mapMetrics: mapMetrics, pollCount: 0}
}

func (m *metrics) update(rng *rand.Rand) {
	m.updateSpecificMemStats()
	m.updateRandomValue(rng)
	m.updateCounterValue()
}

func (m *metrics) updateSpecificMemStats() {

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	m.mapMetrics["gauge"]["Alloc"] = convertToString(ms.Alloc)
	m.mapMetrics["gauge"]["BuckHashSys"] = convertToString(ms.BuckHashSys)
	m.mapMetrics["gauge"]["Frees"] = convertToString(ms.Frees)
	m.mapMetrics["gauge"]["GCCPUFraction"] = convertToString(ms.GCCPUFraction)
	m.mapMetrics["gauge"]["GCSys"] = convertToString(ms.GCSys)
	m.mapMetrics["gauge"]["HeapAlloc"] = convertToString(ms.HeapAlloc)
	m.mapMetrics["gauge"]["HeapIdle"] = convertToString(ms.HeapIdle)
	m.mapMetrics["gauge"]["HeapInuse"] = convertToString(ms.HeapInuse)
	m.mapMetrics["gauge"]["HeapObjects"] = convertToString(ms.HeapObjects)
	m.mapMetrics["gauge"]["HeapReleased"] = convertToString(ms.HeapReleased)
	m.mapMetrics["gauge"]["HeapSys"] = convertToString(ms.HeapSys)
	m.mapMetrics["gauge"]["LastGC"] = convertToString(ms.LastGC)
	m.mapMetrics["gauge"]["Lookups"] = convertToString(ms.Lookups)
	m.mapMetrics["gauge"]["MCacheInuse"] = convertToString(ms.MCacheInuse)
	m.mapMetrics["gauge"]["MCacheSys"] = convertToString(ms.MCacheSys)
	m.mapMetrics["gauge"]["MSpanInuse"] = convertToString(ms.MSpanInuse)
	m.mapMetrics["gauge"]["MSpanSys"] = convertToString(ms.MSpanSys)
	m.mapMetrics["gauge"]["Mallocs"] = convertToString(ms.Mallocs)
	m.mapMetrics["gauge"]["NextGC"] = convertToString(ms.NextGC)
	m.mapMetrics["gauge"]["NumForcedGC"] = convertToString(ms.NumForcedGC)
	m.mapMetrics["gauge"]["NumGC"] = convertToString(ms.NumGC)
	m.mapMetrics["gauge"]["OtherSys"] = convertToString(ms.OtherSys)
	m.mapMetrics["gauge"]["PauseTotalNs"] = convertToString(ms.PauseTotalNs)
	m.mapMetrics["gauge"]["StackInuse"] = convertToString(ms.StackInuse)
	m.mapMetrics["gauge"]["StackSys"] = convertToString(ms.StackSys)
	m.mapMetrics["gauge"]["Sys"] = convertToString(ms.Sys)
	m.mapMetrics["gauge"]["TotalAlloc"] = convertToString(ms.TotalAlloc)
}

func convertToString(value any) string {
	res := ""
	switch v := value.(type) {
	case uint64:
		res = fmt.Sprintf("%f", float64(v))
	case float64:
		res = fmt.Sprintf("%f", v)
	}
	return res
}

func (m *metrics) updateRandomValue(rng *rand.Rand) {
	value := float64(rng.Intn(100))
	m.mapMetrics["gauge"]["Random"] = convertToString(value)
}

func (m *metrics) updateCounterValue() {
	m.pollCount++
	m.mapMetrics["counter"]["Counter"] = fmt.Sprintf("%d", m.pollCount)
}
