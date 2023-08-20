package file

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/h3ll0kitt1/observability/internal/models"
)

type FileStorage struct {
	filename string
}

func NewStorage(filename string) *FileStorage {
	return &FileStorage{
		filename: filename,
	}
}

func (fs *FileStorage) GetList(ctx context.Context) ([]models.MetricsWithValue, error) {
	list := make([]models.MetricsWithValue, 0)

	consumer, err := newConsumer(fs.filename)
	if err != nil {
		return nil, err
	}
	defer consumer.close()

	for {
		metric, err := consumer.readMetric()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		metricWithValue := models.ToMetricWithValue(*metric)
		list = append(list, metricWithValue)
	}
	return list, nil
}

func (fs *FileStorage) UpdateList(ctx context.Context, list []models.MetricsWithValue) error {
	producer, err := newProducer(fs.filename)
	if err != nil {
		return err
	}
	defer producer.close()

	for _, metric := range list {
		m := models.ToMetric(metric)
		if err := producer.writeMetric(&m); err != nil {
			return err
		}
	}
	return nil
}

func (fs *FileStorage) Ping() error { return nil }

func (fs *FileStorage) SetRetryCount(attempts int) {}

func (fs *FileStorage) SetRetryStartWaitTime(sleep time.Duration) {}

func (fs *FileStorage) SetRetryIncreseWaitTime(delta time.Duration) {}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func newConsumer(filename string) (*consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *consumer) readMetric() (*models.Metrics, error) {
	var metric models.Metrics
	if err := c.decoder.Decode(&metric); err != nil {
		return nil, err
	}
	return &metric, nil
}

func (c *consumer) close() error {
	return c.file.Close()
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func newProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) writeMetric(metric *models.Metrics) error {
	return p.encoder.Encode(&metric)
}

func (p *producer) close() error {
	return p.file.Close()
}
