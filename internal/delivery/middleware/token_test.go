package middleware_test

import (
	"testing"
	"time"

	"pvz/internal/delivery/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	userID := "12345"
	role := "employee"

	tokenStr, err := middleware.GenerateToken(userID, role)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	require.Equal(t, userID, claims["user_id"])
	require.Equal(t, role, claims["role"])

	exp := int64(claims["exp"].(float64))
	require.True(t, time.Unix(exp, 0).After(time.Now()))
}
