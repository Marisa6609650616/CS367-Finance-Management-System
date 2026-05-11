// middleware/auth.go
// Mock auth middleware — ใช้ชั่วคราวจนกว่าเสฎฐวุฒิจะทำเสร็จ
package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Mock user — สมมติว่า login อยู่ด้วย user id = 1
		ctx := context.WithValue(r.Context(), UserIDKey, 1)
		next(w, r.WithContext(ctx))
	}
}
