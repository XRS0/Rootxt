package article

import (
	"context"

	app "github.com/rootix/portfolio/internal/application"
	"github.com/rootix/portfolio/internal/domain/article"
)

type UnpublishInput struct {
	ID int64
}

type UnpublishUseCase struct {
	Repo  article.Repository
	Clock app.Clock
}

func (uc UnpublishUseCase) Execute(ctx context.Context, input UnpublishInput) (*article.Article, error) {
	if uc.Clock == nil {
		uc.Clock = app.SystemClock{}
	}

	entity, err := uc.Repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	entity.MarkDraft(uc.Clock.Now())

	if err := uc.Repo.Update(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}
