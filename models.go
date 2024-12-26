package shopify

import "time"

// Product represents a Shopify product
type Product struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Handle      string    `json:"handle"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Vendor      string    `json:"vendor"`
	ProductType string    `json:"product_type"`
	Tags        []string  `json:"tags"`
	Variants    []Variant `json:"variants"`
	Images      []Image   `json:"images"`
	Options     []Option  `json:"options"`
}

// Variant represents a product variant
type Variant struct {
	ID                  int64     `json:"id"`
	ProductID           int64     `json:"product_id"`
	Title               string    `json:"title"`
	Price               string    `json:"price"`
	SKU                 string    `json:"sku"`
	Position            int       `json:"position"`
	CompareAtPrice      string    `json:"compare_at_price"`
	FulfillmentService  string    `json:"fulfillment_service"`
	InventoryManagement string    `json:"inventory_management"`
	Option1             string    `json:"option1"`
	Option2             string    `json:"option2"`
	Option3             string    `json:"option3"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Taxable             bool      `json:"taxable"`
	Barcode             string    `json:"barcode"`
	Grams               int       `json:"grams"`
	Weight              float64   `json:"weight"`
	WeightUnit          string    `json:"weight_unit"`
	InventoryQuantity   int       `json:"inventory_quantity"`
}

// Image represents a product image
type Image struct {
	ID         int64     `json:"id"`
	ProductID  int64     `json:"product_id"`
	Position   int       `json:"position"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	Src        string    `json:"src"`
	VariantIDs []int64   `json:"variant_ids"`
}

// Option represents a product option
type Option struct {
	ID        int64    `json:"id"`
	ProductID int64    `json:"product_id"`
	Name      string   `json:"name"`
	Position  int      `json:"position"`
	Values    []string `json:"values"`
}
