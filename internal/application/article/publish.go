package article

import (
	"context"

	app "github.com/rootix/portfolio/internal/application"
	"github.com/rootix/portfolio/internal/domain/article"
)

type MarkdownRenderer interface {
	Render(markdown string) (string, error)
}

type PublishInput struct {
	ID int64
}

type PublishUseCase struct {
	Repo     article.Repository
	Clock    app.Clock
	Renderer MarkdownRenderer
}

func (uc PublishUseCase) Execute(ctx context.Context, input PublishInput) (*article.Article, error) {
	if uc.Clock == nil {
		uc.Clock = app.SystemClock{}
	}

	entity, err := uc.Repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	html, err := uc.Renderer.Render(entity.MarkdownContent)
	if err != nil {
		return nil, err
	}

	if err := entity.Publish(html, uc.Clock.Now()); err != nil {
		return nil, err
	}

	if err := uc.Repo.Update(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}
