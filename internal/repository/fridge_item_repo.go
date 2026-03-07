package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/hasanraj3100/fridge-inventory/internal/db"
	"github.com/hasanraj3100/fridge-inventory/internal/domain"
)

type FridgeItemRepository interface {
	Create(ctx context.Context, item *domain.FridgeItem) (*domain.FridgeItem, error)
	GetByID(ctx context.Context, id int64) (*domain.FridgeItem, error)
	GetByUserID(ctx context.Context, userID int64) ([]*domain.FridgeItem, error)
	BelongsToUser(ctx context.Context, itemID int64, userID int64) (bool, error)
	DecreaseQuantity(ctx context.Context, itemID int64, amount float64) error
	Update(ctx context.Context, item *domain.FridgeItem) error
}

type fridgeItemRepository struct {
	DB db.DBTX
}

func NewFridgeItemRepository(db db.DBTX) FridgeItemRepository {
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
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
	) RETURNING id`

	res := repo.DB.QueryRowContext(ctx, query, item.Name, item.Category, item.Quantity, item.Unit, item.UserID, item.BoughtAt, item.ExpiresAt, item.MinThreshold, item.CreatedAt, item.UpdatedAt)

	err := res.Scan(&item.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert item to database: %w", err)
	}

	return item, nil
}

func (repo *fridgeItemRepository) GetByUserID(ctx context.Context, userID int64) ([]*domain.FridgeItem, error) {
	query := `SELECT * FROM fridge_items WHERE user_id = $1`

	rows, err := repo.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fridge items by user ID: %w", err)
	}
	defer rows.Close()

	var items []*domain.FridgeItem
	for rows.Next() {
		var item domain.FridgeItem
		err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.Quantity, &item.Unit, &item.UserID, &item.BoughtAt, &item.ExpiresAt, &item.MinThreshold, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fridge item: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}

func (repo *fridgeItemRepository) Update(ctx context.Context, item *domain.FridgeItem) error {
	item.UpdatedAt = time.Now().UTC()

	cnt := 1
	query := `UPDATE fridge_items SET `

	if item.Name != "" {
		query += fmt.Sprintf("name = $%d, ", cnt)
		cnt++
	}

	if item.Category != "" {
		query += fmt.Sprintf("category = $%d, ", cnt)
		cnt++
	}

	if item.Quantity != 0 {
		query += fmt.Sprintf("quantity = $%d, ", cnt)
		cnt++
	}

	if item.Unit != "" {
		query += fmt.Sprintf("unit = $%d, ", cnt)
		cnt++
	}

	if !item.BoughtAt.IsZero() {
		query += fmt.Sprintf("bought_at = $%d, ", cnt)
		cnt++
	}

	if !item.ExpiresAt.IsZero() {
		query += fmt.Sprintf("expires_at = $%d, ", cnt)
		cnt++
	}

	if item.MinThreshold != 0 {
		query += fmt.Sprintf("min_threshold = $%d, ", cnt)
		cnt++
	}

	query += fmt.Sprintf("updated_at = $%d ", cnt)
	query += fmt.Sprintf("WHERE id = $%d", cnt+1)

	_, err := repo.DB.ExecContext(ctx, query, item.Name, item.Category, item.Quantity, item.Unit, item.BoughtAt, item.ExpiresAt, item.MinThreshold, item.UpdatedAt, item.ID)
	if err != nil {
		return fmt.Errorf("failed to update fridge item: %w", err)
	}

	return nil
}

func (repo *fridgeItemRepository) GetByID(ctx context.Context, id int64) (*domain.FridgeItem, error) {
	query := `SELECT * FROM fridge_items WHERE id = $1`

	var item domain.FridgeItem
	err := repo.DB.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.Name, &item.Category, &item.Quantity, &item.Unit, &item.UserID, &item.BoughtAt, &item.ExpiresAt, &item.MinThreshold, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get fridge item by ID: %w", err)
	}

	return &item, nil
}

func (repo *fridgeItemRepository) BelongsToUser(ctx context.Context, itemID int64, userID int64) (bool, error) {
	query := `SELECT COUNT(1) FROM fridge_items WHERE id = $1 AND user_id = $2`

	var count int
	err := repo.DB.QueryRowContext(ctx, query, itemID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check item ownership: %w", err)
	}

	return count > 0, nil
}

func (repo *fridgeItemRepository) DecreaseQuantity(ctx context.Context, itemID int64, amount float64) error {
	query := `
		UPDATE fridge_items
		SET quantity = quantity - $1
		WHERE id = $2
		  AND quantity >= $1
	`

	res, err := repo.DB.ExecContext(ctx, query, amount, itemID)
	if err != nil {
		return fmt.Errorf("failed to decrease item quantity: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("insufficient quantity or item not found")
	}

	return nil
}
