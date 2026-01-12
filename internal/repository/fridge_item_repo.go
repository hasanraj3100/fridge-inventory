package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/hasanraj3100/fridge-inventory/internal/domain"
	"github.com/jmoiron/sqlx"
)

type FridgeItemRepository interface {
	Create(ctx context.Context, item *domain.FridgeItem) (*domain.FridgeItem, error)
}

type fridgeItemRepository struct {
	DB *sqlx.DB
}

func NewFridgeItemRepository(db *sqlx.DB) FridgeItemRepository {
	return &fridgeItemRepository{
		DB: db,
	}
}

func (repo *fridgeItemRepository) Create(ctx context.Context, item *domain.FridgeItem) (*domain.FridgeItem, error) {
	item.CreatedAt = time.Now().UTC()
	item.UpdatedAt = time.Now().UTC()

	query := `INSERT INTO fridge_items (
	name, category, quantity, unit, user_id, bought_at, expires_at, min_stock, created_at, updated_at
	)
	VALUES (
	:name, :category, :quantity, :unit, :user_id, :bought_at, :expires_at, :min_stock, :created_at, :updated_at
	) RETURNING id`

	res, err := repo.DB.NamedQueryContext(ctx, query, item)
	if err != nil {
		return nil, fmt.Errorf("failed to insert item to database: %w", err)
	}
	defer res.Close()

	if res.Next() {
		err = res.Scan(&item.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve inserted item ID: %w", err)
		}
	}

	return item, nil
}
