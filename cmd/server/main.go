package main

import (
	"CS367-Finance-Management-System/config"
	"CS367-Finance-Management-System/pkg/utils"
	"fmt"
)

func main() {
	config.LoadEnv()

	db := config.ConnectDB()
	defer db.Close()

	token, err := utils.GenerateToken(
		1,
		"test@email.com",
		"user",
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("JWT Token:")
	fmt.Println(token)
}