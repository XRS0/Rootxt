package article

import (
	"context"

	"github.com/rootix/portfolio/internal/domain/article"
)

type ListPublishedInput struct {
	Limit  int
	Offset int
}

type ListAllInput struct {
	Limit  int
	Offset int
}

type ListUseCase struct {
	Repo article.Repository
}

func (uc ListUseCase) Published(ctx context.Context, input ListPublishedInput) ([]*article.Article, error) {
	return uc.Repo.ListPublished(ctx, input.Limit, input.Offset)
}

func (uc ListUseCase) All(ctx context.Context, input ListAllInput) ([]*article.Article, error) {
	return uc.Repo.ListAll(ctx, input.Limit, input.Offset)
}
