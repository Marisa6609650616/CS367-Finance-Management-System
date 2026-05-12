package handlers

import (
	"CS367-Finance-Management-System/internal/models"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// DB จะถูก inject เข้ามาจาก main.go ตอน setup routes
var DB *sql.DB

// PUT /api/transactions/{id}
func UpdateTransaction(c *gin.Context) {
	userID := c.GetInt("userID")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid transaction ID."})
		return
	}

	var req models.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body."})
		return
	}

	// ต้องมีอย่างน้อย 1 field
	if req.Type == nil && req.Amount == nil && req.Category == nil &&
		req.Description == nil && req.Date == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No fields to update provided."})
		return
	}

	// Validate amount
	if req.Amount != nil && *req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Amount must be a positive number."})
		return
	}

	// Validate type
	if req.Type != nil && *req.Type != "income" && *req.Type != "expense" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Type must be income or expense."})
		return
	}

	// ตรวจสอบว่า transaction มีอยู่จริง
	var ownerID int
	err = DB.QueryRow("SELECT user_id FROM transactions WHERE id = ?", id).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Transaction not found."})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error."})
		return
	}

	// ตรวจสอบว่าเป็นเจ้าของ
	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden. You can only edit your own transactions."})
		return
	}

	// สร้าง dynamic UPDATE query
	setClauses := []string{}
	args := []interface{}{}

	if req.Type != nil {
		setClauses = append(setClauses, "type = ?")
		args = append(args, *req.Type)
	}
	if req.Amount != nil {
		setClauses = append(setClauses, "amount = ?")
		args = append(args, *req.Amount)
	}
	if req.Category != nil {
		setClauses = append(setClauses, "category = ?")
		args = append(args, *req.Category)
	}
	if req.Description != nil {
		setClauses = append(setClauses, "description = ?")
		args = append(args, *req.Description)
	}
	if req.Date != nil {
		setClauses = append(setClauses, "date = ?")
		args = append(args, *req.Date)
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE transactions SET %s WHERE id = ?", strings.Join(setClauses, ", "))

	if _, err := DB.Exec(query, args...); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error."})
		return
	}

	// ดึงข้อมูลที่อัปเดตแล้วกลับมา
	var t models.Transaction
	err = DB.QueryRow(
		`SELECT id, user_id, type, amount, category, description, date, created_at
		 FROM transactions WHERE id = ?`, id,
	).Scan(&t.ID, &t.UserID, &t.Type, &t.Amount, &t.Category, &t.Description, &t.Date, &t.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction updated successfully.",
		"data":    t,
	})
}

// DELETE /api/transactions/{id}
func DeleteTransaction(c *gin.Context) {
	userID := c.GetInt("userID")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid transaction ID."})
		return
	}

	// ตรวจสอบว่า transaction มีอยู่จริง
	var ownerID int
	err = DB.QueryRow("SELECT user_id FROM transactions WHERE id = ?", id).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Transaction not found."})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error."})
		return
	}

	// ตรวจสอบว่าเป็นเจ้าของ
	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden. You can only delete your own transactions."})
		return
	}

	if _, err := DB.Exec("DELETE FROM transactions WHERE id = ?", id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully."})
}
