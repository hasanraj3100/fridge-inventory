package service

import (
	"context"
	"fmt"

	"github.com/hasanraj3100/fridge-inventory/internal/api/dto"
	"github.com/hasanraj3100/fridge-inventory/internal/domain"
	"github.com/hasanraj3100/fridge-inventory/internal/repository"
)

type FridgeItemService interface {
	AddItem(
		ctx context.Context,
		params dto.FridgeItemAddRequest,
	) (*domain.FridgeItem, error)
}

type fridgeItemService struct {
	fridgeItemRepo repository.FridgeItemRepository
}

func NewFridgeItemService(fridgeItemRepo repository.FridgeItemRepository) FridgeItemService {
	return &fridgeItemService{
		fridgeItemRepo: fridgeItemRepo,
	}
}

func (s *fridgeItemService) AddItem(ctx context.Context, params dto.FridgeItemAddRequest) (*domain.FridgeItem, error) {
	newItem := &domain.FridgeItem{
		Name:      params.Name,
		Category:  params.Category,
		Quantity:  params.Quantity,
		Unit:      params.Unit,
		UserID:    params.UserID,
		BoughtAt:  params.BoughtAt,
		ExpiresAt: params.ExpiresAt,
		MinStock:  params.MinStock,
	}

	createdItem, err := s.fridgeItemRepo.Create(ctx, newItem)
	if err != nil {
		return nil, fmt.Errorf("failed to add fridge item: %w", err)
	}

	return createdItem, nil
}
