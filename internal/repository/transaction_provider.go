package repository

import (
	"context"

	"github.com/hasanraj3100/fridge-inventory/internal/db"
	"github.com/jmoiron/sqlx"
)

type TransactionProvider interface {
	WithinTransaction(ctx context.Context, fn func(tx db.DBTX) error) error
}

type SQLTransactionProvider struct {
	DB *sqlx.DB
}

func NewSQLTransactionProvider(db *sqlx.DB) TransactionProvider {
	return &SQLTransactionProvider{
		DB: db,
	}
}

func (p *SQLTransactionProvider) WithinTransaction(ctx context.Context, fn func(tx db.DBTX) error) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
