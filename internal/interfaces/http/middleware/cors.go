package middleware

import (
	"net/http"
	"strings"
)

func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	allowed := parseOrigins(allowedOrigin)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				origin = "*"
			}
			headerOrigin := resolveOrigin(origin, allowed)
			if headerOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", headerOrigin)
				if headerOrigin != "*" {
					w.Header().Set("Vary", "Origin")
				}
			}
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseOrigins(raw string) []string {
	if raw == "" {
		return []string{"*"}
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		origins = append(origins, value)
	}
	if len(origins) == 0 {
		return []string{"*"}
	}
	return origins
}

func resolveOrigin(origin string, allowed []string) string {
	if len(allowed) == 0 {
		return ""
	}
	for _, entry := range allowed {
		if entry == "*" {
			return "*"
		}
		if origin == entry {
			return origin
		}
	}
	return ""
}
