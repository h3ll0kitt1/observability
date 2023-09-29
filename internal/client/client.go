package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/shirou/gopsutil/mem"

	"github.com/h3ll0kitt1/observability/internal/config"
	"github.com/h3ll0kitt1/observability/internal/hash"
	"github.com/h3ll0kitt1/observability/internal/models"
)

var (
	ErrServerUnavailable = errors.New("error doing post request")
)

type customClient struct {
	httpClient *resty.Client
	endpoint   string
	key        string
}

type metricKey struct {
	id    string
	mtype string
}

type metrics struct {
	mapMetrics *mapRW
	arrMetrics []models.Metrics
	pollCount  int64
}

type mapRW struct {
	metrics map[metricKey]models.Metrics
	sync.RWMutex
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
			go metrics.update(rng)
		case <-sendTicker.C:
			metrics.sendToServerWithRate(client, cfg.RateLimit)
		}
	}
}

func newCustomClient(cfg *config.ClientConfig) customClient {
	httpClient := resty.New()

	httpClient.
		SetRetryCount(cfg.RetryCount).
		SetRetryWaitTime(cfg.RetryWaitTime).
		SetRetryMaxWaitTime(cfg.RetryMaxWaitTime)

	return customClient{
		httpClient: httpClient,
		endpoint:   cfg.Endpoint,
		key:        cfg.Key,
	}
}

func (m *metrics) sendToServerWithRate(client customClient, limit int) {

	ch := make(chan models.Metrics, 256)
	for _, metric := range m.arrMetrics {
		ch <- metric
	}
	close(ch)

	for i := 0; i < limit; i++ {
		go client.worker(i, ch)
	}
}

func (c customClient) worker(i int, ch chan models.Metrics) {

	for metric := range ch {
		err := c.doRequestPOST(metric)
		if err != nil {
			log.Printf("%s\n", err)
		}
	}
}

func (c customClient) doRequestPOST(metric models.Metrics) error {

	jsonData, err := json.Marshal(metric)
	if err != nil {
		return errors.New("error converting slice of metrics to json")
	}

	if c.key != "" {
		hash := hash.ComputeSHA256(jsonData, c.key)
		c.httpClient.R().SetHeader("HashSHA256", hash)
	}

	gzipData, err := GzipCompress(jsonData)
	if err != nil {
		return errors.New("error compressing json to gzip")
	}

	_, err = c.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(gzipData).
		Post(c.endpoint + "/update/")

	if err != nil {
		return ErrServerUnavailable
	}
	return nil
}

func newMetrics() metrics {
	mapMetrics := newMapRW()
	arrMetrics := make([]models.Metrics, 0)
	return metrics{
		mapMetrics: mapMetrics,
		arrMetrics: arrMetrics,
		pollCount:  0}
}

func newMapRW() *mapRW {
	var m mapRW
	m.metrics = make(map[metricKey]models.Metrics)
	return &m
}

func (m *metrics) update(rng *rand.Rand) {
	m.arrMetrics = m.arrMetrics[:0]
	m.updateSpecificMemStats()
	m.updateRandomValue(rng)
	m.updateCounterValue()
	go m.updateMemoryCPUInfo()

}

func (m *metrics) updateSpecificMemStats() {

	searchedFields := map[string]bool{
		"Alloc":         true,
		"BuckHashSys":   true,
		"Frees":         true,
		"GCCPUFraction": true,
		"GCSys":         true,
		"HeapAlloc":     true,
		"HeapIdle":      true,
		"HeapInuse":     true,
		"HeapObjects":   true,
		"HeapReleased":  true,
		"HeapSys":       true,
		"LastGC":        true,
		"Lookups":       true,
		"MCacheInuse":   true,
		"MCacheSys":     true,
		"MSpanInuse":    true,
		"MSpanSys":      true,
		"Mallocs":       true,
		"NextGC":        true,
		"NumForcedGC":   true,
		"NumGC":         true,
		"OtherSys":      true,
		"PauseTotalNs":  true,
		"StackInuse":    true,
		"StackSys":      true,
		"Sys":           true,
		"TotalAlloc":    true,
	}

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	v := reflect.ValueOf(&ms).Elem()

	for i := 0; i < v.NumField(); i++ {

		id := v.Type().Field(i).Name
		if _, ok := searchedFields[id]; !ok {
			continue
		}

		mtype := "gauge"
		value := getFloat64(v.Field(i).Interface())

		key := metricKey{id: id, mtype: mtype}
		metric := models.Metrics{
			ID:    id,
			MType: mtype,
			Value: &value,
		}

		m.mapMetrics.RLock()
		m.mapMetrics.metrics[key] = metric
		m.mapMetrics.RUnlock()
		m.arrMetrics = append(m.arrMetrics, metric)
	}
}

func (m *metrics) updateMemoryCPUInfo() {

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("%s\n", err)
	}

	searchedFields := map[string]bool{
		"Total":       true,
		"Free":        true,
		"UsedPercent": true,
	}

	v := reflect.ValueOf(vmStat).Elem()

	for i := 0; i < v.NumField(); i++ {
		id := v.Type().Field(i).Name
		if _, ok := searchedFields[id]; !ok {
			continue
		}

		mtype := "gauge"
		value := getFloat64(v.Field(i).Interface())

		key := metricKey{id: id, mtype: mtype}
		metric := models.Metrics{
			ID:    id,
			MType: mtype,
			Value: &value,
		}

		m.mapMetrics.RLock()
		m.mapMetrics.metrics[key] = metric
		m.mapMetrics.RUnlock()
		m.arrMetrics = append(m.arrMetrics, metric)
	}
}

func getFloat64(value any) float64 {
	switch i := value.(type) {
	case float64:
		return float64(i)
	case float32:
		return float64(i)
	case int64:
		return float64(i)
	case int32:
		return float64(i)
	case uint64:
		return float64(i)
	case uint32:
		return float64(i)
	default:
		return -1
	}
}

func (m *metrics) updateRandomValue(rng *rand.Rand) {
	id, mtype := "RandomValue", "gauge"
	value := float64(rng.Intn(100))

	key := metricKey{id: id, mtype: mtype}
	metric := models.Metrics{
		ID:    id,
		MType: mtype,
		Value: &value,
	}

	m.mapMetrics.RLock()
	m.mapMetrics.metrics[key] = metric
	m.mapMetrics.RUnlock()
	m.arrMetrics = append(m.arrMetrics, metric)
}

func (m *metrics) updateCounterValue() {
	m.pollCount++

	id, mtype := "PollCount", "counter"
	value := m.pollCount

	key := metricKey{id: id, mtype: mtype}
	metric := models.Metrics{
		ID:    id,
		MType: mtype,
		Delta: &value,
	}

	m.mapMetrics.RLock()
	m.mapMetrics.metrics[key] = metric
	m.mapMetrics.RUnlock()
	m.arrMetrics = append(m.arrMetrics, metric)
}

func GzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	w := gzip.NewWriter(&buf)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
