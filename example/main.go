package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gregtuc/shopify-scraper"
)

func main() {
	// Create a new client with custom timeout
	client := shopify.NewClient(
		shopify.WithTimeout(10 * time.Second),
	)

	// Example Shopify store (replace with an actual store domain)
	domain := "example-store.myshopify.com"

	// Get all products
	products, err := client.GetProducts(domain)
	if err != nil {
		log.Fatalf("Error getting products: %v", err)
	}

	// Print products in a pretty format
	for _, product := range products {
		fmt.Printf("\nProduct: %s\n", product.Title)
		fmt.Printf("Handle: %s\n", product.Handle)
		fmt.Printf("Vendor: %s\n", product.Vendor)
		fmt.Printf("Type: %s\n", product.ProductType)

		if len(product.Variants) > 0 {
			fmt.Printf("Price: %s\n", product.Variants[0].Price)
		}

		fmt.Printf("Tags: %v\n", product.Tags)
		fmt.Println("----------------------------------------")
	}

	// Get a single product by handle
	product, err := client.GetProduct(domain, "example-product-handle")
	if err != nil {
		log.Fatalf("Error getting product: %v", err)
	}

	// Pretty print the product as JSON
	productJSON, _ := json.MarshalIndent(product, "", "  ")
	fmt.Printf("\nSingle Product Details:\n%s\n", string(productJSON))
}
