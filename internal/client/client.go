package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
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
	pollCount  int64
}

type mapRW struct {
	metrics map[metricKey]models.Metrics
	mu      sync.RWMutex
}

func Run(cfg *config.ClientConfig) {

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	pollTicker := time.NewTicker(cfg.PollInterval)
	sendTicker := time.NewTicker(cfg.ReportInterval)
	defer pollTicker.Stop()
	defer sendTicker.Stop()

	client := newCustomClient(cfg)
	metrics := newMetrics()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pollTicker.C:
			go metrics.update(ctx, rng)
		case <-sendTicker.C:
			metrics.sendToServerWithRate(ctx, client, cfg.RateLimit)
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

func (m *metrics) sendToServerWithRate(ctx context.Context, client customClient, limit int) {

	ch := make(chan models.Metrics, 256)

	for i := 0; i < limit; i++ {
		go client.sentToServerWorker(ctx, ch)
	}

	for _, metric := range m.mapMetrics.metrics {
		ch <- metric
	}
	close(ch)
}

func (c customClient) sentToServerWorker(ctx context.Context, ch <-chan models.Metrics) {

	for metric := range ch {
		err := c.doRequestPOST(ctx, metric)
		if err != nil {
			log.Printf("%s\n", err)
		}
	}
}

func (c customClient) doRequestPOST(ctx context.Context, metric models.Metrics) error {

	jsonData, err := json.Marshal(metric)
	if err != nil {
		return errors.New("error converting metric to json")
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
		SetContext(ctx).
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
	return metrics{
		mapMetrics: mapMetrics,
		pollCount:  0}
}

func newMapRW() *mapRW {
	var m mapRW
	m.metrics = make(map[metricKey]models.Metrics)
	return &m
}

func (m *metrics) update(ctx context.Context, rng *rand.Rand) {
	m.updateSpecificMemStats(ctx)
	m.updateRandomValue(ctx, rng)
	m.updateCounterValue(ctx)
	go m.updateMemoryCPUInfo(ctx)
}

func (m *metrics) updateSpecificMemStats(ctx context.Context) {

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

		m.mapMetrics.mu.RLock()
		m.mapMetrics.metrics[key] = metric
		m.mapMetrics.mu.RUnlock()
	}
}

func (m *metrics) updateMemoryCPUInfo(ctx context.Context) {

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

		m.mapMetrics.mu.RLock()
		m.mapMetrics.metrics[key] = metric
		m.mapMetrics.mu.RUnlock()
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

func (m *metrics) updateRandomValue(ctx context.Context, rng *rand.Rand) {
	id, mtype := "RandomValue", "gauge"
	value := float64(rng.Intn(100))

	key := metricKey{id: id, mtype: mtype}
	metric := models.Metrics{
		ID:    id,
		MType: mtype,
		Value: &value,
	}

	m.mapMetrics.mu.RLock()
	m.mapMetrics.metrics[key] = metric
	m.mapMetrics.mu.RUnlock()
}

func (m *metrics) updateCounterValue(ctx context.Context) {
	atomic.AddInt64(&m.pollCount, 1)

	id, mtype := "PollCount", "counter"
	value := atomic.LoadInt64(&m.pollCount)

	key := metricKey{id: id, mtype: mtype}
	metric := models.Metrics{
		ID:    id,
		MType: mtype,
		Delta: &value,
	}

	m.mapMetrics.mu.RLock()
	m.mapMetrics.metrics[key] = metric
	m.mapMetrics.mu.RUnlock()
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
