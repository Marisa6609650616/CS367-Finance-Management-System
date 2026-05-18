// cmd/server/main.go
package main

import (
	"CS367-Finance-Management-System/config"
	"CS367-Finance-Management-System/internal/handlers"
	"CS367-Finance-Management-System/internal/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	db := config.ConnectDB()
	defer db.Close()

	authHandler := &handlers.AuthHandler{DB: db}
	transactionHandler := handlers.NewTransactionHandler(db)
	summaryHandler := handlers.NewSumaryMonthyHandler(db)
	handlers.DB = db // สำหรับ UpdateTransaction และ DeleteTransaction

	r := gin.Default()
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
		protected := api.Group("/")
		protected.Use(middleware.RequireAuth())
		{
			protected.POST("/transactions", transactionHandler.CreateTransaction)
			protected.PUT("/transactions/:id", handlers.UpdateTransaction)
			protected.DELETE("/transactions/:id", handlers.DeleteTransaction)
			protected.POST("/summary/monthly", summaryHandler.SummaryMonthly)
		}
	}

	port := config.GetEnv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running → http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed:", err)
	}
}
