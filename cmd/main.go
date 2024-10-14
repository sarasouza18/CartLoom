package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	kafka_go "github.com/segmentio/kafka-go"

	"cartloom/dynamodb"
	"cartloom/kafka"
	"cartloom/redis"
	"cartloom/shopify"
	"cartloom/utils"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create context for the application
	ctx := context.Background()

	// Initialize Redis and DynamoDB with environment variables
	initializeRedisAndDynamoDB(ctx)

	// Register Shopify webhook
	registerShopifyWebhook()

	// Start Prometheus metrics and Kafka services
	startMetricsServer()
	startKafka(ctx)

	// Set up logging
	setupLogging()

	// Keep the program running
	select {}
}

// initializeRedisAndDynamoDB initializes Redis and DynamoDB with environment variables
func initializeRedisAndDynamoDB(ctx context.Context) {
	// Get Redis address from environment variable
	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		log.Fatalf("Missing Redis address from environment variables")
	}

	// Initialize Redis client
	_, err := redis.InitRedisWithAddress(ctx, redisAddress)
	if err != nil {
		log.Fatalf("Error initializing Redis: %v", err)
	}

	// Get DynamoDB region from environment variable
	dynamoRegion := os.Getenv("DYNAMODB_REGION")
	if dynamoRegion == "" {
		log.Fatalf("Missing DynamoDB region from environment variables")
	}

	// Initialize DynamoDB client and set up table and replication
	db, err := dynamodb.NewDynamoDBClient(ctx, dynamoRegion)
	if err != nil {
		log.Fatalf("Error initializing DynamoDB: %v", err)
	}

	tableName := "Orders"
	if err := dynamodb.CreateTable(ctx, db, tableName); err != nil {
		log.Fatalf("Failed to create DynamoDB table: %v", err)
	}

	if err := dynamodb.EnableGlobalReplication(ctx, db, tableName, "us-west-2"); err != nil {
		log.Fatalf("Failed to enable global replication for table: %v", err)
	}

	log.Println("Global DynamoDB table created successfully with replication!")
}

// registerShopifyWebhook registers a product update webhook for Shopify
func registerShopifyWebhook() {
	shopName := os.Getenv("SHOP_NAME")
	accessToken := os.Getenv("SHOPIFY_ACCESS_TOKEN")
	webhookURL := os.Getenv("WEBHOOK_URL")
	if shopName == "" || accessToken == "" || webhookURL == "" {
		log.Fatalf("Missing environment variables for Shopify configuration")
	}

	if err := shopify.RegisterProductUpdateWebhook(shopName, accessToken, webhookURL); err != nil {
		log.Fatalf("Failed to register product update webhook: %v", err)
	}
	log.Println("Product update webhook registered successfully!")

	http.HandleFunc("/shopify/product/update", func(w http.ResponseWriter, r *http.Request) {
		shopify.HandleProductUpdateWebhook(w, r, nil, nil) // No need for Redis or DynamoDB in this call
	})

	go func() {
		log.Println("Starting HTTP server, waiting for webhooks...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()
}

// startMetricsServer starts the Prometheus metrics server
func startMetricsServer() {
	utils.InitMetricsServer()
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Println("Starting Prometheus metrics server on :8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start Prometheus metrics server: %v", err)
		}
	}()
}

// startKafka starts the Kafka consumer and producer
func startKafka(ctx context.Context) {
	writer := kafka_go.NewWriter(kafka_go.WriterConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "orders",
	})

	reader := kafka_go.NewReader(kafka_go.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "orders",
		GroupID: "order-consumer-group",
	})

	go func() {
		if err := kafka.ConsumeMessages(ctx, reader, nil, nil); err != nil {
			log.Fatalf("Error consuming Kafka messages: %v", err)
		}
	}()

	if err := kafka.ProduceMessages(ctx, writer); err != nil {
		log.Fatalf("Error producing Kafka messages: %v", err)
	}

	if err := writer.Close(); err != nil {
		log.Fatalf("Failed to close Kafka writer: %v", err)
	}
}

// setupLogging configures logging to a file
func setupLogging() {
	logFilePath := "./logs/app.log"
	f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Logging initialized")
}
