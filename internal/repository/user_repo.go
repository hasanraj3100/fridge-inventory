package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hasanraj3100/fridge-inventory/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
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

	err := userRepo.DB.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}
