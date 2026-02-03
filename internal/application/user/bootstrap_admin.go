package user

import (
	"context"
	"errors"

	app "github.com/rootix/portfolio/internal/application"
	"github.com/rootix/portfolio/internal/domain/user"
)

type BootstrapAdminInput struct {
	Email    string
	Password string
}

type BootstrapAdminUseCase struct {
	Repo   user.Repository
	Hasher PasswordHasher
	Clock  app.Clock
}

func (uc BootstrapAdminUseCase) Execute(ctx context.Context, input BootstrapAdminInput) (*user.User, error) {
	if uc.Clock == nil {
		uc.Clock = app.SystemClock{}
	}
	if input.Email == "" || input.Password == "" {
		return nil, user.ErrInvalidEmail
	}

	email := user.NormalizeEmail(input.Email)
	hash, err := uc.Hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	entity, err := uc.Repo.GetByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, user.ErrNotFound) {
			return nil, err
		}
		newUser, err := user.New(email, hash, uc.Clock.Now())
		if err != nil {
			return nil, err
		}
		if err := uc.Repo.Create(ctx, newUser); err != nil {
			return nil, err
		}
		return newUser, nil
	}

	entity.PasswordHash = hash
	entity.Touch(uc.Clock.Now())
	if err := uc.Repo.Update(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}
