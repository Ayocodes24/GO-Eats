package auth

import (
	"log/slog"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

func ValidateToken(token string) (bool, int64) {
	tokenInfo, err := jwt.ParseWithClaims(token, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		slog.Error("auth.ValidateToken", "error", err.Error())
		return false, 0
	}
	if claims, ok := tokenInfo.Claims.(*UserClaims); ok {
		return true, claims.UserID
	}
	return false, 0
}
