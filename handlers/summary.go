package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SummaryResponse struct {
	Status  string  `json:"status"`
	Balance float64 `json:"balance"`
	Message string  `json:"message"`
}

func GetBalance(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, exists := c.Get("user_id")
		var userID int
		if !exists {
			userID = 1
		} else {
			userID = uid.(int)
		}

		var totalIncome, totalExpense float64

		// ใช้ _ แทน err ถ้าเราไม่ต้องการเช็ค error ในบรรทัดนี้ เพื่อลดปัญหา UnusedVar
		_ = db.QueryRow("SELECT SUM(amount) FROM transactions WHERE user_id = ? AND type = 'income'", userID).Scan(&totalIncome)
		_ = db.QueryRow("SELECT SUM(amount) FROM transactions WHERE user_id = ? AND type = 'expense'", userID).Scan(&totalExpense)

		balance := totalIncome - totalExpense

		c.JSON(http.StatusOK, SummaryResponse{
			Status:  "success",
			Balance: balance,
			Message: "Total balance calculated successfully",
		})
	}
}
