package db

import (
	"time"

	"github.com/uptrace/bun"
	"github.com/rootix/portfolio/internal/domain/article"
	"github.com/rootix/portfolio/internal/domain/user"
)

type ArticleModel struct {
	bun.BaseModel `bun:"table:articles"`
	ID            int64     `bun:",pk,autoincrement"`
	Title         string    `bun:",notnull"`
	Slug          string    `bun:",unique,notnull"`
	Markdown      string    `bun:"markdown_content,notnull"`
	HTML          string    `bun:"html_content"`
	Status        string    `bun:",notnull"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func (m *ArticleModel) ToDomain() *article.Article {
	if m == nil {
		return nil
	}
	return &article.Article{
		ID:              m.ID,
		Title:           m.Title,
		Slug:            m.Slug,
		MarkdownContent: m.Markdown,
		HTMLContent:     m.HTML,
		Status:          article.Status(m.Status),
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func ArticleModelFromDomain(entity *article.Article) *ArticleModel {
	if entity == nil {
		return nil
	}
	return &ArticleModel{
		ID:        entity.ID,
		Title:     entity.Title,
		Slug:      entity.Slug,
		Markdown:  entity.MarkdownContent,
		HTML:      entity.HTMLContent,
		Status:    string(entity.Status),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

type UserModel struct {
	bun.BaseModel `bun:"table:users"`
	ID            int64     `bun:",pk,autoincrement"`
	Email         string    `bun:",unique,notnull"`
	PasswordHash  string    `bun:"password_hash,notnull"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func (m *UserModel) ToDomain() *user.User {
	if m == nil {
		return nil
	}
	return &user.User{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func UserModelFromDomain(entity *user.User) *UserModel {
	if entity == nil {
		return nil
	}
	return &UserModel{
		ID:           entity.ID,
		Email:        entity.Email,
		PasswordHash: entity.PasswordHash,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}
