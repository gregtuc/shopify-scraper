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

	// Example Shopify store
	domain := "allbirds.com"

	fmt.Printf("Fetching products from %s...\n", domain)

	// Get all products
	products, err := client.GetProducts(domain)
	if err != nil {
		log.Fatalf("Error getting products: %v", err)
	}

	fmt.Printf("\nFound %d products\n", len(products))

	// Print products in a pretty format
	for _, product := range products {
		fmt.Printf("\nProduct: %s\n", product.Title)
		fmt.Printf("Handle: %s\n", product.Handle)
		fmt.Printf("Vendor: %s\n", product.Vendor)
		fmt.Printf("Type: %s\n", product.ProductType)

		if len(product.Variants) > 0 {
			variant := product.Variants[0]
			fmt.Printf("Price: %s\n", variant.Price)
			if variant.CompareAtPrice != "" {
				fmt.Printf("Compare at price: %s\n", variant.CompareAtPrice)
			}
			fmt.Printf("SKU: %s\n", variant.SKU)
			fmt.Printf("Inventory: %d\n", variant.InventoryQuantity)
		}

		if len(product.Images) > 0 {
			fmt.Printf("First image URL: %s\n", product.Images[0].Src)
		}

		if len(product.Options) > 0 {
			fmt.Printf("Options:\n")
			for _, opt := range product.Options {
				fmt.Printf("  - %s: %v\n", opt.Name, opt.Values)
			}
		}

		fmt.Printf("Tags: %v\n", product.Tags)
		fmt.Println("----------------------------------------")
	}

	// Get a single product by handle (using the first product's handle if available)
	if len(products) > 0 {
		handle := products[0].Handle
		fmt.Printf("\nFetching detailed information for product: %s\n", handle)

		product, err := client.GetProduct(domain, handle)
		if err != nil {
			log.Printf("Error getting product details: %v", err)
		} else {
			// Pretty print the product as JSON
			productJSON, _ := json.MarshalIndent(product, "", "  ")
			fmt.Printf("\nSingle Product Details:\n%s\n", string(productJSON))
		}
	}
}
