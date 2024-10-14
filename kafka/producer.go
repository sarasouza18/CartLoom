package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// ProduceMessages generates and sends a series of order messages to Kafka
func ProduceMessages(ctx context.Context, writer *kafka.Writer) error {
	for i := 0; i < 5; i++ {
		if err := produceOrder(ctx, writer, i); err != nil {
			return err
		}
		time.Sleep(1 * time.Second) // Simulate delay between orders
	}
	return nil
}

// produceOrder generates and sends a single order message to Kafka
func produceOrder(ctx context.Context, writer *kafka.Writer, orderID int) error {
	message := kafka.Message{
		Key:   []byte(fmt.Sprintf("OrderID-%d", orderID)),
		Value: []byte(fmt.Sprintf("Order Created: #%d", orderID)),
	}

	if err := writer.WriteMessages(ctx, message); err != nil {
		log.Printf("Failed to send order %d to Kafka: %v", orderID, err)
		return err
	}

	log.Printf("Order %d sent to Kafka", orderID)
	return nil
}
