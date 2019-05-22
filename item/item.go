package item

import (
	"time"
)

// Item is a model for a catalog item
type Item struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price" gorm:"type:DECIMAL(7,2)"`
	Count       int     `json:"count"`
	//ImageURL    []string   `json:"image_urls"`
	//Tags        []string   `json:"tags"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
