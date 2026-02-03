package user

import (
	"context"
	"errors"

	"github.com/rootix/portfolio/internal/domain/user"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginUseCase struct {
	Repo   user.Repository
	Hasher PasswordHasher
}

func (uc LoginUseCase) Execute(ctx context.Context, input LoginInput) (*user.User, error) {
	if input.Email == "" || input.Password == "" {
		return nil, ErrInvalidCredentials
	}

	email := user.NormalizeEmail(input.Email)
	entity, err := uc.Repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := uc.Hasher.Compare(entity.PasswordHash, input.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return entity, nil
}
