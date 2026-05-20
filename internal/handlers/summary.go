package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BalanceHandler struct {
	DB *sql.DB
}

func NewBalanceHandler(db *sql.DB) *BalanceHandler {
	return &BalanceHandler{DB: db}
}

// GET /api/summary/balance
// คืนยอดคงเหลือสุทธิ = รายรับรวม - รายจ่ายรวม ของ user ที่ login อยู่
func (h *BalanceHandler) GetBalance(c *gin.Context) {
	userID := c.GetInt("user_id")

	var totalIncome, totalExpense float64

	_ = h.DB.QueryRow(
		"SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE user_id = ? AND type = 'income'",
		userID,
	).Scan(&totalIncome)

	_ = h.DB.QueryRow(
		"SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE user_id = ? AND type = 'expense'",
		userID,
	).Scan(&totalExpense)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"balance": totalIncome - totalExpense,
		"income":  totalIncome,
		"expense": totalExpense,
	})
}
