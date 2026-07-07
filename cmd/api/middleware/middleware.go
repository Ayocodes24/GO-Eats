package middleware

import (
	"net/http"
	"strings"

	"github.com/Ayocodes24/GO-Eats/pkg/auth"
	"github.com/gin-gonic/gin"
)

// UserClaims re-exported so callers that import this package still compile.
type UserClaims = auth.UserClaims

func ValidateToken(token string) (bool, int64) {
	return auth.ValidateToken(token)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header!!"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization"})
			c.Abort()
			return
		}
		token := tokenParts[1]

		ok, userID := auth.ValidateToken(token)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
