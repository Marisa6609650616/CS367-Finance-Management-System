package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// SummaryResponse สำหรับส่งค่ากลับไปหาผู้ใช้งาน
type SummaryResponse struct {
	Status  string  `json:"status"`
	Balance float64 `json:"balance"`
	Message string  `json:"message"`
}

// GetBalance ฟังก์ชันของมาริษาสำหรับคำนวณยอดคงเหลือสุทธิ
func GetBalance(db *sql.DB, userID int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var totalIncome, totalExpense float64

		// ดึงยอดรายรับรวม
		err := db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE user_id = ? AND type = 'income'", userID).Scan(&totalIncome)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// ดึงยอดรายจ่ายรวม
		err = db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE user_id = ? AND type = 'expense'", userID).Scan(&totalExpense)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		balance := totalIncome - totalExpense

		res := SummaryResponse{
			Status:  "success",
			Balance: balance,
			Message: "Total balance calculated successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
