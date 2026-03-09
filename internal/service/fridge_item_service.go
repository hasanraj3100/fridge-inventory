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
	GetByUserID(ctx context.Context, userID int64) ([]*domain.FridgeItem, error)
	UpdateItem(ctx context.Context, itemID int64, userID int64, params dto.FridgeItemUpdateRequest) (*domain.FridgeItem, error)
}

type fridgeItemService struct {
	fridgeItemRepo repository.FridgeItemRepository
	txManager      repository.TransactionProvider
}

func NewFridgeItemService(fridgeItemRepo repository.FridgeItemRepository, txManager repository.TransactionProvider) FridgeItemService {
	return &fridgeItemService{
		fridgeItemRepo: fridgeItemRepo,
		txManager:      txManager,
	}
}

func (s *fridgeItemService) AddItem(ctx context.Context, params dto.FridgeItemAddRequest) (*domain.FridgeItem, error) {
	boughtAt, err := time.Parse(time.DateOnly, params.BoughtAt)
	if err != nil {
		return nil, fmt.Errorf("invalid bought_at date format: %w", err)
	}

	expiresAt, err := time.Parse(time.DateOnly, params.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid expires_at date format: %w", err)
	}

	newItem := &domain.FridgeItem{
		Name:         params.Name,
		Category:     params.Category,
		Quantity:     params.Quantity,
		Unit:         params.Unit,
		UserID:       params.UserID,
		BoughtAt:     boughtAt,
		ExpiresAt:    expiresAt,
		MinThreshold: params.MinThreshold,
	}

	createdItem, err := s.fridgeItemRepo.Create(ctx, newItem)
	if err != nil {
		return nil, fmt.Errorf("failed to add fridge item: %w", err)
	}

	return createdItem, nil
}

func (s *fridgeItemService) GetByUserID(ctx context.Context, userID int64) ([]*domain.FridgeItem, error) {
	return s.fridgeItemRepo.GetByUserID(ctx, userID)
}

func (s *fridgeItemService) UpdateItem(ctx context.Context, itemID int64, userID int64, params dto.FridgeItemUpdateRequest) (*domain.FridgeItem, error) {
	// Check if item exists and belongs to user
	exists, err := s.fridgeItemRepo.BelongsToUser(ctx, itemID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify item ownership: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("item not found or does not belong to user")
	}

	// Build update item with only changed fields
	updateItem := &domain.FridgeItem{
		ID: int(itemID),
	}

	if params.Name != nil {
		updateItem.Name = *params.Name
	}

	if params.Category != nil {
		updateItem.Category = *params.Category
	}

	if params.Quantity != nil {
		updateItem.Quantity = *params.Quantity
	}

	if params.Unit != nil {
		updateItem.Unit = *params.Unit
	}

	if params.BoughtAt != nil {
		boughtAt, err := time.Parse(time.DateOnly, *params.BoughtAt)
		if err != nil {
			return nil, fmt.Errorf("invalid bought_at date format: %w", err)
		}
		updateItem.BoughtAt = boughtAt
	}

	if params.ExpiresAt != nil {
		expiresAt, err := time.Parse(time.DateOnly, *params.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("invalid expires_at date format: %w", err)
		}
		updateItem.ExpiresAt = expiresAt
	}

	if params.MinThreshold != nil {
		updateItem.MinThreshold = *params.MinThreshold
	}

	// Update the item
	err = s.fridgeItemRepo.Update(ctx, updateItem)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	// Return updated item
	updatedItem, err := s.fridgeItemRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated item: %w", err)
	}

	return updatedItem, nil
}
