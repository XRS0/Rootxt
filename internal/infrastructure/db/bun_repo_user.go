package db

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/rootix/portfolio/internal/domain/user"
)

type UserRepository struct {
	DB *bun.DB
}

func (r UserRepository) Create(ctx context.Context, entity *user.User) error {
	model := UserModelFromDomain(entity)
	if _, err := r.DB.NewInsert().Model(model).Exec(ctx); err != nil {
		return err
	}
	entity.ID = model.ID
	return nil
}

func (r UserRepository) Update(ctx context.Context, entity *user.User) error {
	model := UserModelFromDomain(entity)
	_, err := r.DB.NewUpdate().Model(model).Where("id = ?", entity.ID).Exec(ctx)
	return err
}

func (r UserRepository) GetByID(ctx context.Context, id int64) (*user.User, error) {
	model := new(UserModel)
	if err := r.DB.NewSelect().Model(model).Where("id = ?", id).Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return nil, user.ErrNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	model := new(UserModel)
	if err := r.DB.NewSelect().Model(model).Where("email = ?", email).Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return nil, user.ErrNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}
