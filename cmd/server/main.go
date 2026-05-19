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

	// inject DB สำหรับ handler ที่ใช้ global var
	handlers.DB = db

	authHandler := &handlers.AuthHandler{DB: db}
	transactionHandler := handlers.NewTransactionHandler(db)
	summaryMonthlyHandler := handlers.NewSumaryMonthyHandler(db)
	balanceHandler := handlers.NewBalanceHandler(db)

	r := gin.Default()

	api := r.Group("/api")
	{
		// Public routes — ไม่ต้อง login
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes — ต้องมี JWT token
		protected := api.Group("/")
		protected.Use(middleware.RequireAuth())
		{
			// Transaction APIs (จันทร์พงศ์ + นัทธ์ชนัน)
			protected.POST("/transactions", transactionHandler.CreateTransaction)
			protected.PUT("/transactions/:id", handlers.UpdateTransaction)
			protected.DELETE("/transactions/:id", handlers.DeleteTransaction)

			// Summary APIs (ธนดล + มาริษา)
			protected.GET("/summary/monthly", summaryMonthlyHandler.SummaryMonthly)
			protected.GET("/summary/balance", balanceHandler.GetBalance)
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
