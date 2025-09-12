package broker

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(brokers []string, groupID, topic string) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           brokers,
		GroupID:           groupID,
		Topic:             topic,
		MinBytes:          10e3, // 10KB
		MaxBytes:          10e6, // 10MB
		CommitInterval:    time.Second,
		HeartbeatInterval: time.Second * 3,
	})
	return &KafkaConsumer{reader: r}
}

func (c *KafkaConsumer) Consume(ctx context.Context, handler func(kafka.Message) error) error {
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return err
		}

		if err := handler(m); err != nil {
			log.Printf("❌ Handler failed: %v", err)
			continue
		}

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Printf("❌ Commit failed: %v", err)
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
