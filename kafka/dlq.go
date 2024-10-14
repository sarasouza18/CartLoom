package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

// DLQWriter handles sending messages to the Dead Letter Queue (DLQ)
type DLQWriter struct {
	writer *kafka.Writer
}

// NewDLQWriter initializes a new Kafka writer for the DLQ
func NewDLQWriter(brokers []string, topic string) *DLQWriter {
	return &DLQWriter{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
	}
}

// SendToDLQ sends the failed message to the Dead Letter Queue (DLQ)
func SendToDLQ(ctx context.Context, msg kafka.Message) {
	dlq := NewDLQWriter([]string{"kafka:9092"}, "dlq-orders")
	defer dlq.Close()

	if err := dlq.SendMessage(ctx, msg); err != nil {
		log.Printf("Failed to send message to DLQ: %v", err)
	}
}

// SendMessage sends a message to the DLQ
func (dlq *DLQWriter) SendMessage(ctx context.Context, msg kafka.Message) error {
	if err := dlq.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("Failed to write message to DLQ: %v", err)
		return err
	}

	log.Printf("Message sent to DLQ: Key=%s", string(msg.Key))
	return nil
}

// Close closes the Kafka writer
func (dlq *DLQWriter) Close() error {
	if err := dlq.writer.Close(); err != nil {
		log.Printf("Failed to close DLQ writer: %v", err)
		return err
	}
	return nil
}
