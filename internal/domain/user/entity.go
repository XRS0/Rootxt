package user

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidEmail    = errors.New("email is required")
	ErrInvalidPassword = errors.New("password is required")
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func New(email, passwordHash string, now time.Time) (*User, error) {
	cleanEmail := strings.TrimSpace(email)
	if cleanEmail == "" {
		return nil, ErrInvalidEmail
	}
	cleanHash := strings.TrimSpace(passwordHash)
	if cleanHash == "" {
		return nil, ErrInvalidPassword
	}

	return &User{
		Email:        strings.ToLower(cleanEmail),
		PasswordHash: cleanHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (u *User) Touch(now time.Time) {
	u.UpdatedAt = now
}
