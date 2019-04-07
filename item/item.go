package item

// Item is a model for a catalog item
type Item struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Count       int      `json:"count"`
	ImageURL    []string `json:"image_urls"`
	Tags        []string `json:"tags"`
}
