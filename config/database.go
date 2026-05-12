package config

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func ConnectDB() *sql.DB {
	dbPath := GetEnv("DB_PATH")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("Open DB error:", err)
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
