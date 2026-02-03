package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	appArticle "github.com/rootix/portfolio/internal/application/article"
	"github.com/rootix/portfolio/internal/domain/article"
)

type PublicArticleHandler struct {
	List appArticle.ListUseCase
	Get  appArticle.GetUseCase
}

type publicArticleListItem struct {
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Excerpt   string    `json:"excerpt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type publicArticleResponse struct {
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	HTML      string    `json:"html"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (h PublicArticleHandler) ListPublished(w http.ResponseWriter, r *http.Request) {
	limit := parseIntQuery(r, "limit", 20)
	offset := parseIntQuery(r, "offset", 0)

	items, err := h.List.Published(r.Context(), appArticle.ListPublishedInput{Limit: limit, Offset: offset})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list articles")
		return
	}

	response := make([]publicArticleListItem, 0, len(items))
	for _, item := range items {
		response = append(response, publicArticleListItem{
			Title:     item.Title,
			Slug:      item.Slug,
			Excerpt:   excerpt(item.MarkdownContent, 180),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	writeJSON(w, http.StatusOK, response)
}

func (h PublicArticleHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	if strings.TrimSpace(slug) == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	item, err := h.Get.BySlug(r.Context(), appArticle.GetBySlugInput{Slug: slug})
	if err != nil {
		writeError(w, http.StatusNotFound, "article not found")
		return
	}
	if item.Status != article.StatusPublished {
		writeError(w, http.StatusNotFound, "article not found")
		return
	}

	writeJSON(w, http.StatusOK, publicArticleResponse{
		Title:     item.Title,
		Slug:      item.Slug,
		HTML:      item.HTMLContent,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	})
}

func parseIntQuery(r *http.Request, key string, fallback int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return fallback
	}
	return parsed
}

func excerpt(markdown string, max int) string {
	text := strings.TrimSpace(markdown)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	text = strings.Join(strings.Fields(text), " ")
	if len(text) <= max {
		return text
	}
	cut := text[:max]
	if idx := strings.LastIndex(cut, " "); idx > 0 {
		cut = cut[:idx]
	}
	return cut + "..."
}
