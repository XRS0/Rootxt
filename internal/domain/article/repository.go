package article

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("article not found")

// Repository defines persistence operations for articles.
type Repository interface {
	Create(ctx context.Context, article *Article) error
	Update(ctx context.Context, article *Article) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*Article, error)
	GetBySlug(ctx context.Context, slug string) (*Article, error)
	ListPublished(ctx context.Context, limit, offset int) ([]*Article, error)
	ListAll(ctx context.Context, limit, offset int) ([]*Article, error)
}
