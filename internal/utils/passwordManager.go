package utils

import "golang.org/x/crypto/bcrypt"

type PasswordManager struct {
	cost int
}

func NewPasswordManager(cost int) *PasswordManager {
	if cost < bcrypt.MaxCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &PasswordManager{cost: cost}
}

func (pm *PasswordManager) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), pm.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (pm *PasswordManager) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
