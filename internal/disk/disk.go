package disk

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/h3ll0kitt1/observability/internal/models"
	"github.com/h3ll0kitt1/observability/internal/storage"
)

func Load(filename string, storage storage.Storage) error {
	consumer, err := newConsumer(filename)
	if err != nil {
		return err
	}
	defer consumer.close()

	for {
		metric, err := consumer.readMetric()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		storage.Update(metric)
	}
	return nil
}

func Flush(filename string, storage storage.Storage) {
	producer, err := newProducer(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.close()

	metrics := storage.GetList()
	for _, metric := range metrics {
		log.Print("Flush metric to disk", metric)
		if err := producer.writeMetric(metric); err != nil {
			log.Fatal(err)
		}
	}
}

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
