package domain

import (
	"time"
)

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductStock struct {
	ID          int64     `json:"id"`
	ProductID   int64     `json:"product_id"`
	Quantity    int       `json:"quantity"`
	Version     int       `json:"version"`
	LastUpdated time.Time `json:"last_updated"`
}

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
