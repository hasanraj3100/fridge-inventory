package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/hasanraj3100/fridge-inventory/internal/db"
	"github.com/hasanraj3100/fridge-inventory/internal/domain"
)

type ItemUsageRepository interface {
	Create(ctx context.Context, itemUsage *domain.ItemUsage) error
	GetByID(ctx context.Context, id int64) (*domain.ItemUsage, error)
	GetByUserIDWithPagination(ctx context.Context, userID int64, limit int, offset int) ([]map[string]interface{}, error)
	CountByUserID(ctx context.Context, userID int64) (int64, error)
	Update(ctx context.Context, itemUsage *domain.ItemUsage) error
	Delete(ctx context.Context, id int64) error
}

type itemUsageRepository struct {
	DB db.DBTX
}

func NewItemUsageRepository(db db.DBTX) ItemUsageRepository {
	return &itemUsageRepository{
		DB: db,
	}
}

func (repo *itemUsageRepository) Create(ctx context.Context, itemUsage *domain.ItemUsage) error {
	itemUsage.UsedAt = time.Now().UTC()

	query := `INSERT INTO item_usage (item_id, quantity_used, reason, used_at)
	VALUES ($1, $2, $3, $4 ) RETURNING id`

	res := repo.DB.QueryRowContext(ctx, query, itemUsage.ItemID, itemUsage.QuantityUsed, itemUsage.Reason, itemUsage.UsedAt)
	if res == nil {
		return fmt.Errorf("failed to insert item_usage to database: no result returned")
	}
	err := res.Scan(&itemUsage.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve inserted item_usage ID: %w", err)
	}

	return nil
}

func (repo *itemUsageRepository) GetByUserIDWithPagination(ctx context.Context, userID int64, limit int, offset int) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			iu.id,
			iu.item_id,
			fi.name as item_name,
			iu.quantity_used,
			fi.unit,
			iu.reason,
			iu.used_at
		FROM item_usage iu
		INNER JOIN fridge_items fi ON iu.item_id = fi.id
		WHERE fi.user_id = $1
		ORDER BY iu.used_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := repo.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get item usage: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, itemID int64
		var itemName, unit, reason string
		var quantityUsed float32
		var usedAt time.Time

		err := rows.Scan(&id, &itemID, &itemName, &quantityUsed, &unit, &reason, &usedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item usage: %w", err)
		}

		results = append(results, map[string]interface{}{
			"id":            id,
			"item_id":       itemID,
			"item_name":     itemName,
			"quantity_used": quantityUsed,
			"unit":          unit,
			"reason":        reason,
			"used_at":       usedAt,
		})
	}

	return results, nil
}

func (repo *itemUsageRepository) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM item_usage iu
		INNER JOIN fridge_items fi ON iu.item_id = fi.id
		WHERE fi.user_id = $1
	`

	var count int64
	err := repo.DB.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count item usage: %w", err)
	}

	return count, nil
}

func (repo *itemUsageRepository) GetByID(ctx context.Context, id int64) (*domain.ItemUsage, error) {
	query := `SELECT id, item_id, quantity_used, reason, used_at FROM item_usage WHERE id = $1`

	var itemUsage domain.ItemUsage
	err := repo.DB.QueryRowContext(ctx, query, id).Scan(
		&itemUsage.ID,
		&itemUsage.ItemID,
		&itemUsage.QuantityUsed,
		&itemUsage.Reason,
		&itemUsage.UsedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get item usage by ID: %w", err)
	}

	return &itemUsage, nil
}

func (repo *itemUsageRepository) Update(ctx context.Context, itemUsage *domain.ItemUsage) error {
	query := `UPDATE item_usage SET quantity_used = $1, reason = $2 WHERE id = $3`

	res, err := repo.DB.ExecContext(ctx, query, itemUsage.QuantityUsed, itemUsage.Reason, itemUsage.ID)
	if err != nil {
		return fmt.Errorf("failed to update item usage: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("item usage not found")
	}

	return nil
}

func (repo *itemUsageRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM item_usage WHERE id = $1`

	res, err := repo.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item usage: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("item usage not found")
	}

	return nil
}
