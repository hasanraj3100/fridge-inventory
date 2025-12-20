package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/hasanraj3100/fridge-inventory/internal/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"

type AuthConfig struct {
	JWTManager *utils.JWTManager
	Skipper    func(r *http.Request) bool
}

func Auth(jwtManager *utils.JWTManager) func(http.Handler) http.Handler {
	return AuthWithConfig(AuthConfig{
		JWTManager: jwtManager,
		Skipper:    nil,
	})
}

func AuthWithConfig(config AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.Skipper != nil && config.Skipper(r) {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}
			token := parts[1]
			userID, err := config.JWTManager.VerifyToken(token)
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}
