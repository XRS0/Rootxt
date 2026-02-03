package article

import (
	"context"

	"github.com/rootix/portfolio/internal/domain/article"
)

type GetByIDInput struct {
	ID int64
}

type GetBySlugInput struct {
	Slug string
}

type GetUseCase struct {
	Repo article.Repository
}

func (uc GetUseCase) ByID(ctx context.Context, input GetByIDInput) (*article.Article, error) {
	return uc.Repo.GetByID(ctx, input.ID)
}

func (uc GetUseCase) BySlug(ctx context.Context, input GetBySlugInput) (*article.Article, error) {
	return uc.Repo.GetBySlug(ctx, input.Slug)
}
