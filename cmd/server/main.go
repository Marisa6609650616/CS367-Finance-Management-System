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
	// transactionHandler := handlers.NewTransactionHandler(db)
	summaryHandler := handlers.NewSumaryMonthyHandler(db)

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
			// protected.POST("/api/transactions", transactionHandler.CreateTransaction)
			protected.PUT("/api/transactions/:id", handlers.UpdateTransaction)
			protected.DELETE("/api/transactions/:id", handlers.DeleteTransaction)
			protected.POST("/api/summary/monthly", summaryHandler.SummaryMonthly)
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
