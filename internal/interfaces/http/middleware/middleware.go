package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"todo_list/config"
	"todo_list/internal/interfaces/http/repository"
)

func TokenMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		slog.Info("request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)

		if r.URL.Path == "/users/login" || r.URL.Path == "/users/signup" {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		userID, err := isValidToken(token, cfg)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isValidToken(tokenStr string, cfg *config.Config) (uint, error) {
	claims, err := repository.ExtractClaims(tokenStr, cfg)
	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}
	return claims.UserID, nil
}
