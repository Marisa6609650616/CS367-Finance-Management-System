package handlers

import (
	"CS367-Finance-Management-System/middleware"
	"database/sql"
	"encoding/json"
	"net/http"
)

type SummaryResponse struct {
	Status  string  `json:"status"`
	Balance float64 `json:"balance"`
	Message string  `json:"message"`
}

// ปรับให้กลับมาเป็น http.HandlerFunc เพื่อให้เข้ากับ AuthMiddleware ของเพื่อน
func GetBalance(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ดึง userID จาก Context ที่เพื่อนทำไว้ใน middleware
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
			return
		}

		var totalIncome, totalExpense float64

		// คำนวณรายรับ
		_ = db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE user_id = ? AND type = 'income'", userID).Scan(&totalIncome)
		// คำนวณรายจ่าย
		_ = db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE user_id = ? AND type = 'expense'", userID).Scan(&totalExpense)

		balance := totalIncome - totalExpense

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SummaryResponse{
			Status:  "success",
			Balance: balance,
			Message: "Total balance calculated successfully",
		})
	}
}
