package routes

import (
	"CS367-Finance-Management-System/internal/handlers"
	"CS367-Finance-Management-System/internal/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupTransactionRoutes(r *gin.Engine, db *sql.DB) {
	handlers.DB = db

	tx := r.Group("/api/transactions")
	tx.Use(middleware.AuthMiddleware())
	{
		tx.PUT("/:id", handlers.UpdateTransaction)
		tx.DELETE("/:id", handlers.DeleteTransaction)
	}
}
