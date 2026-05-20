package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func (h *SummaryMonthyHandler) SummaryMonthly(c *gin.Context) {
	userID := c.GetInt("user_id")

	month := c.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	rows, err := h.DB.Query(`
		SELECT type, SUM(amount)
		FROM transactions
		WHERE user_id = ?
		  AND strftime('%Y-%m', date) = ?
		GROUP BY type
	`, userID, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error."})
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Summary retrieved successfully.",
		"summary": SummaryMonthlyResponse{
			Month:   month,
			Income:  income,
			Expense: expense,
			Balance: income - expense,
		},
	})
}
