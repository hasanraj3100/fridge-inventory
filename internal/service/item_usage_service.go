package service

import (
	"context"
	"fmt"

	"github.com/hasanraj3100/fridge-inventory/internal/api/dto"
	"github.com/hasanraj3100/fridge-inventory/internal/db"
	"github.com/hasanraj3100/fridge-inventory/internal/domain"
	"github.com/hasanraj3100/fridge-inventory/internal/repository"
)

type ItemUsageService interface {
	CreateItemUsage(ctx context.Context, userID int64, params dto.ItemUsageCreateRequest) error
}

type itemUsageService struct {
	itemUsageRepo repository.ItemUsageRepository
	itemRepo      repository.FridgeItemRepository
	txManager     repository.TransactionProvider
}

func NewItemUsageService(
	itemUsageRepo repository.ItemUsageRepository,
	itemRepo repository.FridgeItemRepository,
	txManager repository.TransactionProvider,
) ItemUsageService {
	return &itemUsageService{
		itemUsageRepo: itemUsageRepo,
		itemRepo:      itemRepo,
		txManager:     txManager,
	}
}

func (s *itemUsageService) CreateItemUsage(ctx context.Context, userID int64, params dto.ItemUsageCreateRequest) error {
	err := s.txManager.WithinTransaction(ctx, func(tx db.DBTX) error {
		fridgeItemRepo := repository.NewFridgeItemRepository(tx)
		itemUsageRepo := repository.NewItemUsageRepository(tx)

		item, err := fridgeItemRepo.GetByID(ctx, params.ItemID)
		if err != nil {
			return err
		}
		if item == nil || item.UserID != userID {
			return fmt.Errorf("item not found or does not belong to user")
		}

		if params.QuantityUsed > item.Quantity {
			return fmt.Errorf("quantity used is greater than available quantity")
		}

		item.Quantity -= params.QuantityUsed

		err = fridgeItemRepo.Update(ctx, item)
		if err != nil {
			return err
		}

		itemUsage := &domain.ItemUsage{
			ItemID:       params.ItemID,
			QuantityUsed: params.QuantityUsed,
			Reason:       domain.ItemUsageReason(params.Reason),
		}

		err = itemUsageRepo.Create(ctx, itemUsage)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
