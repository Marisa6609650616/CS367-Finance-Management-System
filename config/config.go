package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}