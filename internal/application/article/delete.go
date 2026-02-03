package article

import (
	"context"

	"github.com/rootix/portfolio/internal/domain/article"
)

type DeleteInput struct {
	ID int64
}

type DeleteUseCase struct {
	Repo article.Repository
}

func (uc DeleteUseCase) Execute(ctx context.Context, input DeleteInput) error {
	return uc.Repo.Delete(ctx, input.ID)
}
