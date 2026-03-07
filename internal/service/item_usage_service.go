package service

import (
	"context"
	"errors"

	"github.com/hasanraj3100/fridge-inventory/internal/api/dto"
	"github.com/hasanraj3100/fridge-inventory/internal/db"
	"github.com/hasanraj3100/fridge-inventory/internal/domain"
	"github.com/hasanraj3100/fridge-inventory/internal/repository"
)

var (
	ErrItemNotFound         = errors.New("item not found")
	ErrUnauthorized         = errors.New("item does not belong to user")
	ErrInsufficientQuantity = errors.New("quantity used is greater than available quantity")
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
			return ErrItemNotFound
		}
		if item.UserID != userID {
			return ErrUnauthorized
		}

		if params.QuantityUsed > item.Quantity {
			return ErrInsufficientQuantity
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
