package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	LoadConfig()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		Config.DBHost, Config.DBPort, Config.DBUser, Config.DBPass, Config.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Cannot connect:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Ping error:", err)
	}

	log.Println("Connected to PostgreSQL!")
	DB = db
}
