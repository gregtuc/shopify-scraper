package shopify

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient(
		WithTimeout(5*time.Second),
		WithUserAgent("test-agent"),
		WithPageSize(50),
	)

	if client.pageSize != 50 {
		t.Errorf("expected pageSize 50, got %d", client.pageSize)
	}

	if client.userAgent != "test-agent" {
		t.Errorf("expected userAgent test-agent, got %s", client.userAgent)
	}

	if client.httpClient.Timeout != 5*time.Second {
		t.Errorf("expected timeout 5s, got %s", client.httpClient.Timeout)
	}
}

func TestFormatDomain(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://www.example.com", "example.com"},
		{"http://example.com", "example.com"},
		{"www.example.com", "example.com"},
		{"example.com", "example.com"},
	}

	for _, test := range tests {
		result := formatDomain(test.input)
		if result != test.expected {
			t.Errorf("formatDomain(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func setupTestServer(t *testing.T, path string, response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			t.Errorf("Expected path %s, got %s", path, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
}

func TestGetProducts(t *testing.T) {
	response := `{
		"products": [
			{
				"id": 123,
				"title": "Test Product",
				"handle": "test-product",
				"variants": [
					{
						"id": 456,
						"price": "19.99"
					}
				]
			}
		]
	}`

	server := setupTestServer(t, "/products.json", response)
	defer server.Close()

	client := NewClient()
	products, err := client.GetProducts(server.URL)

	if err != nil {
		t.Fatalf("GetProducts returned error: %v", err)
	}

	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}

	product := products[0]
	if product.ID != 123 {
		t.Errorf("expected product ID 123, got %d", product.ID)
	}

	if product.Title != "Test Product" {
		t.Errorf("expected product title 'Test Product', got %s", product.Title)
	}
}

func TestGetProduct(t *testing.T) {
	response := `{
		"product": {
			"id": 123,
			"title": "Test Product",
			"handle": "test-product",
			"variants": [
				{
					"id": 456,
					"price": "19.99"
				}
			]
		}
	}`

	server := setupTestServer(t, "/products/test-product.json", response)
	defer server.Close()

	client := NewClient()
	product, err := client.GetProduct(server.URL, "test-product")

	if err != nil {
		t.Fatalf("GetProduct returned error: %v", err)
	}

	if product.ID != 123 {
		t.Errorf("expected product ID 123, got %d", product.ID)
	}

	if product.Title != "Test Product" {
		t.Errorf("expected product title 'Test Product', got %s", product.Title)
	}
}

func TestGetCollections(t *testing.T) {
	response := `{
		"collections": [
			{
				"id": 123,
				"title": "Test Collection",
				"handle": "test-collection"
			}
		]
	}`

	server := setupTestServer(t, "/collections.json", response)
	defer server.Close()

	client := NewClient()
	collections, err := client.GetCollections(server.URL)

	if err != nil {
		t.Fatalf("GetCollections returned error: %v", err)
	}

	if len(collections) != 1 {
		t.Fatalf("expected 1 collection, got %d", len(collections))
	}

	collection := collections[0]
	if collection.ID != 123 {
		t.Errorf("expected collection ID 123, got %d", collection.ID)
	}

	if collection.Title != "Test Collection" {
		t.Errorf("expected collection title 'Test Collection', got %s", collection.Title)
	}
}

func TestSearchProducts(t *testing.T) {
	response := `{
		"resources": {
			"results": {
				"products": [
					{
						"id": 123,
						"title": "Test Product",
						"handle": "test-product"
					}
				]
			}
		}
	}`

	server := setupTestServer(t, "/search/suggest.json", response)
	defer server.Close()

	client := NewClient()
	products, err := client.SearchProducts(server.URL, "test")

	if err != nil {
		t.Fatalf("SearchProducts returned error: %v", err)
	}

	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}

	product := products[0]
	if product.ID != 123 {
		t.Errorf("expected product ID 123, got %d", product.ID)
	}

	if product.Title != "Test Product" {
		t.Errorf("expected product title 'Test Product', got %s", product.Title)
	}
}
