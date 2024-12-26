package shopify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client represents a Shopify scraper client
type Client struct {
	httpClient *http.Client
	userAgent  string
	pageSize   int
}

// ClientOption is a function that modifies the client
type ClientOption func(*Client)

// NewClient creates a new Shopify scraper client
func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		pageSize:  250, // Maximum allowed by Shopify
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// WithTimeout sets the client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithUserAgent sets the client user agent
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithPageSize sets the number of products per page
func WithPageSize(size int) ClientOption {
	return func(c *Client) {
		if size > 250 {
			size = 250 // Shopify's maximum
		}
		c.pageSize = size
	}
}

// formatDomain ensures the domain is in the correct format
func formatDomain(domain string) string {
	// Remove protocol if present
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")

	// Remove www. if present
	domain = strings.TrimPrefix(domain, "www.")

	return domain
}

// makeRequest makes an HTTP request with proper headers
func (c *Client) makeRequest(method, url string, referer string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", referer)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	return c.httpClient.Do(req)
}

// GetProducts fetches all products from a Shopify store with pagination
func (c *Client) GetProducts(domain string) ([]Product, error) {
	domain = formatDomain(domain)
	var allProducts []Product
	page := 1

	for {
		url := fmt.Sprintf("https://%s/products.json?limit=%d&page=%d", domain, c.pageSize, page)
		referer := fmt.Sprintf("https://%s/", domain)

		resp, err := c.makeRequest("GET", url, referer)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		var response struct {
			Products []Product `json:"products"`
		}

		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %w", err)
		}

		if len(response.Products) == 0 {
			break // No more products
		}

		allProducts = append(allProducts, response.Products...)
		page++

		// Add a small delay to be nice to the server
		time.Sleep(100 * time.Millisecond)
	}

	return allProducts, nil
}

// GetProduct fetches a single product by handle
func (c *Client) GetProduct(domain, handle string) (*Product, error) {
	domain = formatDomain(domain)
	url := fmt.Sprintf("https://%s/products/%s.json", domain, handle)
	referer := fmt.Sprintf("https://%s/products/%s", domain, handle)

	resp, err := c.makeRequest("GET", url, referer)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var response struct {
		Product Product `json:"product"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &response.Product, nil
}

// GetCollections fetches all collections from a Shopify store
func (c *Client) GetCollections(domain string) ([]Collection, error) {
	domain = formatDomain(domain)
	url := fmt.Sprintf("https://%s/collections.json", domain)
	referer := fmt.Sprintf("https://%s/", domain)

	resp, err := c.makeRequest("GET", url, referer)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var response struct {
		Collections []Collection `json:"collections"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return response.Collections, nil
}

// GetCollectionProducts fetches all products in a collection
func (c *Client) GetCollectionProducts(domain string, collectionHandle string) ([]Product, error) {
	domain = formatDomain(domain)
	var allProducts []Product
	page := 1

	for {
		url := fmt.Sprintf("https://%s/collections/%s/products.json?limit=%d&page=%d",
			domain, collectionHandle, c.pageSize, page)
		referer := fmt.Sprintf("https://%s/collections/%s", domain, collectionHandle)

		resp, err := c.makeRequest("GET", url, referer)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		var response struct {
			Products []Product `json:"products"`
		}

		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %w", err)
		}

		if len(response.Products) == 0 {
			break // No more products
		}

		allProducts = append(allProducts, response.Products...)
		page++

		// Add a small delay to be nice to the server
		time.Sleep(100 * time.Millisecond)
	}

	return allProducts, nil
}

// SearchProducts searches for products using the store's search endpoint
func (c *Client) SearchProducts(domain, query string) ([]Product, error) {
	domain = formatDomain(domain)
	url := fmt.Sprintf("https://%s/search/suggest.json?q=%s&resources[type]=product",
		domain, query)
	referer := fmt.Sprintf("https://%s/search?q=%s", domain, query)

	resp, err := c.makeRequest("GET", url, referer)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var response struct {
		Resources struct {
			Results struct {
				Products []Product `json:"products"`
			} `json:"results"`
		} `json:"resources"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return response.Resources.Results.Products, nil
}
