package utils

import (
	"CS367-Finance-Management-System/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int, email string, role string) (string, error) {
	secret := config.GetEnv("JWT_SECRET")

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}