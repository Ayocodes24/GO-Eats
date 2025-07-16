package user

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateJWT(userID int64) (string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24hr expiry
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
