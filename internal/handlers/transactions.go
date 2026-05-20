// internal/handlers/transactions.go
package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	// รับค่า user_id จาก middleware
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userID := userIDVal.(int)

	// Parse request body
	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "รูปแบบข้อมูลไม่ถูกต้อง"})
		return
	}

	// Validate
	if req.CategoryID == 0 || req.Type == "" || req.Amount <= 0 || req.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "กรุณากรอกข้อมูลให้ครบ (category_id, type, amount, date)"})
		return
	}
	if req.Type != "income" && req.Type != "expense" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "type ต้องเป็น income หรือ expense เท่านั้น"})
		return
	}

	// เช็คว่า category มีอยู่จริง
	var count int
	err := h.DB.QueryRow("SELECT COUNT(*) FROM categories WHERE id = ?", req.CategoryID).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "ไม่พบหมวดหมู่นี้"})
		return
	}

	// บันทึกลง DB
	result, err := h.DB.Exec(
		"INSERT INTO transactions (user_id, category_id, type, amount, note, date) VALUES (?, ?, ?, ?, ?, ?)",
		userID, req.CategoryID, req.Type, req.Amount, req.Note, req.Date,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "เกิดข้อผิดพลาดในการบันทึกข้อมูล"})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{
		"message": "เพิ่มรายการสำเร็จ",
		"transaction": gin.H{
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
