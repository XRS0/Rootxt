package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	appArticle "github.com/rootix/portfolio/internal/application/article"
	appUser "github.com/rootix/portfolio/internal/application/user"
	"github.com/rootix/portfolio/internal/infrastructure/auth"
	"github.com/rootix/portfolio/internal/infrastructure/db"
	"github.com/rootix/portfolio/internal/infrastructure/http"
	"github.com/rootix/portfolio/internal/infrastructure/markdown"
	"github.com/rootix/portfolio/internal/interfaces/http/handlers"
)

func main() {
	ctx := context.Background()

	database, err := db.Open(ctx)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	articleRepo := db.ArticleRepository{DB: database}
	userRepo := db.UserRepository{DB: database}

	hasher := auth.BcryptHasher{}
	if adminEmail := os.Getenv("ADMIN_EMAIL"); adminEmail != "" {
		if adminPassword := os.Getenv("ADMIN_PASSWORD"); adminPassword != "" {
			_, err := appUser.BootstrapAdminUseCase{
				Repo:   userRepo,
				Hasher: hasher,
			}.Execute(ctx, appUser.BootstrapAdminInput{Email: adminEmail, Password: adminPassword})
			if err != nil {
				log.Fatalf("failed to bootstrap admin user: %v", err)
			}
			log.Printf("admin user ensured: %s", adminEmail)
		}
	}

	renderer := markdown.NewRenderer()

	jwtSecret := getenv("JWT_SECRET", "change-me")
	jwtTTL := getenvInt("JWT_TTL_MINUTES", 120)
	jwtManager := auth.JWTManager{
		Secret: []byte(jwtSecret),
		TTL:    time.Duration(jwtTTL) * time.Minute,
	}

	publicArticles := handlers.PublicArticleHandler{
		List: appArticle.ListUseCase{Repo: articleRepo},
		Get:  appArticle.GetUseCase{Repo: articleRepo},
	}

	adminArticles := handlers.AdminArticleHandler{
		List:      appArticle.ListUseCase{Repo: articleRepo},
		Get:       appArticle.GetUseCase{Repo: articleRepo},
		Create:    appArticle.CreateUseCase{Repo: articleRepo},
		Update:    appArticle.UpdateUseCase{Repo: articleRepo},
		Publish:   appArticle.PublishUseCase{Repo: articleRepo, Renderer: renderer},
		Unpublish: appArticle.UnpublishUseCase{Repo: articleRepo},
		Delete:    appArticle.DeleteUseCase{Repo: articleRepo},
	}

	adminAuth := handlers.AdminAuthHandler{
		LoginUseCase: appUser.LoginUseCase{Repo: userRepo, Hasher: hasher},
		JWT:          jwtManager,
	}

	router := httpserver.NewRouter(httpserver.RouterParams{
		PublicArticles: publicArticles,
		AdminArticles:  adminArticles,
		AdminAuth:      adminAuth,
		JWT:            jwtManager,
		CORSOrigin:     os.Getenv("CORS_ORIGIN"),
	})

	addr := ":" + getenv("PORT", "8080")
	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server stopped: %v", err)
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getenvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
