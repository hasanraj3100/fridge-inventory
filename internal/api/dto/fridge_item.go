package dto

import (
	"github.com/hasanraj3100/fridge-inventory/internal/domain"
)

type FridgeItemAddRequest struct {
	Name      string              `json:"name" validate:"required,min=1,max=50"`
	Category  domain.FoodCategory `json:"category" validate:"required,oneof=Dairy Produce Meat Pantry Other"`
	Quantity  float32             `json:"quantity" validate:"required,gt=0"`
	Unit      string              `json:"unit" validate:"required,min=1,max=20"`
	UserID    int64               `json:"-"`
	BoughtAt  string              `json:"bought_at" validate:"required,datetime=2006-01-02"`
	ExpiresAt string              `json:"expires_at" validate:"required,datetime=2006-01-02"`
	MinThreshold float32             `json:"min_threshold" validate:"required,gte=0"`
}
