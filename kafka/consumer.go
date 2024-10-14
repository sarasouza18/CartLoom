package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"golang.org/x/time/rate"
)

// ConsumeMessages reads messages from Kafka and processes them with rate limiting and retries.
func ConsumeMessages(ctx context.Context, reader *kafka.Reader, rdb *redis.Client, db *dynamodb.Client) error {
	limiter := rate.NewLimiter(5, 1) // Allow 5 messages per second

	for {
		if err := limiter.Wait(ctx); err != nil {
			log.Printf("Rate limiting error: %v", err)
			continue
		}

		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		log.Printf("Received message: Key=%s, Value=%s", string(msg.Key), string(msg.Value))

		if err := retryProcessOrder(ctx, rdb, db, string(msg.Key)); err != nil {
			log.Printf("Failed to process order %s after multiple attempts. Sending to DLQ.", string(msg.Key))
			SendToDLQ(ctx, msg)
		}
	}
}

// retryProcessOrder attempts to process the order with retries on failure
func retryProcessOrder(ctx context.Context, rdb *redis.Client, db *dynamodb.Client, orderID string) error {
	const maxRetries = 3
	const retryDelay = 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := processOrder(ctx, rdb, db, orderID)
		if err == nil {
			return nil // Order processed successfully
		}

		log.Printf("Error processing order %s (attempt %d/%d): %v", orderID, attempt, maxRetries, err)
		time.Sleep(retryDelay) // Delay before retrying
	}

	return logError("Max retries reached", orderID)
}

// processOrder processes the order using Redis for caching and DynamoDB for persistence
func processOrder(ctx context.Context, rdb *redis.Client, db *dynamodb.Client, orderID string) error {
	status, err := rdb.Get(ctx, orderID).Result()

	switch {
	case err == redis.Nil:
		log.Printf("Order %s not found in Redis. Processing...", orderID)
		if err := updateOrderStatus(ctx, rdb, db, orderID, "Processed"); err != nil {
			return err
		}
		log.Printf("Order %s processed and saved", orderID)

	case err != nil:
		log.Printf("Redis error for order %s: %v", orderID, err)
		return err

	default:
		log.Printf("Order %s already processed. Status: %s", orderID, status)
	}

	return nil
}

// updateOrderStatus updates the order status in Redis and saves it to DynamoDB
func updateOrderStatus(ctx context.Context, rdb *redis.Client, db *dynamodb.Client, orderID, status string) error {
	if err := rdb.Set(ctx, orderID, status, 0).Err(); err != nil {
		return logError("Failed to update Redis", orderID)
	}

	log.Printf("Order %s status updated to '%s' in Redis", orderID, status)

	item := map[string]types.AttributeValue{
		"OrderID": &types.AttributeValueMemberS{Value: orderID},
		"Status":  &types.AttributeValueMemberS{Value: status},
	}

	if _, err := db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("Orders"),
		Item:      item,
	}); err != nil {
		return logError("Failed to save order to DynamoDB", orderID)
	}

	return nil
}

// logError logs the error and returns it
func logError(message, orderID string) error {
	err := fmt.Errorf("%s for order %s", message, orderID)
	log.Println(err)
	return err
}
