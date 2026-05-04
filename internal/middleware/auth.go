package middleware

import (
	"CS367-Finance-Management-System/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		userIDFloat, _ := claims["user_id"].(float64)
		c.Set("user_id", int(userIDFloat))
		c.Set("email", claims["email"])
		c.Set("role", claims["role"])

		c.Next()
	}
}
