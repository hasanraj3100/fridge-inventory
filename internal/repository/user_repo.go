// Package repository acts as a data access layer
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hasanraj3100/fridge-inventory/internal/db"
	"github.com/hasanraj3100/fridge-inventory/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUserName(ctx context.Context, username string) (*domain.User, error)
}

type userRepository struct {
	DB db.DBTX
}

func NewUserRepository(db db.DBTX) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (userRepo *userRepository) Create(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now().UTC()

	query := `INSERT INTO users (username, email, password_hash, created_at)
	VALUES ($1, $2, $3, $4) RETURNING id`

	err := userRepo.DB.QueryRowContext(ctx, query,
		user.UserName,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to insert user to database: %w", err)
	}

	return nil
}

func (userRepo *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := `SELECT * FROM users WHERE email = $1`

	row := userRepo.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&user.ID, &user.UserName, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (userRepo *userRepository) GetByUserName(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User

	query := `SELECT * FROM users WHERE username = $1`

	row := userRepo.DB.QueryRowContext(ctx, query, username)
	err := row.Scan(&user.ID, &user.UserName, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}
