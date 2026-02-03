package article

import (
	"context"
	"errors"
	"fmt"

	app "github.com/rootix/portfolio/internal/application"
	"github.com/rootix/portfolio/internal/domain/article"
)

type UpdateInput struct {
	ID       int64
	Title    string
	Markdown string
}

type UpdateUseCase struct {
	Repo  article.Repository
	Clock app.Clock
}

func (uc UpdateUseCase) Execute(ctx context.Context, input UpdateInput) (*article.Article, error) {
	if uc.Clock == nil {
		uc.Clock = app.SystemClock{}
	}

	entity, err := uc.Repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if err := entity.Update(input.Title, input.Markdown, uc.Clock.Now()); err != nil {
		return nil, err
	}

	if entity.Status == article.StatusPublished {
		entity.MarkDraft(uc.Clock.Now())
	}

	uniqueSlug, err := ensureUniqueSlugForUpdate(ctx, uc.Repo, entity.Slug, entity.ID)
	if err != nil {
		return nil, err
	}
	entity.Slug = uniqueSlug

	if err := uc.Repo.Update(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func ensureUniqueSlugForUpdate(ctx context.Context, repo article.Repository, base string, id int64) (string, error) {
	candidate := base
	index := 0
	for {
		existing, err := repo.GetBySlug(ctx, candidate)
		if err != nil {
			if errors.Is(err, article.ErrNotFound) {
				return candidate, nil
			}
			return "", err
		}
		if existing == nil || existing.ID == id {
			return candidate, nil
		}
		index++
		candidate = fmt.Sprintf("%s-%d", base, index)
	}
}
