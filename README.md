# shopify-scraper

A powerful Go library for efficiently scraping Shopify stores by leveraging their internal GraphQL API endpoints.

## Features

- Fetch product data directly from Shopify's internal APIs
- No HTML parsing required - pure JSON responses
- Rate limiting and error handling
- Clean, idiomatic Go interface

## Installation

```bash
go get github.com/yourusername/shopify-scraper
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/yourusername/shopify-scraper"
)

func main() {
    client := shopify.NewClient()
    
    // Get products from a Shopify store
    products, err := client.GetProducts("store-name.myshopify.com")
    if err != nil {
        panic(err)
    }
    
    for _, product := range products {
        fmt.Printf("Product: %s, Price: %s\n", product.Title, product.Price)
    }
}
```

## License

MIT License - See LICENSE file for details