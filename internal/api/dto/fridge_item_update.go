package dto

import "github.com/hasanraj3100/fridge-inventory/internal/domain"

type FridgeItemUpdateRequest struct {
	Name         *string              `json:"name,omitempty" validate:"omitempty,min=1,max=50"`
	Category     *domain.FoodCategory `json:"category,omitempty" validate:"omitempty,oneof=Dairy Produce Meat Pantry Other"`
	Quantity     *float32             `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Unit         *string              `json:"unit,omitempty" validate:"omitempty,min=1,max=20"`
	BoughtAt     *string              `json:"bought_at,omitempty" validate:"omitempty,datetime=2006-01-02"`
	ExpiresAt    *string              `json:"expires_at,omitempty" validate:"omitempty,datetime=2006-01-02"`
	MinThreshold *float32             `json:"min_threshold,omitempty" validate:"omitempty,gte=0"`
}
