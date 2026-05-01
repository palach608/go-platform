package broker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type HandlerFunc func(ctx context.Context, value []byte) error

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (c *Consumer) Consume(ctx context.Context, handler HandlerFunc) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("kafka consumer error: %v", err)
			continue
		}
		if err := handler(ctx, msg.Value); err != nil {
			log.Printf("kafka handler error: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func Unmarshal[T any](data []byte) (T, error) {
	var v T
	err := json.Unmarshal(data, &v)
	return v, err
}
