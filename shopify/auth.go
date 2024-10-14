package shopify

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// OAuthURL generates the authentication URL for Shopify OAuth flow
func OAuthURL(apiKey, shopName string) string {
	redirectURI := url.QueryEscape("https://your-app.com/shopify/callback")
	return fmt.Sprintf("https://%s.myshopify.com/admin/oauth/authorize?client_id=%s&scope=read_products,write_products&redirect_uri=%s", shopName, apiKey, redirectURI)
}

// HandleOAuthCallback handles the OAuth callback from Shopify and exchanges the code for an access token
func HandleOAuthCallback(w http.ResponseWriter, r *http.Request, apiKey, apiSecret string) {
	code, shop, err := extractOAuthParams(r)
	if err != nil {
		log.Printf("Invalid OAuth parameters: %v", err)
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
		return
	}

	token, err := exchangeCodeForToken(shop, code, apiKey, apiSecret)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		http.Error(w, "Failed to authenticate with Shopify", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully authenticated shop %s with token: %s", shop, token)
	fmt.Fprintf(w, "Shop %s authenticated", shop)
}

// extractOAuthParams extracts the required parameters from the OAuth callback request
func extractOAuthParams(r *http.Request) (string, string, error) {
	query := r.URL.Query()
	code := query.Get("code")
	shop := query.Get("shop")

	if code == "" || shop == "" {
		return "", "", fmt.Errorf("missing code or shop parameter")
	}
	return code, shop, nil
}

// exchangeCodeForToken exchanges the authorization code for an access token
func exchangeCodeForToken(shop, code, apiKey, apiSecret string) (string, error) {
	tokenURL := fmt.Sprintf("https://%s.myshopify.com/admin/oauth/access_token", shop)
	data := url.Values{
		"client_id":     {apiKey},
		"client_secret": {apiSecret},
		"code":          {code},
	}

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return "", fmt.Errorf("failed to send token exchange request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // io.ReadAll is preferred over ioutil.ReadAll
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %v", err)
	}

	// Extract the token - consider using json decoding for a real application
	return string(body), nil
}
