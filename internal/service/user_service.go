// Package service
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hasanraj3100/fridge-inventory/internal/api/dto"
	"github.com/hasanraj3100/fridge-inventory/internal/domain"
	"github.com/hasanraj3100/fridge-inventory/internal/repository"
	"github.com/hasanraj3100/fridge-inventory/internal/utils"
)

var (
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserService interface {
	Register(ctx context.Context, params dto.RegisterRequest) error
	Login(ctx context.Context, params dto.LoginRequest) (string, *domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userService struct {
	userRepo        repository.UserRepository
	passwordManager *utils.PasswordManager
	jwtManager      *utils.JWTManager
	txManager       repository.TransactionProvider
}

func NewUserService(
	userRepo repository.UserRepository,
	passwordManager *utils.PasswordManager,
	jwtManager *utils.JWTManager,
	txManager repository.TransactionProvider,
) UserService {
	return &userService{
		userRepo:        userRepo,
		passwordManager: passwordManager,
		jwtManager:      jwtManager,
		txManager:       txManager,
	}
}

func (s *userService) Register(ctx context.Context, params dto.RegisterRequest) error {
	existingUser, err := s.userRepo.GetByEmail(ctx, params.Email)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	hashedPassword, err := s.passwordManager.HashPassword(params.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := &domain.User{
		UserName:     params.Username,
		Email:        params.Email,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	newUser.PasswordHash = ""

	return nil
}

func (s *userService) Login(ctx context.Context, params dto.LoginRequest) (string, *domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, params.Email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	if user == nil {
		return "", nil, ErrInvalidCredentials
	}

	if err := s.passwordManager.VerifyPassword(user.PasswordHash, params.Password); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(&utils.PayLoad{
		Sub:      user.ID,
		Username: user.UserName,
		Email:    user.Email,
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	user.PasswordHash = ""

	return token, user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}
