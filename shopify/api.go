package shopify

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// FetchProductDetails fetches the details of a specific product from Shopify
func FetchProductDetails(shop, accessToken, productID string) (string, error) {
	url := buildProductURL(shop, productID)

	req, err := buildRequest(url, accessToken)
	if err != nil {
		return "", fmt.Errorf("failed to build request: %v", err)
	}

	resp, err := executeRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		return "", err
	}

	log.Printf("Fetched product details for product %s", productID)
	return body, nil
}

// buildProductURL constructs the product URL for the Shopify API
func buildProductURL(shop, productID string) string {
	return fmt.Sprintf("https://%s.myshopify.com/admin/api/2023-01/products/%s.json", shop, productID)
}

// buildRequest creates a new HTTP GET request with the appropriate headers
func buildRequest(url, accessToken string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Shopify-Access-Token", accessToken)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// executeRequest sends the HTTP request and returns the response
func executeRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}

// readResponseBody reads the body from the HTTP response
func readResponseBody(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(body), nil
}
