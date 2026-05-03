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

	r := gin.Default()

	authHandler := handlers.NewAuthHandler(db)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes — ใช้ JWT Middleware
		// (transaction + summary routes จะเพิ่มในสัปดาห์ที่ 2-3)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			_ = protected
		}
	}

	port := config.GetEnv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
