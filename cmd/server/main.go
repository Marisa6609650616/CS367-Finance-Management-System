package main

import (
	"CS367-Finance-Management-System/config"
	"CS367-Finance-Management-System/handlers"
	"CS367-Finance-Management-System/middleware"
	"log"
	"net/http"
)

func main() {
	config.LoadEnv()
	db := config.ConnectDB()
	defer db.Close()

	transactionHandler := handlers.NewTransactionHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/transactions", middleware.AuthMiddleware(transactionHandler.CreateTransaction))

	log.Println("🚀 Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
