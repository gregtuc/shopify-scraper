# shopify-scraper

A powerful Go library for scraping Shopify stores by leveraging their client-side APIs. This library provides a clean, efficient way to fetch product data, collections, and more from any Shopify store without requiring API keys or authentication.

## Features

- 🚀 **Zero Authentication Required**: Works with any Shopify store's public endpoints
- 📦 **Comprehensive Data**: Full product details, variants, images, collections, and more
- 🔄 **Automatic Pagination**: Handles stores with any number of products
- 🔍 **Search Support**: Use the store's native search functionality
- 📑 **Collection Support**: Browse and fetch products by collection
- 🛡️ **Rate Limiting**: Built-in delays to be respectful to servers
- 🎯 **Type Safety**: Full Go types for all Shopify data structures
- 🔧 **Configurable**: Customize timeouts, page sizes, and user agents

## Installation

```bash
go get github.com/gtucker/shopify-scraper
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/gtucker/shopify-scraper"
)

func main() {
    // Create a new client
    client := shopify.NewClient()
    
    // Get products from a Shopify store
    products, err := client.GetProducts("allbirds.com")
    if err != nil {
        panic(err)
    }
    
    // Print product details
    for _, product := range products {
        fmt.Printf("Product: %s\n", product.Title)
        if len(product.Variants) > 0 {
            fmt.Printf("Price: %s\n", product.Variants[0].Price)
        }
    }
}
```

## Advanced Usage

### Client Configuration

```go
client := shopify.NewClient(
    shopify.WithTimeout(10 * time.Second),
    shopify.WithPageSize(50),
    shopify.WithUserAgent("Custom User Agent"),
)
```

### Fetching Collections

```go
collections, err := client.GetCollections("allbirds.com")
if err != nil {
    panic(err)
}

for _, collection := range collections {
    fmt.Printf("Collection: %s\n", collection.Title)
    
    // Get products in this collection
    products, err := client.GetCollectionProducts("allbirds.com", collection.Handle)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Found %d products in collection\n", len(products))
}
```

### Searching Products

```go
products, err := client.SearchProducts("allbirds.com", "wool runners")
if err != nil {
    panic(err)
}
fmt.Printf("Found %d matching products\n", len(products))
```

### Getting Single Product

```go
product, err := client.GetProduct("allbirds.com", "mens-wool-runners")
if err != nil {
    panic(err)
}
fmt.Printf("Product: %s\n", product.Title)
```

## Available Methods

- `GetProducts(domain string) ([]Product, error)`
- `GetProduct(domain, handle string) (*Product, error)`
- `GetCollections(domain string) ([]Collection, error)`
- `GetCollectionProducts(domain, collectionHandle string) ([]Product, error)`
- `SearchProducts(domain, query string) ([]Product, error)`

## Data Models

### Product
```go
type Product struct {
    ID          int64
    Title       string
    Handle      string
    Description string
    Vendor      string
    ProductType string
    Tags        []string
    Variants    []Variant
    Images      []Image
    Options     []Option
    // ... and more
}
```

### Variant
```go
type Variant struct {
    ID                int64
    Title            string
    Price            string
    SKU              string
    CompareAtPrice   string
    InventoryQuantity int
    // ... and more
}
```

See [models.go](models.go) for complete type definitions.

## Best Practices

1. **Rate Limiting**: The library includes a built-in 100ms delay between requests. Adjust the page size for your needs.
2. **Error Handling**: Always check error returns. The library provides detailed error messages.
3. **Domain Format**: Domains can be provided with or without protocol/www (e.g., "store.com" or "https://www.store.com").

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - See LICENSE file for details

## Acknowledgments

This library is designed to be respectful to Shopify stores by:
- Using the same endpoints as their web frontend
- Including reasonable delays between requests
- Properly identifying itself with user agents
- Not attempting to bypass any rate limiting