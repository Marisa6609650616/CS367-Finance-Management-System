package middleware

import (
	"CS367-Finance-Management-System/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware ตรวจสอบ JWT Token จาก Header: Authorization: Bearer <token>
// เสฎฐวุฒิจะเป็นคนออก Token จาก POST /api/auth/login
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Access denied. No token provided.",
			})
			return
		}

		// ต้องขึ้นต้นด้วย "Bearer "
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token format. Use: Bearer <token>",
			})
			return
		}

		tokenStr := parts[1]
		secret := config.GetEnv("JWT_SECRET")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid or expired token.",
			})
			return
		}

		// ดึง user_id จาก claims แล้วเก็บไว้ใน context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token claims.",
			})
			return
		}

		userID := int(claims["user_id"].(float64))
		c.Set("userID", userID)
		c.Next()
	}
}
