package handlers

import (
	"CS367-Finance-Management-System/middleware"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type SummaryMonthyHandler struct {
	DB *sql.DB
}

func NewSumaryMonthyHandler(db *sql.DB) *SummaryMonthyHandler {
	return &SummaryMonthyHandler{DB: db}
}

type SummaryMonthlyResponse struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

func (h *SummaryMonthyHandler) SummaryMonthly(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}

	// รับ query param ?month=2025-04 (ถ้าไม่ส่งมาใช้เดือนปัจจุบัน)
	month := r.URL.Query().Get("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	// Query รวม income และ expense ของเดือนนั้น
	rows, err := h.DB.Query(`
		SELECT type, SUM(amount)
		FROM transactions
		WHERE user_id = ?
		  AND strftime('%Y-%m', date) = ?
		GROUP BY type
	`, userID, month)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "เกิดข้อผิดพลาดในการดึงข้อมูล"})
		return
	}
	defer rows.Close()

	var income, expense float64
	for rows.Next() {
		var txType string
		var total float64
		if err := rows.Scan(&txType, &total); err != nil {
			continue
		}
		switch txType {
		case "income":
			income = total
		case "expense":
			expense = total
		}
	}

	summary := SummaryMonthlyResponse{
		Month:   month,
		Income:  income,
		Expense: expense,
		Balance: income - expense,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "ดึงข้อมูลสำเร็จ",
		"summary": summary,
	})
}
