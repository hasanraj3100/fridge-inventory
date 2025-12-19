// Package service
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/hasanraj3100/fridge-inventory/internal/domain"
	"github.com/hasanraj3100/fridge-inventory/internal/repository"
	"github.com/hasanraj3100/fridge-inventory/internal/utils"
)

var (
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
)

type UserService interface {
	Register(ctx context.Context, username, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type userService struct {
	userRepo        repository.UserRepository
	passwordManager *utils.PasswordManager
	jwtManager      *utils.JWTManager
}

func NewUserService(userRepo repository.UserRepository, passwordManager *utils.PasswordManager, jwtManager *utils.JWTManager) UserService {
	return &userService{
		userRepo:        userRepo,
		passwordManager: passwordManager,
		jwtManager:      jwtManager,
	}
}

func (s *userService) Register(ctx context.Context, username, email, password string) (*domain.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("username, email, and password cannot be empty")
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := s.passwordManager.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		UserName:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}
	newUser.PasswordHash = ""

	return newUser, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	if err := s.passwordManager.VerifyPassword(user.PasswordHash, password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(&utils.PayLoad{
		Sub:      user.ID,
		Username: user.UserName,
		Email:    user.Email,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return token, nil
}
