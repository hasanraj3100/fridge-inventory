package service

import (
	"context"
	"fmt"
	"time"

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
	boughtAt, err := time.Parse("2006-01-02", params.BoughtAt)
	if err != nil {
		return nil, fmt.Errorf("invalid bought_at date format: %w", err)
	}

	expiresAt, err := time.Parse("2006-01-02", params.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid expires_at date format: %w", err)
	}

	newItem := &domain.FridgeItem{
		Name:      params.Name,
		Category:  params.Category,
		Quantity:  params.Quantity,
		Unit:      params.Unit,
		UserID:    params.UserID,
		BoughtAt:  boughtAt,
		ExpiresAt: expiresAt,
		MinStock:  params.MinStock,
	}

	createdItem, err := s.fridgeItemRepo.Create(ctx, newItem)
	if err != nil {
		return nil, fmt.Errorf("failed to add fridge item: %w", err)
	}

	return createdItem, nil
}
