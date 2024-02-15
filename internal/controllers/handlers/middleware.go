package middleware

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"go-server/internal/models"
	"go-server/pkg/logging"

)

func AuthMiddleware(next httprouter.Handle, logger *logging.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

		if claims.UserRole == "admin" {
			next(w, r, params)
			return
		} else if claims.UserRole == "user" {
			if isPathForAdmin(r.URL.Path) {
				http.Error(w, "Access denied for user role", http.StatusForbidden)
				return
			}
			next(w, r, params)
			return
		} else {
			http.Error(w, "Unknown user role", http.StatusForbidden)
			return
		}
	}
}

func isPathForAdmin(path string) bool {
	adminURLs := []string{"/api/admin/users"}

	for _, url := range adminURLs {
		if path == url || strings.HasPrefix(path, url+"/") {
			return true
		}
	}

	return false
}
