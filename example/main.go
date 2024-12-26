package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gregtuc/shopify-scraper"
)

func main() {
	// Create a new client with custom timeout and page size
	client := shopify.NewClient(
		shopify.WithTimeout(10*time.Second),
		shopify.WithPageSize(50), // Smaller page size for demonstration
	)

	// Example Shopify store
	domain := "allbirds.com"

	// 1. First, let's get all collections
	fmt.Printf("Fetching collections from %s...\n", domain)
	collections, err := client.GetCollections(domain)
	if err != nil {
		log.Printf("Error getting collections: %v\n", err)
	} else {
		fmt.Printf("\nFound %d collections\n", len(collections))
		for _, collection := range collections {
			fmt.Printf("\nCollection: %s\n", collection.Title)
			fmt.Printf("Handle: %s\n", collection.Handle)
			if collection.Image != nil {
				fmt.Printf("Image: %s\n", collection.Image.Src)
			}
			fmt.Println("----------------------------------------")
		}

		// 2. If we found any collections, let's get products from the first one
		if len(collections) > 0 {
			collection := collections[0]
			fmt.Printf("\nFetching products from collection: %s\n", collection.Title)
			products, err := client.GetCollectionProducts(domain, collection.Handle)
			if err != nil {
				log.Printf("Error getting collection products: %v\n", err)
			} else {
				fmt.Printf("\nFound %d products in collection\n", len(products))
				printProducts(products)
			}
		}
	}

	// 3. Let's try searching for products
	searchQuery := "wool"
	fmt.Printf("\nSearching for products with query: %s\n", searchQuery)
	searchResults, err := client.SearchProducts(domain, searchQuery)
	if err != nil {
		log.Printf("Error searching products: %v\n", err)
	} else {
		fmt.Printf("\nFound %d products matching search\n", len(searchResults))
		printProducts(searchResults)
	}

	// 4. Finally, let's get all products (with pagination)
	fmt.Printf("\nFetching all products from %s...\n", domain)
	allProducts, err := client.GetProducts(domain)
	if err != nil {
		log.Printf("Error getting all products: %v\n", err)
	} else {
		fmt.Printf("\nFound %d total products\n", len(allProducts))
		printProducts(allProducts)

		// 5. Get detailed information for the first product
		if len(allProducts) > 0 {
			handle := allProducts[0].Handle
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
}

func printProducts(products []shopify.Product) {
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
			if variant.Available {
				fmt.Printf("Status: Available\n")
			} else {
				fmt.Printf("Status: Out of Stock\n")
			}
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
}
