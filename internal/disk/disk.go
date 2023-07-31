package disk

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/h3ll0kitt1/observability/internal/models"
	"github.com/h3ll0kitt1/observability/internal/storage"
)

type metric struct {
	ID    string `json:"id"`
	Value any    `json:"value"`
}

func Load(filename string, storage storage.Storage) error {
	consumer, err := newConsumer(filename)
	if err != nil {
		return err
	}
	defer consumer.close()

	for {
		event, err := consumer.readEvent()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		storage.Update(event.ID, event.Value)
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
		if err := producer.writeEvent(metric); err != nil {
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

func (c *consumer) readEvent() (*metric, error) {
	event := &metric{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
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

func (p *producer) writeEvent(event *models.Metrics) error {
	return p.encoder.Encode(&event)
}

func (p *producer) close() error {
	return p.file.Close()
}
