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
	GetByUserID(ctx context.Context, userID int64) ([]*domain.FridgeItem, error)
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
	name, category, quantity, unit, user_id, bought_at, expires_at, min_threshold, created_at, updated_at
	)
	VALUES (
	:name, :category, :quantity, :unit, :user_id, :bought_at, :expires_at, :min_threshold, :created_at, :updated_at
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

func (repo *fridgeItemRepository) GetByUserID(ctx context.Context, userID int64) ([]*domain.FridgeItem, error) {
	query := `SELECT * FROM fridge_items WHERE user_id = $1`

	rows, err := repo.DB.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query fridge items: %w", err)
	}

	defer rows.Close()

	var items []*domain.FridgeItem = make([]*domain.FridgeItem, 0)

	for rows.Next() {
		var item domain.FridgeItem
		if err := rows.StructScan(&item); err != nil {
			return nil, fmt.Errorf("failed to scan fridge item: %w", err)
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
	}

	return items, nil
}
