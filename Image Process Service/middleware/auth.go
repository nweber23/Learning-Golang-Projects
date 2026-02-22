package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func JWTMiddleware(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				JSONError(w, http.StatusUnauthorized, "Missing Authorization header")
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				JSONError(w, http.StatusUnauthorized, "Invalid Authorization header format")
				return
			}
			tokenString := parts[1]

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				JSONError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			r = SetUserIDInContext(r, claims.UserID)
			next.ServeHTTP(w, r)
		})
	}
}
