// handlers/transactions.go
package handlers

import (
	"CS367-Finance-Management-System/middleware"
	"database/sql"
	"encoding/json"
	"net/http"
)

type TransactionRequest struct {
	CategoryID int     `json:"category_id"`
	Type       string  `json:"type"`
	Amount     float64 `json:"amount"`
	Note       string  `json:"note"`
	Date       string  `json:"date"`
}

type TransactionHandler struct {
	DB *sql.DB
}

func NewTransactionHandler(db *sql.DB) *TransactionHandler {
	return &TransactionHandler{DB: db}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	// รับค่า user_id จาก middleware
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}

	// Parse request body
	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "รูปแบบข้อมูลไม่ถูกต้อง"})
		return
	}

	// Validate
	if req.CategoryID == 0 || req.Type == "" || req.Amount <= 0 || req.Date == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "กรุณากรอกข้อมูลให้ครบ (category_id, type, amount, date)"})
		return
	}
	if req.Type != "income" && req.Type != "expense" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "type ต้องเป็น income หรือ expense เท่านั้น"})
		return
	}

	// เช็คว่า category มีอยู่จริง
	var count int
	err := h.DB.QueryRow("SELECT COUNT(*) FROM categories WHERE id = ?", req.CategoryID).Scan(&count)
	if err != nil || count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "ไม่พบหมวดหมู่นี้"})
		return
	}

	// บันทึกลง DB
	result, err := h.DB.Exec(
		"INSERT INTO transactions (user_id, category_id, type, amount, note, date) VALUES (?, ?, ?, ?, ?, ?)",
		userID, req.CategoryID, req.Type, req.Amount, req.Note, req.Date,
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "เกิดข้อผิดพลาดในการบันทึกข้อมูล"})
		return
	}

	id, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "เพิ่มรายการสำเร็จ",
		"transaction": map[string]interface{}{
			"id":          id,
			"user_id":     userID,
			"category_id": req.CategoryID,
			"type":        req.Type,
			"amount":      req.Amount,
			"note":        req.Note,
			"date":        req.Date,
		},
	})
}
