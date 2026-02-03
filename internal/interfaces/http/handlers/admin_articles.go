package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	appArticle "github.com/rootix/portfolio/internal/application/article"
	"github.com/rootix/portfolio/internal/domain/article"
)

type AdminArticleHandler struct {
	List      appArticle.ListUseCase
	Get       appArticle.GetUseCase
	Create    appArticle.CreateUseCase
	Update    appArticle.UpdateUseCase
	Publish   appArticle.PublishUseCase
	Unpublish appArticle.UnpublishUseCase
	Delete    appArticle.DeleteUseCase
}

type articlePayload struct {
	Title    string `json:"title"`
	Markdown string `json:"markdown"`
}

type adminArticleResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Markdown  string    `json:"markdown"`
	HTML      string    `json:"html"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (h AdminArticleHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	limit := parseIntQuery(r, "limit", 50)
	offset := parseIntQuery(r, "offset", 0)

	items, err := h.List.All(r.Context(), appArticle.ListAllInput{Limit: limit, Offset: offset})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list articles")
		return
	}

	response := make([]adminArticleResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toAdminArticleResponse(item))
	}

	writeJSON(w, http.StatusOK, response)
}

func (h AdminArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var payload articlePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	item, err := h.Create.Execute(r.Context(), appArticle.CreateInput{Title: payload.Title, Markdown: payload.Markdown})
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to create article")
		return
	}

	writeJSON(w, http.StatusCreated, toAdminArticleResponse(item))
}

func (h AdminArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var payload articlePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	item, err := h.Update.Execute(r.Context(), appArticle.UpdateInput{ID: id, Title: payload.Title, Markdown: payload.Markdown})
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to update article")
		return
	}

	writeJSON(w, http.StatusOK, toAdminArticleResponse(item))
}

func (h AdminArticleHandler) PublishArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	item, err := h.Publish.Execute(r.Context(), appArticle.PublishInput{ID: id})
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to publish article")
		return
	}

	writeJSON(w, http.StatusOK, toAdminArticleResponse(item))
}

func (h AdminArticleHandler) UnpublishArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	item, err := h.Unpublish.Execute(r.Context(), appArticle.UnpublishInput{ID: id})
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to unpublish article")
		return
	}

	writeJSON(w, http.StatusOK, toAdminArticleResponse(item))
}

func (h AdminArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.Delete.Execute(r.Context(), appArticle.DeleteInput{ID: id}); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete article")
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func toAdminArticleResponse(item *article.Article) adminArticleResponse {
	return adminArticleResponse{
		ID:        item.ID,
		Title:     item.Title,
		Slug:      item.Slug,
		Markdown:  item.MarkdownContent,
		HTML:      item.HTMLContent,
		Status:    string(item.Status),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
