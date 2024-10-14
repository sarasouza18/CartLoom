package shopify

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-redis/redis/v8"
)

// RegisterProductUpdateWebhook registers a product update webhook for Shopify
func RegisterProductUpdateWebhook(shop, accessToken, webhookURL string) error {
	webhookData := buildWebhookData(webhookURL)
	url := buildWebhookURL(shop)

	req, err := buildWebhookRequest(url, webhookData, accessToken)
	if err != nil {
		return fmt.Errorf("failed to build webhook request: %v", err)
	}

	resp, err := sendWebhookRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to register webhook: %s", string(body))
	}

	log.Printf("Product update webhook registered successfully for shop %s", shop)
	return nil
}

// HandleProductUpdateWebhook processes product update webhooks from Shopify
func HandleProductUpdateWebhook(w http.ResponseWriter, r *http.Request, rdb *redis.Client, db *dynamodb.Client) {
	body, err := readRequestBody(r)
	if err != nil {
		log.Printf("Error reading webhook request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	log.Printf("Received product update webhook: %s", string(body))

	ctx := context.Background()

	if err := updateProductInRedis(ctx, rdb, "product-id-123", "updated-product-data"); err != nil {
		log.Printf("Error updating product in Redis: %v", err)
	}

	if err := updateProductInDynamoDB(ctx, db, "123", "updated-product-data"); err != nil {
		log.Printf("Error updating product in DynamoDB: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product update processed"))
}

// buildWebhookData creates the JSON payload for the webhook
func buildWebhookData(webhookURL string) string {
	return fmt.Sprintf(`{
		"webhook": {
			"topic": "products/update",
			"address": "%s",
			"format": "json"
		}
	}`, webhookURL)
}

// buildWebhookURL constructs the Shopify webhook URL
func buildWebhookURL(shop string) string {
	return fmt.Sprintf("https://%s.myshopify.com/admin/api/2023-01/webhooks.json", shop)
}

// buildWebhookRequest creates a POST request for the webhook
func buildWebhookRequest(url, webhookData, accessToken string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(webhookData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Shopify-Access-Token", accessToken)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// sendWebhookRequest sends the webhook request to Shopify
func sendWebhookRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send webhook request: %v", err)
	}
	return resp, nil
}

// readRequestBody reads and returns the body of an HTTP request
func readRequestBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return body, nil
}

// updateProductInRedis simulates updating a product in Redis
func updateProductInRedis(ctx context.Context, rdb *redis.Client, productID, data string) error {
	err := rdb.Set(ctx, productID, data, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to update product in Redis: %v", err)
	}
	log.Printf("Product %s updated in Redis successfully", productID)
	return nil
}

// updateProductInDynamoDB simulates updating a product in DynamoDB
func updateProductInDynamoDB(ctx context.Context, db *dynamodb.Client, productID, data string) error {
	_, err := db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("Products"),
		Item: map[string]types.AttributeValue{
			"ProductID": &types.AttributeValueMemberS{Value: productID},
			"Data":      &types.AttributeValueMemberS{Value: data},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update product in DynamoDB: %v", err)
	}
	log.Printf("Product %s updated in DynamoDB successfully", productID)
	return nil
}
