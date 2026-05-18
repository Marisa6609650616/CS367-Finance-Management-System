package middleware

import (
	"CS367-Finance-Management-System/pkg/utils"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ดึง token จาก Header: Authorization: Bearer <token>
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"message":"Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			http.Error(w, `{"message":"Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// ดึง user_id จาก claims แล้วใส่ใน context
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, `{"message":"Invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, int(userIDFloat))
		next(w, r.WithContext(ctx))
	}
}
