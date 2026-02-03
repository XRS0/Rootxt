package article

import (
	"context"
	"errors"
	"fmt"

	app "github.com/rootix/portfolio/internal/application"
	"github.com/rootix/portfolio/internal/domain/article"
)

type CreateInput struct {
	Title    string
	Markdown string
}

type CreateUseCase struct {
	Repo  article.Repository
	Clock app.Clock
}

func (uc CreateUseCase) Execute(ctx context.Context, input CreateInput) (*article.Article, error) {
	if uc.Clock == nil {
		uc.Clock = app.SystemClock{}
	}

	entity, err := article.NewDraft(input.Title, input.Markdown, uc.Clock.Now())
	if err != nil {
		return nil, err
	}

	uniqueSlug, err := ensureUniqueSlug(ctx, uc.Repo, entity.Slug, 0)
	if err != nil {
		return nil, err
	}
	entity.Slug = uniqueSlug

	if err := uc.Repo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func ensureUniqueSlug(ctx context.Context, repo article.Repository, base string, start int) (string, error) {
	candidate := base
	index := start
	for {
		existing, err := repo.GetBySlug(ctx, candidate)
		if err != nil {
			if errors.Is(err, article.ErrNotFound) {
				return candidate, nil
			}
			return "", err
		}
		if existing == nil {
			return candidate, nil
		}
		index++
		candidate = fmt.Sprintf("%s-%d", base, index)
	}
}
