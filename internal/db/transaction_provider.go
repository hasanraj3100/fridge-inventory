package db

import (
	"context"
	"database/sql"
)

type TransactionProvider interface {
	WithinTransaction(ctx context.Context, fn func(tx DBTX) error) error
}

type SQLTransactionProvider struct {
	DB *sql.DB
}

func NewSQLTransactionProvider(db *sql.DB) TransactionProvider {
	return &SQLTransactionProvider{
		DB: db,
	}
}

func (p *SQLTransactionProvider) WithinTransaction(ctx context.Context, fn func(tx DBTX) error) error {
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
