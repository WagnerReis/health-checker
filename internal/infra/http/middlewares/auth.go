package middlewares

import (
	"context"
	"health-checker/config"
	"health-checker/internal/infra/http/helpers"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := config.LoadConfig()
		token := r.Header.Get("Authorization")
		if token == "" {
			helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		tokenData, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.AccessTokenSecret), nil
		})
		if err != nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		claims, ok := tokenData.Claims.(*Claims)
		if !ok {
			helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		var UserIDKey = "userID"
		var UserEmailKey = "userEmail"

		if tokenData.Valid {
			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
	})
}
