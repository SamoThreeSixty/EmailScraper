package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/samothreesixty/EmailScraper/internal/db"
)

func Connect() (*db.Queries, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	dsn := os.Getenv("DB_URL")

	dbPost, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("db connection failed: %w", err)
	}

	// Test the connection
	if err := dbPost.Ping(); err != nil {
		return nil, fmt.Errorf("db unreachable: %w", err)
	}

	dbConn := db.New(dbPost)

	return dbConn, nil
}
