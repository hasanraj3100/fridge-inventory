package service

import (
	"context"
	"errors"
	"time"

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
	GetItemUsageByUserID(ctx context.Context, userID int64, page int, pageSize int) (*dto.PaginatedItemUsageResponse, error)
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

func (s *itemUsageService) GetItemUsageByUserID(ctx context.Context, userID int64, page int, pageSize int) (*dto.PaginatedItemUsageResponse, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get total count
	totalItems, err := s.itemUsageRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get paginated data
	results, err := s.itemUsageRepo.GetByUserIDWithPagination(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	var items []dto.ItemUsageWithDetails
	for _, result := range results {
		items = append(items, dto.ItemUsageWithDetails{
			ID:           result["id"].(int64),
			ItemID:       result["item_id"].(int64),
			ItemName:     result["item_name"].(string),
			QuantityUsed: result["quantity_used"].(float32),
			Unit:         result["unit"].(string),
			Reason:       result["reason"].(string),
			UsedAt:       result["used_at"].(time.Time).Format(time.RFC3339),
		})
	}

	// Calculate total pages
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize != 0 {
		totalPages++
	}

	return &dto.PaginatedItemUsageResponse{
		Data: items,
		Pagination: dto.PaginationResponse{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}, nil
}
