package middleware

import (
	"context"
	"net/http"
	"strings"

	userpb "github.com/Ayocodes24/GO-Eats/pkg/proto/user"
	"github.com/gin-gonic/gin"
)

// GRPCAuthMiddleware validates JWT tokens by calling the User Service over gRPC.
// This is the cross-service authentication call — Order Service never reads
// the JWT secret directly; it delegates validation to the User Service.
func GRPCAuthMiddleware(userClient userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			c.Abort()
			return
		}

		resp, err := userClient.ValidateToken(context.Background(), &userpb.ValidateTokenRequest{
			Token: parts[1],
		})
		if err != nil || !resp.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", resp.UserId)
		c.Next()
	}
}
