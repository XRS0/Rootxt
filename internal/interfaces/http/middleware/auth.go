package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rootix/portfolio/internal/infrastructure/auth"
)

type contextKey string

const userIDKey contextKey = "userID"

func UserIDFromContext(ctx context.Context) (int64, bool) {
	value := ctx.Value(userIDKey)
	if value == nil {
		return 0, false
	}
	id, ok := value.(int64)
	return id, ok
}

func RequireAuth(manager auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(authorization, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			userID, err := manager.Verify(parts[1])
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
