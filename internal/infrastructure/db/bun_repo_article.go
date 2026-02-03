package db

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/rootix/portfolio/internal/domain/article"
)

type ArticleRepository struct {
	DB *bun.DB
}

func (r ArticleRepository) Create(ctx context.Context, entity *article.Article) error {
	model := ArticleModelFromDomain(entity)
	if _, err := r.DB.NewInsert().Model(model).Exec(ctx); err != nil {
		return err
	}
	entity.ID = model.ID
	return nil
}

func (r ArticleRepository) Update(ctx context.Context, entity *article.Article) error {
	model := ArticleModelFromDomain(entity)
	_, err := r.DB.NewUpdate().Model(model).
		Where("id = ?", entity.ID).
		Exec(ctx)
	return err
}

func (r ArticleRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.DB.NewDelete().Model((*ArticleModel)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (r ArticleRepository) GetByID(ctx context.Context, id int64) (*article.Article, error) {
	model := new(ArticleModel)
	if err := r.DB.NewSelect().Model(model).Where("id = ?", id).Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return nil, article.ErrNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r ArticleRepository) GetBySlug(ctx context.Context, slug string) (*article.Article, error) {
	model := new(ArticleModel)
	if err := r.DB.NewSelect().Model(model).Where("slug = ?", slug).Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			return nil, article.ErrNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r ArticleRepository) ListPublished(ctx context.Context, limit, offset int) ([]*article.Article, error) {
	models := make([]*ArticleModel, 0)
	query := r.DB.NewSelect().Model(&models).Where("status = ?", string(article.StatusPublished)).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return mapArticleModels(models), nil
}

func (r ArticleRepository) ListAll(ctx context.Context, limit, offset int) ([]*article.Article, error) {
	models := make([]*ArticleModel, 0)
	query := r.DB.NewSelect().Model(&models).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return mapArticleModels(models), nil
}

func mapArticleModels(models []*ArticleModel) []*article.Article {
	items := make([]*article.Article, 0, len(models))
	for _, model := range models {
		items = append(items, model.ToDomain())
	}
	return items
}
