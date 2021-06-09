package testhelpers

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

func NewReader(ctx context.Context, topic string, offset *int64) (*kafka.Reader, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic,
	})
	if offset != nil {
		if err := reader.SetOffset(*offset); err != nil {
			return nil, err
		}
	} else {
		if err := reader.SetOffsetAt(ctx, time.Now()); err != nil {
			return nil, err
		}
	}
	return reader, nil
}

func NewWriter(topic string) (*kafka.Writer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Balancer: &kafka.LeastBytes{},
		Topic:    topic,
	}
	return writer, nil
}
