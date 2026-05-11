package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB() *sql.DB {
	dbPath := GetEnv("DB_PATH")

	dsn := fmt.Sprintf("file:%s?_foreign_keys=on", dbPath)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal("Open DB error:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	RunSchema(db)

	log.Println("Database connected")
	return db
}

func RunSchema(db *sql.DB) {
	schema, err := os.ReadFile("migrations/schema.sql")
	if err != nil {
		log.Fatal("Read schema error:", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal("Run schema error:", err)
	}

	log.Println("Schema loaded")
}