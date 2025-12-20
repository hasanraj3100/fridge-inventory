// Package domain defines core domain models for the application.
package domain

import "time"

type User struct {
	ID           int       `json:"id" db:"id"`
	UserName     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
