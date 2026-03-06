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
