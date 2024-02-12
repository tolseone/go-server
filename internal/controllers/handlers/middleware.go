package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-server/internal/models"
	"go-server/pkg/logging"
)

func AuthMiddleware(next http.Handler, logger *logging.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := model.ParseToken(tokenString)
		if err != nil {
			logger.Errorf("Invalid or expired token: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
