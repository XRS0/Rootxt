package httpserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rootix/portfolio/internal/infrastructure/auth"
	"github.com/rootix/portfolio/internal/interfaces/http/handlers"
	"github.com/rootix/portfolio/internal/interfaces/http/middleware"
)

type RouterParams struct {
	PublicArticles handlers.PublicArticleHandler
	AdminArticles  handlers.AdminArticleHandler
	AdminAuth      handlers.AdminAuthHandler
	JWT            auth.JWTManager
	CORSOrigin     string
}

func NewRouter(params RouterParams) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.JSON)
	if params.CORSOrigin != "" {
		router.Use(middleware.CORS(params.CORSOrigin))
	}
	router.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	api := router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/articles", params.PublicArticles.ListPublished).Methods(http.MethodGet)
	api.HandleFunc("/articles/{slug}", params.PublicArticles.GetBySlug).Methods(http.MethodGet)

	adminPublic := api.PathPrefix("/admin").Subrouter()
	adminPublic.HandleFunc("/login", params.AdminAuth.Login).Methods(http.MethodPost)

	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.RequireAuth(params.JWT))
	admin.HandleFunc("/articles", params.AdminArticles.ListAll).Methods(http.MethodGet)
	admin.HandleFunc("/articles", params.AdminArticles.CreateArticle).Methods(http.MethodPost)
	admin.HandleFunc("/articles/{id}", params.AdminArticles.UpdateArticle).Methods(http.MethodPut)
	admin.HandleFunc("/articles/{id}", params.AdminArticles.DeleteArticle).Methods(http.MethodDelete)
	admin.HandleFunc("/articles/{id}/publish", params.AdminArticles.PublishArticle).Methods(http.MethodPost)
	admin.HandleFunc("/articles/{id}/unpublish", params.AdminArticles.UnpublishArticle).Methods(http.MethodPost)

	return router
}
