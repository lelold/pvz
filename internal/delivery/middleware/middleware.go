package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"message":"missing or invalid token"}`, http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"message":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		userID, ok1 := claims["user_id"].(string)
		role, ok2 := claims["role"].(string)
		if !ok1 || !ok2 {
			http.Error(w, `{"message":"invalid token payload"}`, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, roleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(required string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, err := GetUserRole(r.Context())
		if err != nil || role != required {
			http.Error(w, `{"message":"forbidden"}`, http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func GetUserID(ctx context.Context) (string, error) {
	val := ctx.Value(userIDKey)
	userID, ok := val.(string)
	if !ok {
		return "", errors.New("user_id not found")
	}
	return userID, nil
}

func GetUserRole(ctx context.Context) (string, error) {
	val := ctx.Value(roleKey)
	role, ok := val.(string)
	if !ok {
		return "", errors.New("role not found")
	}
	return role, nil
}
