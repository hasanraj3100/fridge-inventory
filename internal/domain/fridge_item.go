// Package domain defines core domain models for the application.
package domain

import "time"

type FoodCategory string

var (
	CategoryDairy   FoodCategory = "Dairy"
	CategoryProduce FoodCategory = "Produce"
	CategoryMeat    FoodCategory = "Meat"
	CategoryPantry  FoodCategory = "Pantry"
	CategoryOther   FoodCategory = "Other"
)

type FridgeItem struct {
	ID        int          `json:"id" db:"id"`
	Name      string       `json:"name" db:"name"`
	Category  FoodCategory `json:"category" db:"category"`
	Quantity  float32      `json:"quantity" db:"quantity"`
	Unit      string       `json:"unit" db:"unit"`
	UserID    int64        `json:"user_id" db:"user_id"`
	BoughtAt  time.Time    `json:"bought_at" db:"bought_at"`
	ExpiresAt time.Time    `json:"expires_at" db:"expires_at"`
	MinStock  float32      `json:"min_stock" db:"min_stock"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
}
